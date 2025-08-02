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
	GeneratedText   string
	SuggestedTitles []string
	RelatedIdeas    []string
	CreatedAt       string
}

type BlogAIFeedback struct {
	ContentID string
	UserID    string
	Rating    int
	Feedback  string
}

type BlogAIPrompt struct {
	ID        string
	UserID    string
	Topic     string
	CreatedAt string
}

type BlogAIUsecase interface {
	GenerateContent(ctx context.Context, req BlogAIGenerate) (*BlogAIContent, *DomainError)
	GetGeneratedContentByID(ctx context.Context, id string) (*BlogAIContent, *DomainError)
	GetUserPromptHistory(ctx context.Context, userID string) ([]BlogAIPrompt, *DomainError)
	SubmitFeedback(ctx context.Context, feedback BlogAIFeedback) *DomainError
}

type BlogAIRepository interface {
	StoreGeneratedContent(ctx context.Context, content *BlogAIContent) *DomainError
	GetGeneratedContentByID(ctx context.Context, id string) (*BlogAIContent, *DomainError)
	GetPromptsByUserID(ctx context.Context, userID string) ([]BlogAIPrompt, *DomainError)
	SaveFeedback(ctx context.Context, feedback BlogAIFeedback) *DomainError
}
