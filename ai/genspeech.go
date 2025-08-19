package ai

import(
	"context"
	"strings"

	"google.golang.org/genai"
)

func GenSpeech(p Params) ([]byte, error) {
	client, err := NewClient(p.ApiKey)
	
	if err != nil {
		return nil, err
	}

	voice := p.Voice

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

	parts := []*genai.Part{ genai.NewPartFromText(p.UserPrompt) }

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	if p.SysPrompt != "" {
		config.SystemInstruction = genai.NewContentFromText(p.SysPrompt, genai.RoleUser)
	}

	model := strings.TrimSpace(p.Model)
	if model == "" {
		model = "gemini-2.5-flash-preview-tts"
	}

	ctx := context.Background()

	result, err := client.Models.GenerateContent(
		ctx,
		model,
		contents,
		config,
	)

	if err != nil {
		return nil, err
	}

	audioData := result.Candidates[0].Content.Parts[0].InlineData.Data

	return audioData, nil
}
