package utils

import (
	"errors"
	"fmt"
	"g6/blog-api/Infrastructure/database/mongo/mapper"

	"go.mongodb.org/mongo-driver/bson"
)

func SerializeBlogComment(comment mapper.BlogCommentModel) (string, error) {
	data, err := bson.MarshalExtJSON(comment, false, false)
	if err != nil {
		return "", fmt.Errorf("failed to serialize blog comment to BSON: %w", err)
	}
	return string(data), nil
}

func DeserializeBlogComment(serialized string) (mapper.BlogCommentModel, error) {
	var comment mapper.BlogCommentModel
	err := bson.UnmarshalExtJSON([]byte(serialized), true, &comment)
	if err != nil {
		return mapper.BlogCommentModel{}, fmt.Errorf("failed to deserialize blog comment: %w", err)
	}
	return comment, nil
}

type blogCommentListWrapper struct {
	Comments []mapper.BlogCommentModel `json:"comments"`
}

func SerializeBlogCommentList(comments []mapper.BlogCommentModel) (string, error) {
	wrapper := blogCommentListWrapper{Comments: comments}
	data, err := bson.MarshalExtJSON(wrapper, false, false)
	if err != nil {
		return "", errors.New("failed to serialize blog comment list")
	}
	return string(data), nil
}

func DeserializeBlogCommentList(serialized string) ([]mapper.BlogCommentModel, error) {
	if serialized == "" {
		return nil, fmt.Errorf("empty serialized string")
	}

	var wrapper blogCommentListWrapper
	err := bson.UnmarshalExtJSON([]byte(serialized), true, &wrapper)
	if err != nil {
		return nil, errors.New("failed to deserialize blog comment list")
	}
	return wrapper.Comments, nil
}
