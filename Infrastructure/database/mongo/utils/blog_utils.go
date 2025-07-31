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
func PaginationOpts(page, page_size int, recency domain.Recency) *options.FindOptions {
	if page <= 0 {
		page = 1
	}
	if page_size <= 0 {
		page_size = 10
	}
	skip := int64((page - 1) * page_size)

	return options.Find().
		SetSkip(skip).
		SetLimit(int64(page_size)).
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

// popularity score pipeline

// PopularityStages returns MongoDB aggregation stages for computing and sorting by popularity score.
func PopularityStages() []bson.D {
	return []bson.D{
		{{
			Key: "$addFields", Value: bson.M{
				"popularity_score": bson.M{
					"$add": bson.A{
						"$comment_count",
						"$likes",
						bson.M{"$multiply": bson.A{-1, "$dislikes"}},
					},
				},
			},
		}},
		{{
			Key: "$sort", Value: bson.M{"popularity_score": -1},
		}},
	}
}
