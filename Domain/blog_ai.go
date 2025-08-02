package domain

import (
	"context"
)

type BlogAIGenerate struct {
	UserID   string
	Topic    string
	Keywords []string
}

type BlogAIContent struct {
	ID              string
	UserID          string
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

type BlogAIFeedback struct {
	ID        string
	UserID    string
	ContentID string
	Rating    int
	Feedback  string
	CreatedAt string
}

type BlogAIUsecase interface {
	GenerateContent(ctx context.Context, req BlogAIGenerate) (*BlogAIContent, *DomainError)
	GetGeneratedContentByID(ctx context.Context, id string) (*BlogAIContent, *DomainError)
	SubmitFeedback(ctx context.Context, feedback BlogAIFeedback) (*BlogAIFeedback, *DomainError)
}

type BlogAIRepository interface {
	StoreGeneratedContent(ctx context.Context, content *BlogAIContent) (*BlogAIContent, *DomainError)
	GetGeneratedContentByID(ctx context.Context, id string) (*BlogAIContent, *DomainError)
	SaveFeedback(ctx context.Context, feedback BlogAIFeedback) (*BlogAIFeedback, *DomainError)
}
