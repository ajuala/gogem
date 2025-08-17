package ai

import (
	"context"
	"os"

	"google.golang.org/genai"
)

func EditImage(prompt, imagePath, apiKey string) ([]byte, string, error) {

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

	ctx := context.Background()
	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash-preview-image-generation",
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
