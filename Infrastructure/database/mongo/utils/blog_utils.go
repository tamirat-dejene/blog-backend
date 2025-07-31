package utils

import (
	domain "g6/blog-api/Domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// RecencySort returns a sort option based on the recency type.
// It sorts by the "created_at" field in ascending or descending order.
func RecencySort(recency domain.Recency) bson.D {
	switch recency {
	case domain.RecencyOldest:
		return bson.D{{Key: "created_at", Value: 1}} // ascending
	default: // domain.Newest or empty
		return bson.D{{Key: "created_at", Value: -1}} // descending
	}
}

// PaginationOpts returns MongoDB options for pagination based on the provided page, pageSize, and recency.
func PaginationOpts(page, pageSize int, recency domain.Recency) *options.FindOptions {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	skip := int64((page - 1) * pageSize)

	return options.Find().
		SetSkip(skip).
		SetLimit(int64(pageSize)).
		SetSort(RecencySort(recency))
}

// BuildBlogFilterQuery constructs a MongoDB query based on the provided BlogFilter.
// It filters blogs by tags, title, and author name.
func BuildBlogFilterQuery(filter *domain.BlogFilter) bson.M {
	query := bson.M{}

	if filter == nil {
		return query
	}

	if len(filter.Tags) > 0 {
		query["tags"] = bson.M{"$in": filter.Tags}
	}

	if filter.Title != "" {
		query["title"] = bson.M{
			"$regex": primitive.Regex{Pattern: filter.Title, Options: "i"},
		}
	}

	if filter.AuthorName != "" {
		query["author_name"] = bson.M{
			"$regex": primitive.Regex{Pattern: filter.AuthorName, Options: "i"},
		}
	}

	return query
}
