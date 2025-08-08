package repository

import (
	"context"
	"errors"
	"fmt"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/database/mongo"
	"g6/blog-api/Infrastructure/database/mongo/mapper"
	"g6/blog-api/Infrastructure/database/mongo/utils"
	"net/http"

	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type blogPostRepo struct {
	db          mongo.Database
	collections *mongo.Collections
}

// Create implements domain.BlogRepository.
func (b *blogPostRepo) Create(ctx context.Context, blog *domain.BlogPost) (*domain.BlogPost, *domain.DomainError) {
	// Map the domain model to the DB model
	blogModel := &mapper.BlogPostModel{}
	if err := blogModel.Parse(blog); err != nil {
		return nil, &domain.DomainError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	// Insert the blog into the collection
	result, err := b.db.Collection(b.collections.BlogPosts).InsertOne(ctx, blogModel)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	// Extract the inserted ID
	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, &domain.DomainError{
			Err:  errors.New("failed to cast inserted ID to ObjectID"),
			Code: http.StatusInternalServerError,
		}
	}

	// Set the generated ID back to domain model
	blog.ID = objectID.Hex()
	return blog, nil
}

// Delete implements domain.BlogRepository.
func (b *blogPostRepo) Delete(ctx context.Context, id string) *domain.DomainError {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &domain.DomainError{
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}

	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return &domain.DomainError{
			Err:  errors.New("unauthorized: missing user ID"),
			Code: http.StatusUnauthorized,
		}
	}

	role, _ := ctx.Value("role").(string)
	isAdmin := role == "admin" || role == "superadmin"

	var filter bson.M

	if isAdmin {
		// Admins can delete any blog post
		filter = bson.M{"_id": oid}
	} else {
		// Regular users can only delete their own blog post
		authorOID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			return &domain.DomainError{
				Err:  errors.New("invalid user ID"),
				Code: http.StatusBadRequest,
			}
		}
		filter = bson.M{"_id": oid, "author_id": authorOID}
	}

	result, err := b.db.Collection(b.collections.BlogPosts).DeleteOne(ctx, filter)
	if err != nil {
		return &domain.DomainError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	if result == 0 {
		return &domain.DomainError{
			Err:  fmt.Errorf("not found or not authorized to delete blog post with ID %s", id),
			Code: http.StatusNotFound,
		}
	}

	return nil
}

// Get implements domain.BlogRepository.
func (b *blogPostRepo) Get(ctx context.Context, filter *domain.BlogPostFilter) ([]domain.BlogPostsPage, *string, *domain.DomainError) {
	collection := b.db.Collection(b.collections.BlogPosts)
	pipeline := utils.BuildBlogRetrievalAggregationPipeline(filter)

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, nil, &domain.DomainError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	defer cursor.Close(ctx)

	var dbResults []mapper.BlogPostModel
	if err := cursor.All(ctx, &dbResults); err != nil {
		return nil, nil, &domain.DomainError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	if len(dbResults) == 0 {
		return nil, nil, &domain.DomainError{
			Err:  errors.New("no blog posts found"),
			Code: http.StatusNotFound,
		}
	}

	// Serialize the results for caching
	serialized, err := utils.SerializeBlogPostsPage(&dbResults)

	if err != nil {
		return nil, nil, &domain.DomainError{
			Err:  errors.New("failed to serialize blog posts page"),
			Code: http.StatusInternalServerError,
		}
	}

	return utils.PaginateBlogs(dbResults, filter.PageSize), &serialized, nil
}

// GetBlogByID implements domain.BlogRepository.
func (b *blogPostRepo) GetBlogByID(ctx context.Context, id string) (*domain.BlogPost, *domain.DomainError) {
	// Validate the ID format
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}

	// Fetch the blog post from the database
	var blogModel *mapper.BlogPostModel
	err = b.db.Collection(b.collections.BlogPosts).FindOne(ctx, bson.M{"_id": oid}).Decode(&blogModel)
	if err == mongo.ErrNoDocuments() {
		return nil, &domain.DomainError{
			Err:  err,
			Code: http.StatusNotFound,
		}
	} else if err != nil {
		return nil, &domain.DomainError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	return blogModel.ToDomain(), nil
}

// Update implements domain.BlogRepository.
func (b *blogPostRepo) Update(ctx context.Context, id string, blog domain.BlogPost) (*domain.BlogPost, *domain.DomainError) {
	// Convert blog ID from hex
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("invalid blog ID: %w", err),
			Code: http.StatusBadRequest,
		}
	}

	// Extract user_id from context and assert it's an ObjectID
	userID, ok := ctx.Value("user_id").(primitive.ObjectID)
	if !ok {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("invalid or missing user_id in context"),
			Code: http.StatusUnauthorized,
		}
	}

	// Filter must ensure both blog ID and author ID match
	filter := bson.M{
		"_id":       oid,
		"author_id": userID,
	}

	// Set update fields
	blog.UpdatedAt = time.Now()
	update := bson.M{
		"$set": bson.M{
			"title":      blog.Title,
			"content":    blog.Content,
			"tags":       blog.Tags,
			"updated_at": primitive.NewDateTimeFromTime(blog.UpdatedAt),
		},
	}

	// Perform update
	res, err := b.db.Collection(b.collections.BlogPosts).UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("failed to update blog post: %w", err),
			Code: http.StatusInternalServerError,
		}
	}

	if res.MatchedCount == 0 {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("blog not found or unauthorized to update"),
			Code: http.StatusForbidden,
		}
	}

	// Return updated blog
	blog.ID = oid.Hex()
	blog.AuthorID = userID.Hex()
	return &blog, nil
}

