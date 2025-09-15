package ai

import (
	"errors"
	"strings"
	"context"
	"encoding/json"

	"google.golang.org/genai"
)



func createHistData(histJSON string) ([]*genai.Content, error) {
	var chatHistory []*genai.Content

	histJSON = strings.TrimSpace(histJSON)
	if histJSON == "" {
		return nil, nil
	}

	err := json.Unmarshal([]byte(histJSON), &chatHistory)

	if err != nil {
		return nil, errors.New("invalid history JSON data")
	}

	return chatHistory, err
}

func CreateChat(sysPrompt string, histJSON string, model string, apiKey string, temp, topK, topP *float32) (*genai.Client, *genai.Chat, error) {
	client, err := NewClient(apiKey)

	if err != nil {
		return nil, nil, err
	}

	hist, err := createHistData(histJSON)
	if err != nil {
		return nil, nil, err
	}

	model = strings.TrimSpace(model)

	if model == "" {
		model = "gemini-2.5-flash"
	}

	var config *genai.GenerateContentConfig

	if sysPrompt != "" || temp != nil || topK != nil || topP != nil {
		config = &genai.GenerateContentConfig{
			Temperature: temp,
			TopK: topK,
			TopP: topP,
		}

		if sysPrompt != "" {
			config.SystemInstruction = genai.NewContentFromText(sysPrompt, genai.RoleUser)
		}
	}

	ctx := context.Background()
	chat, err := client.Chats.Create(ctx, model, config, hist)

	if err != nil {
		return nil, nil, err
	}

	return client, chat, nil
}
