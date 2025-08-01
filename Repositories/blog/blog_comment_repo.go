package repository

import (
	"context"
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/database/mongo"
	"g6/blog-api/Infrastructure/database/mongo/mapper"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type blogCommentRepo struct {
	db          mongo.Database
	collections *collections
}

func (b *blogCommentRepo) Create(ctx context.Context, comment domain.BlogComment) (domain.BlogComment, error) {
	blogComment, err := mapper.BlogCommentFromDomain(&comment)

	if err != nil {
		return domain.BlogComment{}, err
	}
	_, err = b.db.Collection(b.collections.BlogComments).InsertOne(ctx, blogComment)

	if err != nil {
		return domain.BlogComment{}, err
	}
	res := mapper.BlogCommentToDomain(blogComment)
	return *res, nil
}

func (b *blogCommentRepo) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}
	_, err = b.db.Collection(b.collections.BlogComments).DeleteOne(ctx, filter)

	return err
}

func NewBlogCommentRepo(db mongo.Database, col *collections) domain.BlogCommentRepository {
	return &blogCommentRepo{
		db:          db,
		collections: col,
	}

}
