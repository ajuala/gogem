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
		model = "gemini-2.5-flash-image-preview"
	}

	client, err := NewClient(apiKey)
	if err != nil {
		return nil, "", err
	}

	var config *genai.GenerateContentConfig

	if sysPrompt != "" || temp != nil || topK != nil || topP != nil || model == "gemini-2.0-flash-preview-image-generation" {
	config = &genai.GenerateContentConfig{
		Temperature: temp,
		TopK: topK,
		TopP: topP,
	}

	if sysPrompt != "" {
		config.SystemInstruction = genai.NewContentFromText(sysPrompt, genai.RoleUser)
	}

	if model == "gemini-2.0-flash-preview-image-generation" {
		config.ResponseModalities = []string{"TEXT", "IMAGE"}
	}
	}

	sysPrompt = strings.TrimSpace(sysPrompt)


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
