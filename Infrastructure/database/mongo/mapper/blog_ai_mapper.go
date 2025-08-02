package mapper

import (
	domain "g6/blog-api/Domain"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- MongoDB Models ---

type BlogAIContentModel struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	UserID          primitive.ObjectID `bson:"user_id"`
	Topic           string             `bson:"topic"`
	Keywords        []string           `bson:"keywords"`
	Title           string             `bson:"title"`
	Introduction    string             `bson:"introduction"`
	Body            string             `bson:"body"`
	Conclusion      string             `bson:"conclusion"`
	SuggestedTitles []string           `bson:"suggested_titles"`
	RelatedIdeas    []string           `bson:"related_ideas"`
	CreatedAt       primitive.DateTime `bson:"created_at"`
}

type BlogAIFeedbackModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id"`
	ContentID primitive.ObjectID `bson:"content_id"`
	Rating    int                `bson:"rating"`
	Feedback  string             `bson:"feedback"`
	CreatedAt primitive.DateTime `bson:"created_at"`
}

// --- Mapper Functions ---

func BlogAIContentToDomain(content *BlogAIContentModel) *domain.BlogAIContent {
	return &domain.BlogAIContent{
		ID:              content.ID.Hex(),
		UserID:          content.UserID.Hex(),
		Topic:           content.Topic,
		Keywords:        content.Keywords,
		Title:           content.Title,
		Introduction:    content.Introduction,
		Body:            content.Body,
		Conclusion:      content.Conclusion,
		SuggestedTitles: content.SuggestedTitles,
		RelatedIdeas:    content.RelatedIdeas,
		CreatedAt:       content.CreatedAt.Time().Format(time.RFC3339),
	}
}

func BlogAIContentFromDomain(content *domain.BlogAIContent) (*BlogAIContentModel, error) {
	userID, err := primitive.ObjectIDFromHex(content.UserID)
	if err != nil {
		return nil, err
	}

	objectID := primitive.NewObjectID()
	if content.ID != "" {
		objectID, err = primitive.ObjectIDFromHex(content.ID)
		if err != nil {
			return nil, err
		}
	}

	createdAtTime, err := time.Parse(time.RFC3339, content.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &BlogAIContentModel{
		ID:              objectID,
		UserID:          userID,
		Topic:           content.Topic,
		Keywords:        content.Keywords,
		Title:           content.Title,
		Introduction:    content.Introduction,
		Body:            content.Body,
		Conclusion:      content.Conclusion,
		SuggestedTitles: content.SuggestedTitles,
		RelatedIdeas:    content.RelatedIdeas,
		CreatedAt:       primitive.NewDateTimeFromTime(createdAtTime),
	}, nil
}

func BlogAIFeedbackToDomain(feedback *BlogAIFeedbackModel) *domain.BlogAIFeedback {
	return &domain.BlogAIFeedback{
		ID:        feedback.ID.Hex(),
		UserID:    feedback.UserID.Hex(),
		ContentID: feedback.ContentID.Hex(),
		Rating:    feedback.Rating,
		Feedback:  feedback.Feedback,
		CreatedAt: feedback.CreatedAt.Time().Format(time.RFC3339),
	}
}
func BlogAIFeedbackFromDomain(feedback *domain.BlogAIFeedback) (*BlogAIFeedbackModel, error) {
	userID, err := primitive.ObjectIDFromHex(feedback.UserID)
	if err != nil {
		return nil, err
	}

	contentID, err := primitive.ObjectIDFromHex(feedback.ContentID)
	if err != nil {
		return nil, err
	}

	objectID := primitive.NewObjectID()
	if feedback.ID != "" {
		objectID, err = primitive.ObjectIDFromHex(feedback.ID)
		if err != nil {
			return nil, err
		}
	}

	createdAtTime, err := time.Parse(time.RFC3339, feedback.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &BlogAIFeedbackModel{
		ID:        objectID,
		UserID:    userID,
		ContentID: contentID,
		Rating:    feedback.Rating,
		Feedback:  feedback.Feedback,
		CreatedAt: primitive.NewDateTimeFromTime(createdAtTime),
	}, nil
}
