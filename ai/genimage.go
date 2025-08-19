package ai

import (
	"context"
	"strings"

	"google.golang.org/genai"
)

func GenImage(p Params) ([]byte, string, error) {

	prompt := p.UserPrompt
	apiKey := p.ApiKey
	model := strings.TrimSpace(p.Model)

	if model == "" {
		model = "gemini-2.0-flash-preview-image-generation"
	}

	client, err := NewClient(apiKey)
	if err != nil {
		return nil, "", err
	}

	config := &genai.GenerateContentConfig{
		ResponseModalities: []string{"TEXT", "IMAGE"},
	}

	if p.SysPrompt != "" {
		config.SystemInstruction = genai.NewContentFromText(p.SysPrompt, genai.RoleUser)
	}

	ctx := context.Background()

	result, err := client.Models.GenerateContent(
		ctx,
		model,
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
