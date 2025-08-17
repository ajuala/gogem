package ai

import(
	"context"
	"google.golang.org/genai"
)

func GenSpeech(textPrompt, sysPrompt, voice, apiKey string) ([]byte, error) {
	client, err := NewClient(apiKey)
	
	if err != nil {
		return nil, err
	}

	config := &genai.GenerateContentConfig{
		ResponseModalities: []string{"AUDIO"},
		SpeechConfig: &genai.SpeechConfig{
			VoiceConfig: &genai.VoiceConfig{
			PrebuiltVoiceConfig: &genai.PrebuiltVoiceConfig{
				VoiceName: voice,
			},
			},
		},
	}

	parts := []*genai.Part{ genai.NewPartFromText(textPrompt) }

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	if sysPrompt != "" {
		config.SystemInstruction = genai.NewContentFromText(sysPrompt, genai.RoleUser)
	}

	ctx := context.Background()

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash-preview-tts",
		contents,
		config,
	)

	if err != nil {
		return nil, err
	}

	audioData := result.Candidates[0].Content.Parts[0].InlineData.Data

	return audioData, nil
}