// RefreshPopularityScore implements domain.BlogRepository.
func (b *blogPostRepo) RefreshPopularityScore(ctx context.Context, id string) (*domain.BlogPost, *domain.DomainError) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}

	var blogModel mapper.BlogPostModel
	err = b.db.Collection(b.collections.BlogPosts).FindOne(ctx, bson.M{"_id": oid}).Decode(&blogModel)
	if err == mongo.ErrNoDocuments() {
		return nil, &domain.DomainError{
			Err:  fmt.Errorf("blog post with ID %s not found", id),
			Code: http.StatusNotFound,
		}
	} else if err != nil {
		return nil, &domain.DomainError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	ps := utils.CalculatePopularityScore(blogModel.Likes, blogModel.ViewCount, blogModel.CommentCount, blogModel.Dislikes)
	_, err = b.db.Collection(b.collections.BlogPosts).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{
		"$set": bson.M{"popularity_score": ps},
	})

	if err != nil {
		return nil, &domain.DomainError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	blogModel.PopularityScore = ps
	return blogModel.ToDomain(), nil
}

// IncrementViewCount implements domain.BlogRepository.
func (b *blogPostRepo) IncrementViewCount(ctx context.Context, id string) (*domain.BlogPost, *domain.DomainError) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}

	// Increment the view count
	update := bson.M{"$inc": bson.M{"view_count": 1}}
	_, err = b.db.Collection(b.collections.BlogPosts).UpdateOne(ctx, bson.M{"_id": oid}, update)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	return b.GetBlogByID(ctx, id)
}

// UpdateCommentCount implements domain.BlogRepository.
func (b *blogPostRepo) UpdateCommentCount(ctx context.Context, id string, increment bool) (*domain.BlogPost, *domain.DomainError) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}

	var update bson.M
	if increment {
		update = bson.M{"$inc": bson.M{"comment_count": 1}}
	} else {
		update = bson.M{"$inc": bson.M{"comment_count": -1}}
	}

	_, err = b.db.Collection(b.collections.BlogPosts).UpdateOne(ctx, bson.M{"_id": oid}, update)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	return b.GetBlogByID(ctx, id)
}

// UpdateReactionCount implements domain.BlogRepository.
func (b *blogPostRepo) UpdateReactionCount(ctx context.Context, is_like bool, id string, increment bool) (*domain.BlogPost, *domain.DomainError) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}

	field := "dislikes"
	if is_like {
		field = "likes"
	}
	inc := -1
	if increment {
		inc = 1
	}
	update := bson.M{"$inc": bson.M{field: inc}}

	_, err = b.db.Collection(b.collections.BlogPosts).UpdateOne(ctx, bson.M{"_id": oid}, update)
	if err != nil {
		return nil, &domain.DomainError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	return b.GetBlogByID(ctx, id)
}

// NewBlogPostRepo creates a new instance of blogPostRepo.
func NewBlogPostRepo(database mongo.Database, collections *mongo.Collections) domain.BlogPostRepository {
	return &blogPostRepo{
		db:          database,
		collections: collections,
	}
}
