package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"g6/blog-api/Infrastructure/database/mongo/mapper"
)

func SerializeBlogComment(comment *mapper.BlogCommentModel) (string, error) {
	data, err := json.Marshal(comment)
	if err != nil {
		return "", fmt.Errorf("failed to serialize blog comment: %w", err)
	}
	return string(data), nil
}

func DeserializeBlogComment(serialized string) (*mapper.BlogCommentModel, error) {
	var comment mapper.BlogCommentModel
	err := json.Unmarshal([]byte(serialized), &comment)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize blog comment: %w", err)
	}
	return &comment, nil
}

type blogCommentListWrapper struct {
	Comments []*mapper.BlogCommentModel `json:"comments"`
}

func SerializeBlogCommentList(comments []*mapper.BlogCommentModel) (string, error) {
	wrapper := blogCommentListWrapper{Comments: comments}
	data, err := json.Marshal(wrapper)
	if err != nil {
		return "", errors.New("failed to serialize blog comment list")
	}
	return string(data), nil
}

func DeserializeBlogCommentList(serialized string) ([]*mapper.BlogCommentModel, error) {
	if serialized == "" {
		return nil, fmt.Errorf("empty serialized string")
	}

	var wrapper blogCommentListWrapper
	err := json.Unmarshal([]byte(serialized), &wrapper)
	if err != nil {
		return nil, errors.New("failed to deserialize blog comment list")
	}
	return wrapper.Comments, nil
}