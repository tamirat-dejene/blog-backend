package dto

import domain "g6/blog-api/Domain"

type AIBlogPostGenerateRequest struct {
	Topic    string   `json:"topic" binding:"required"`
	Keywords []string `json:"keywords" binding:"required"`
}
type AIBlogPostResponse struct {
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

func (a *AIBlogPostResponse) ToDomain() *domain.AIBlogPost {
	return &domain.AIBlogPost{
		Topic:           a.Topic,
		Keywords:        a.Keywords,
		Title:           a.Title,
		Introduction:    a.Introduction,
		Body:            a.Body,
		Conclusion:      a.Conclusion,
		SuggestedTitles: a.SuggestedTitles,
		RelatedIdeas:    a.RelatedIdeas,
		CreatedAt:       a.CreatedAt,
	}
}

func (a *AIBlogPostResponse) FromDomain(domainPost *domain.AIBlogPost) {
	a.Topic = domainPost.Topic
	a.Keywords = domainPost.Keywords
	a.Title = domainPost.Title
	a.Introduction = domainPost.Introduction
	a.Body = domainPost.Body
	a.Conclusion = domainPost.Conclusion
	a.SuggestedTitles = domainPost.SuggestedTitles
	a.RelatedIdeas = domainPost.RelatedIdeas
	a.CreatedAt = domainPost.CreatedAt
}