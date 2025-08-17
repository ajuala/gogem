package ai

import (
    "net/http"
	"context"
	"google.golang.org/genai"
)

func NewClient(apiKey string) (*genai.Client, error) {
	ctx := context.Background()

	if apiKey == "" {
		return genai.NewClient(ctx, nil)
	} else {
		return genai.NewClient(ctx, &genai.ClientConfig{
			APIKey: apiKey,
			Backend: genai.BackendGeminiAPI,
		})
	}
}



func getMIME(fileData []byte) string {
    return http.DetectContentType(fileData)

}
