package ai

import(
	"context"
	"strings"

	"github.com/ajuala/gogem/types"

	"google.golang.org/genai"
)

// func GenSpeech(p Params) ([]byte, error) {
func GenSpeech(userPrompt, sysPrompt, voice, model, apiKey string, speakers []types.Speaker, temp, topK, topP *float32) ([]byte, error) {



	client, err := NewClient(apiKey)

	if err != nil {
		return nil, err
	}


	config := &genai.GenerateContentConfig{
		ResponseModalities: []string{"AUDIO"},
		Temperature: temp,
		TopK: topK,
		TopP: topP,
	}

	if len(speakers) == 0 {
		config.SpeechConfig = &genai.SpeechConfig{
			VoiceConfig: &genai.VoiceConfig{
				PrebuiltVoiceConfig: &genai.PrebuiltVoiceConfig{
					VoiceName: voice,
				},
			},
		}
	} else {
		multiSpeakerVoiceConfig := genai.MultiSpeakerVoiceConfig{}

		for _, speaker := range speakers {
			multiSpeakerVoiceConfig.SpeakerVoiceConfigs = append(multiSpeakerVoiceConfig.SpeakerVoiceConfigs, &genai.SpeakerVoiceConfig{
				Speaker: speaker.Name,
				VoiceConfig: &genai.VoiceConfig{
					PrebuiltVoiceConfig: &genai.PrebuiltVoiceConfig{
						VoiceName: speaker.Voice,
					},
				},
			})

		}

		config.SpeechConfig = &genai.SpeechConfig{
			MultiSpeakerVoiceConfig: &multiSpeakerVoiceConfig,
		}

	}



	parts := []*genai.Part{ genai.NewPartFromText(userPrompt) }

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	if sysPrompt != "" {
		config.SystemInstruction = genai.NewContentFromText(sysPrompt, genai.RoleUser)
	}

	model = strings.TrimSpace(model)
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
