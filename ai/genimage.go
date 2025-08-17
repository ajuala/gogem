package ai

import (
	"context"
	"google.golang.org/genai"
)

func GenImage(prompt, apiKey string) ([]byte, string, error) {

	client, err := NewClient(apiKey)
	if err != nil {
		return nil, "", err
	}

	config := &genai.GenerateContentConfig{
		ResponseModalities: []string{"TEXT", "IMAGE"},
	}

	ctx := context.Background()

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash-preview-image-generation",
		genai.Text(prompt),
		config,
	)

	if err != nil {
		return nil, "", err
	}

	var (
		t string
		b []byte
	)

for _, part := range result.Candidates[0].Content.Parts {
      if part.Text != "" {
		  t = part.Text
      } else if part.InlineData != nil {
          b = part.InlineData.Data
      }
  }

	return b, t, nil


}
