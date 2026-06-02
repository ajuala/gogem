package ai

import (
	"context"
	"google.golang.org/genai"
)

func Imagen(userPrompt, model, apiKey, res, aspectRatio, imageSize string, personGen genai.PersonGeneration, numOfImages int32)(*genai.GenerateImagesResponse, error) {

	var client *genai.Client

	ctx := context.Background()

	if apiKey == "" {
		c, err := genai.NewClient(ctx, nil)

		if err != nil {
			return nil, err
		}

		client = c
	} else {

		c, err := genai.NewClient(ctx, &genai.ClientConfig{
			APIKey: apiKey,
		})

		if err != nil {
			return nil, err
		}

		client = c
	}

	config := &genai.GenerateImagesConfig{
		NumberOfImages: numOfImages,
	}

	if aspectRatio != "" {
		config.AspectRatio = aspectRatio
	}

	if imageSize != "" {
		config.ImageSize = imageSize
	}

	if personGen != "" {
		config.PersonGeneration = personGen
	}



	return client.Models.GenerateImages(
		ctx,
		model,
		userPrompt,
		config,
	)
}
