package ai

import (
	"context"
	"strings"

	"google.golang.org/genai"
)

// func GenImage(p Params) ([]byte, string, error) {
func GenImage(userPrompt, sysPrompt, model, apiKey string, temp, topK, topP *float32) ([]byte, string, error) {

	model = strings.TrimSpace(model)

	if model == "" {
		model = "gemini-2.0-flash-preview-image-generation"
	}

	client, err := NewClient(apiKey)
	if err != nil {
		return nil, "", err
	}

	config := &genai.GenerateContentConfig{
		ResponseModalities: []string{"TEXT", "IMAGE"},
		Temperature: temp,
		TopK: topK,
		TopP: topP,
	}

	sysPrompt = strings.TrimSpace(sysPrompt)

	if sysPrompt != "" {
		config.SystemInstruction = genai.NewContentFromText(sysPrompt, genai.RoleUser)
	}

	ctx := context.Background()

	result, err := client.Models.GenerateContent(
		ctx,
		model,
		genai.Text(userPrompt),
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
