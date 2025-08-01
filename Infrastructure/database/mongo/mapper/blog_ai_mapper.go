package mapper

import (
	"g6/blog-api/Domain"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- MongoDB Models ---

type BlogAIContentModel struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	UserID          primitive.ObjectID `bson:"user_id"`
	Topic           string             `bson:"topic"`
	Keywords        []string           `bson:"keywords"`
	GeneratedText   string             `bson:"generated_text"`
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
		GeneratedText:   content.GeneratedText,
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
		GeneratedText:   content.GeneratedText,
		SuggestedTitles: content.SuggestedTitles,
		RelatedIdeas:    content.RelatedIdeas,
		CreatedAt:       primitive.NewDateTimeFromTime(createdAtTime),
	}, nil
}
