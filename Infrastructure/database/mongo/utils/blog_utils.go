package utils

import (
	domain "g6/blog-api/Domain"
	"g6/blog-api/Infrastructure/database/mongo/mapper"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	maxScore := 50000.0 // assumed maximum score for normalization
	normalized := (raw / maxScore) * 100
	if normalized < 0 {
		return 0
	}
	if normalized > 100 {
		return 100
	}
	return normalized
}

// BuildBlogRetrievalAggregationPipeline constructs an aggregation pipeline for retrieving blog posts.
func BuildBlogRetrievalAggregationPipeline(filter *domain.BlogPostFilter) []bson.D {
	query := BuildBlogPostFilterQuery(filter)
	pipeline := []bson.D{
		{{Key: "$match", Value: query}},
	}

	if filter.Popular {
		pipeline = append(pipeline, bson.D{{Key: "$sort", Value: bson.D{{Key: "popularity_score", Value: -1}}}})
	} else if filter.Recency != "" {
		pipeline = append(pipeline, bson.D{{Key: "$sort", Value: RecencySort(filter.Recency)}})
	} else {
		pipeline = append(pipeline, bson.D{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}})
	}

	pipeline = append(pipeline,
		// bson.D{{Key: "$skip", Value: max((filter.Page-1)*filter.PageSize, 0)}},
		bson.D{{Key: "$limit", Value: filter.PageSize * filter.Page}},
	)

	return pipeline
}

// PaginateBlogs takes a slice of BlogPostModel and returns a paginated result based on the page size.
func PaginateBlogs(blogs []mapper.BlogPostModel, pageSize int) []domain.BlogPostsPage {
	totalBlogs := len(blogs)
	totalPages := (totalBlogs + pageSize - 1) / pageSize

	var pages []domain.BlogPostsPage

	for page := 1; page <= totalPages; page++ {
		start := (page - 1) * pageSize
		end := min(start+pageSize, totalBlogs)

		paginatedBlogs := blogs[start:end]
		domainBlogs := make([]domain.BlogPost, len(paginatedBlogs))
		for i, blog := range paginatedBlogs {
			domainBlogs[i] = *blog.ToDomain()
		}

		pageObj := domain.BlogPostsPage{
			Blogs:      domainBlogs,
			PageNumber: page,
			PageSize:   len(domainBlogs),
		}
		pages = append(pages, pageObj)
	}

	return pages
}
