package redis

import (
	"fmt"
	domain "g6/blog-api/Domain"
	"sort"
	"strings"
)

type RedisService struct{}

func (r *RedisService) GenerateRedisKey(filter *domain.BlogPostFilter) string {
	sort.Strings(filter.Tags)
	tags := strings.Join(filter.Tags, ",")
	return fmt.Sprintf("blogs:page=%d:size=%d:recency=%s:tags=%s:author=%s:title=%s:popular=%t",
		filter.Page,
		filter.PageSize,
		filter.Recency,
		tags,
		filter.AuthorName,
		filter.Title,
		filter.Popular,
	)
}

func (r *RedisService) GenerateBlogPostKey(id string) string {
	return fmt.Sprintf("blogpost:%s", id)
}

func (r *RedisService) GenerateBlogPostCommentsKey(blogID string) string {
	return fmt.Sprintf("blogpost:%s:comments", blogID)
}

func (r *RedisService) GenerateBlogPostReactionsKey(blogID string) string {
	return fmt.Sprintf("blogpost:%s:reactions", blogID)
}

func (r *RedisService) GenerateBlogPostAuthorKey(authorID string) string {
	return fmt.Sprintf("blogpost:author:%s", authorID)
}
