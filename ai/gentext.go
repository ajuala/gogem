package ai

import (
  "context"
  "os"
  "strings"
  "encoding/json"

  "google.golang.org/genai"
)


// func GenText(p Params) (string, error) {
func GenText(userPrompt, sysPrompt, model, schema, apiKey string, temp, topK, topP *float32) (string, error) {

  client, err := NewClient(apiKey)
  if err != nil {
	  return "", err
  }



  var config *genai.GenerateContentConfig

  if sysPrompt != "" || schema != ""  || temp != nil || topK != nil || topP != nil {
	  config = &genai.GenerateContentConfig{}

	  if sysPrompt == "" {
		  config.SystemInstruction = genai.NewContentFromText(sysPrompt, genai.RoleUser)
	  }

  if schema != "" {
  var outSchema genai.Schema
	  err := json.Unmarshal([]byte(schema), &outSchema)

	  if err != nil {
		  return "", err
	  }


		  config.ResponseMIMEType = "application/json"
		  config.ResponseSchema = &outSchema

  }

  config.Temperature = temp
  config.TopP = topP
  config.TopK = topK
  }




  ctx := context.Background()

  model = strings.TrimSpace(model)
  if model == "" {
  model = "gemini-2.5-flash"
  }

  result, err := client.Models.GenerateContent(
      ctx,
	  model,
      genai.Text(userPrompt),
      config,
  )

  if err != nil {
	  return "", err
  }

  return result.Text(), nil
}

func GenTextMultiModal(prompt, sysPrompt, filepath string) (string, error) {

  ctx := context.Background()
  client, err := genai.NewClient(ctx, nil)
  if err != nil {
	  return "", err
  }

  var config *genai.GenerateContentConfig
  if sysPrompt != "" {
	  config = &genai.GenerateContentConfig{
		  SystemInstruction: genai.NewContentFromText(sysPrompt, genai.RoleUser),
	  }
  }

  fileData, err := os.ReadFile(filepath)
  if err != nil {
	  return "", err
  }

  mime := getMIME(fileData)

  parts := []*genai.Part{
	  genai.NewPartFromText(prompt),
	  &genai.Part{
		  InlineData: &genai.Blob{
			  MIMEType:mime,
			  Data: fileData,
		  },
	  },
  }

  contents := []*genai.Content{
	  genai.NewContentFromParts(parts, genai.RoleUser),
  }

  result, err := client.Models.GenerateContent(
      ctx,
      "gemini-2.5-flash",
	  contents,
      config,
  )

  if err != nil {
	  return "", err
  }

  return result.Text(), nil
}
