package domain

import (
	"context"
)

type AIBlogPostGenerate struct {
	Topic    string
	Keywords []string
}

type AIBlogPost struct {
	Topic           string
	Keywords        []string
	Title           string
	Introduction    string
	Body            string
	Conclusion      string
	SuggestedTitles []string
	RelatedIdeas    []string
	CreatedAt       string
}

type AIBlogPostUsecase interface {
	GeneratePost(ctx context.Context, req AIBlogPostGenerate) (*AIBlogPost, *DomainError)
}