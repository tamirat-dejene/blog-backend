package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"g6/blog-api/Delivery/dto"
	"strings"

	"google.golang.org/genai"
)

type GeminiConfig struct {
	APIKey    string
	ModelName string
}

func (c *GeminiConfig) GenerateWithGemini(ctx context.Context, topic string, keywords []string) (string, error) {
	// Check if the API key is available.
	if c.APIKey == "" {
		return "", fmt.Errorf("API_KEY environment variable is not set")
	}

	// Create a new GenAI client with the provided API key.
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: c.APIKey,
	})

	if err != nil {
		return "", fmt.Errorf("failed to create GenAI client: %w", err)
	}

	// Prepare the prompt using the topic and keywords.
	keywordStr := strings.Join(keywords, ", ")
	prompt := fmt.Sprintf(`
Generate a blog post about "%s" using the following keywords: %s.

The blog post should include:
- A clear and engaging **title**
- An **introduction**
- A multi-paragraph **body** (at least 300–500 words)
- A **conclusion**

Then at the end, provide:
- 3 suggested alternative titles
- 3 related blog post ideas

Return your response **strictly in the following JSON format** and **don't use markdown or any other formatting, just plain JSON**:

{
  "title": "...",
  "introduction": "...",
  "body": "...",
  "conclusion": "...",
  "suggested_titles": ["...", "...", "..."],
  "related_ideas": ["...", "...", "..."]
}
`, topic, keywordStr)

	resp, err := client.Models.GenerateContent(
		ctx,
		c.ModelName,
		genai.Text(prompt),
		nil,
	)

	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	return resp.Text(), nil
}

func (c *GeminiConfig) ParseGeneratedContent(content string, output *dto.AIBlogPostResponse) error {
	// Remove ```json``` and ``` at the start and end of the content.
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)
	
	// Ensure the content is valid JSON.
	if !strings.HasPrefix(content, "{") || !strings.HasSuffix(content, "}") {
		return fmt.Errorf("generated content is not valid JSON")
	}

	// Parse the JSON content into the output structure.
	if err := json.Unmarshal([]byte(content), output); err != nil {
		return fmt.Errorf("failed to parse generated content: %w", err)
	}

	// Ensure all required fields are populated.
	if output.Title == "" || output.Introduction == "" || output.Body == "" || output.Conclusion == "" {
		return fmt.Errorf("generated content is missing required fields")
	}

	return nil
}
