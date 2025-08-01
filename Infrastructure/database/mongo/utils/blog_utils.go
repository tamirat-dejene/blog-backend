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

// BuildBlogPostFilterQuery constructs a MongoDB query based on the provided BlogPostFilter.
// It filters blog posts by tags, title, and author name.
func BuildBlogPostFilterQuery(filter *domain.BlogPostFilter) bson.M {
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

// CalculatePopularityScore computes the popularity score based on likes, views, comments, and dislikes.
func CalculatePopularityScore(likes, views, comments, dislikes int) float64 {
	raw := (float64(likes) * 3.0) + (float64(views) * 2.0) + (float64(comments) * 1.5) - (float64(dislikes) * 2.5)
	maxScore := 50000.0  // assumed maximum score for normalization
	normalized := (raw / maxScore) * 100
	if normalized < 0 {
		return 0
	}
	if normalized > 100 {
		return 100
	}
	return normalized
}