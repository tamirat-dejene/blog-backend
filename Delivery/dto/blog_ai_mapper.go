package dto

type BlogAIGenerateRequest struct {
	Topic    string   `json:"topic" binding:"required"`
	Keywords []string `json:"keywords" binding:"required"`
}

type BlogAIResponseDTO struct {
	ID              string   `json:"id"`
	Topic           string   `json:"topic"`
	Keywords        []string `json:"keywords"`
	GeneratedText   string   `json:"generated_text"`
	SuggestedTitles []string `json:"suggested_titles"`
	RelatedIdeas    []string `json:"related_ideas"`
	CreatedAt       string   `json:"created_at"`
}

type BlogAIFeedbackRequest struct {
	ContentID string `json:"content_id" binding:"required"`
	Rating    int    `json:"rating" binding:"required"`
	Feedback  string `json:"feedback" binding:"required"`
}
