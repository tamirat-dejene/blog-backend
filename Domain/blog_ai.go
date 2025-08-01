package domain

import "context"

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
	GenerateContent(ctx context.Context, req BlogAIGenerate) (*BlogAIContent, error)
	GetGeneratedContentByID(ctx context.Context, id string) (*BlogAIContent, error)
	GetUserPromptHistory(ctx context.Context, userID string) ([]BlogAIPrompt, error)
	SubmitFeedback(ctx context.Context, feedback BlogAIFeedback) error
}

type BlogAIRepository interface {
	StoreGeneratedContent(ctx context.Context, content *BlogAIContent) error
	GetGeneratedContentByID(ctx context.Context, id string) (*BlogAIContent, error)
	GetPromptsByUserID(ctx context.Context, userID string) ([]BlogAIPrompt, error)
	SaveFeedback(ctx context.Context, feedback BlogAIFeedback) error
}
