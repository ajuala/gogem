package ai

import (
	"context"
	"os"
	"strings"

	"google.golang.org/genai"
)

func EditImage(p Params) ([]byte, string, error) {

	prompt := p.UserPrompt
	imagePath := p.FilePath
	sysPrompt := strings.TrimSpace(p.SysPrompt)
	apiKey := p.ApiKey
	model := strings.TrimSpace(p.Model)

	if model == "" {
		model = "gemini-2.0-flash-preview-image-generation"
	}

	client, err := NewClient(apiKey)
	if err != nil {
		return nil, "", err
	}

	imgData, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, "", err
	}

	mime := getMIME(imgData)

	parts := []*genai.Part{
		genai.NewPartFromText(prompt),
		&genai.Part{
			InlineData: &genai.Blob{
				MIMEType: mime,
				Data:     imgData,
			},
		},
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	config := &genai.GenerateContentConfig{
		ResponseModalities: []string{"TEXT", "IMAGE"},
	}

	if sysPrompt != "" {
		config.SystemInstruction = genai.NewContentFromText(sysPrompt, genai.RoleUser)
	}

	ctx := context.Background()
	result, err := client.Models.GenerateContent(
		ctx,
		model,
		contents,
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
