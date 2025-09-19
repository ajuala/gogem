package ai

import (
	"context"
	"strings"
	"strconv"
	"regexp"
	"fmt"

	"google.golang.org/genai"
)

// func GenImage(p Params) ([]byte, string, error) {
func GenImage(userPrompt, sysPrompt, model, apiKey string, temp, topK, topP *float32) ([]byte, string, error) {

	model = strings.TrimSpace(model)

	if model == "" {
		model = "gemini-2.5-flash-image-preview"
	}

	variantPtn := regexp.MustCompile("^gemini-(?:live-)?([1-9]+)\\.(\\d+)")
	variantVerNums := variantPtn.FindStringSubmatch(model)

	if variantVerNums == nil {
		return nil, "", fmt.Errorf("model %s not supported", model)
	}

	majorVer, _ := strconv.Atoi(variantVerNums[1])
	minorVer, _ := strconv.Atoi(variantVerNums[2])

	client, err := NewClient(apiKey)
	if err != nil {
		return nil, "", err
	}

	var config *genai.GenerateContentConfig

	if sysPrompt != "" || temp != nil || topK != nil || topP != nil || majorVer < 2 || (majorVer == 2 && minorVer < 5) {
	config = &genai.GenerateContentConfig{
		Temperature: temp,
		TopK: topK,
		TopP: topP,
	}

	if sysPrompt != "" {
		config.SystemInstruction = genai.NewContentFromText(sysPrompt, genai.RoleUser)
	}

	if majorVer == 1 || minorVer < 5 {
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
