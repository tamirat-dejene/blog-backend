package mapper

import (
	"fmt"
	domain "g6/blog-api/Domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogPostModel struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	Title           string             `bson:"title"`
	Content         string             `bson:"content"`
	AuthorID        primitive.ObjectID `bson:"author_id"`
	AuthorName      string             `bson:"author_name"` // for easy access to author's name: first_name + last_name
	Tags            []string           `bson:"tags,omitempty"`
	CreatedAt       primitive.DateTime `bson:"created_at"`
	UpdatedAt       primitive.DateTime `bson:"updated_at"`
	Likes           int                `bson:"likes"`
	Dislikes        int                `bson:"dislikes"`
	ViewCount       int                `bson:"view_count"`
	CommentCount    int                `bson:"comment_count"`    // for easy access to comment count
	PopularityScore float64            `bson:"popularity_score"` // computed popularity score
}

type BlogCommentModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	BlogID    primitive.ObjectID `bson:"blog_id"`
	AuthorID  primitive.ObjectID `bson:"author_id"`
	Comment   string             `bson:"comment"`
	CreatedAt primitive.DateTime `bson:"created_at"`
}

type BlogUserReactionModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	BlogID    primitive.ObjectID `bson:"blog_id"`
	UserID    primitive.ObjectID `bson:"user_id"`
	IsLike    bool               `bson:"is_like"`
	CreatedAt primitive.DateTime `bson:"created_at"`
}

type ObjectIDModel struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
}

func (b *BlogPostModel) Parse(bp *domain.BlogPost) error {
	b.Title = bp.Title
	b.Content = bp.Content
	authorID, err := primitive.ObjectIDFromHex(bp.AuthorID)
	if err != nil {
		return fmt.Errorf("invalid author ID: %w", err)
	}
	b.AuthorID = authorID
	b.AuthorName = bp.AuthorName
	b.Tags = bp.Tags
	b.CreatedAt = primitive.NewDateTimeFromTime(bp.CreatedAt)
	b.UpdatedAt = primitive.NewDateTimeFromTime(bp.UpdatedAt)
	b.Likes = bp.Likes
	b.Dislikes = bp.Dislikes
	b.ViewCount = bp.ViewCount
	b.CommentCount = bp.CommentCount
	b.PopularityScore = bp.PopularityScore
	return nil
}

func (b *BlogPostModel) ToDomain() *domain.BlogPost {
	return &domain.BlogPost{
		ID:              b.ID.Hex(),
		Title:           b.Title,
		Content:         b.Content,
		AuthorID:        b.AuthorID.Hex(),
		AuthorName:      b.AuthorName,
		Tags:            b.Tags,
		CreatedAt:       b.CreatedAt.Time(),
		UpdatedAt:       b.UpdatedAt.Time(),
		Likes:           b.Likes,
		Dislikes:        b.Dislikes,
		ViewCount:       b.ViewCount,
		CommentCount:    b.CommentCount,
		PopularityScore: b.PopularityScore,
	}
}

func (c *BlogCommentModel) Parse(comment *domain.BlogComment) error {
	c.Comment = comment.Comment
	blogID, err := primitive.ObjectIDFromHex(comment.BlogID)
	if err != nil {
		return fmt.Errorf("invalid blog ID: %w", err)
	}
	c.BlogID = blogID
	authorID, err := primitive.ObjectIDFromHex(comment.AuthorID)
	if err != nil {
		return fmt.Errorf("invalid author ID: %w", err)
	}
	c.AuthorID = authorID
	c.CreatedAt = primitive.NewDateTimeFromTime(comment.CreatedAt)

	if cid, err := primitive.ObjectIDFromHex(comment.ID); err == nil {
		c.ID = cid
	}
	return nil
}

func (c *BlogCommentModel) ToDomain() *domain.BlogComment {
	return &domain.BlogComment{
		ID:        c.ID.Hex(),
		BlogID:    c.BlogID.Hex(),
		AuthorID:  c.AuthorID.Hex(),
		Comment:   c.Comment,
		CreatedAt: c.CreatedAt.Time(),
	}
}

func (b *BlogUserReactionModel) Parse(reaction *domain.BlogUserReaction) error {
	blogID, err := primitive.ObjectIDFromHex(reaction.BlogID)
	if err != nil {
		return fmt.Errorf("invalid blog ID: %w", err)
	}
	b.BlogID = blogID

	userID, err := primitive.ObjectIDFromHex(reaction.UserID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}
	b.UserID = userID

	b.IsLike = reaction.IsLike
	b.CreatedAt = primitive.NewDateTimeFromTime(reaction.CreatedAt)
	return nil
}

func (b *BlogUserReactionModel) ToDomain() *domain.BlogUserReaction {
	return &domain.BlogUserReaction{
		ID:        b.ID.Hex(),
		BlogID:    b.BlogID.Hex(),
		UserID:    b.UserID.Hex(),
		IsLike:    b.IsLike,
		CreatedAt: b.CreatedAt.Time(),
	}
}
