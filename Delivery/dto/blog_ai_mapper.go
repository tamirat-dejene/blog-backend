package dto

import (
	domain "g6/blog-api/Domain"
	"time"
)

type BlogAIGenerateRequest struct {
	Topic    string   `json:"topic" binding:"required"`
	Keywords []string `json:"keywords" binding:"required"`
}

type BlogAIResponseDTO struct {
	ID              string   `json:"id"`
	Topic           string   `json:"topic"`
	Keywords        []string `json:"keywords"`
	Title           string   `json:"title"`
	Introduction    string   `json:"introduction"`
	Body            string   `json:"body"`
	Conclusion      string   `json:"conclusion"`
	SuggestedTitles []string `json:"suggested_titles"`
	RelatedIdeas    []string `json:"related_ideas"`
	CreatedAt       string   `json:"created_at"`
}

type BlogAIFeedbackRequest struct {
	ContentID string `json:"content_id" binding:"required"`
	Rating    int    `json:"rating" binding:"required"`
	Feedback  string `json:"feedback" binding:"required"`
}

func BlogAIContentToDomain(content *BlogAIResponseDTO, userID string, createdAtTime time.Time) *domain.BlogAIContent {
	return &domain.BlogAIContent{
		ID:              content.ID,
		UserID:          userID,
		Topic:           content.Topic,
		Keywords:        content.Keywords,
		Title:           content.Title,
		Introduction:    content.Introduction,
		Body:            content.Body,
		Conclusion:      content.Conclusion,
		SuggestedTitles: content.SuggestedTitles,
		RelatedIdeas:    content.RelatedIdeas,
		CreatedAt:       createdAtTime.Format(time.RFC3339),
	}
}

func BlogAIContentFromDomain(content *domain.BlogAIContent) *BlogAIResponseDTO {
	return &BlogAIResponseDTO{
		ID:              content.ID,
		Topic:           content.Topic,
		Keywords:        content.Keywords,
		Title:           content.Title,
		Introduction:    content.Introduction,
		Body:            content.Body,
		Conclusion:      content.Conclusion,
		SuggestedTitles: content.SuggestedTitles,
		RelatedIdeas:    content.RelatedIdeas,
		CreatedAt:       content.CreatedAt,
	}
}