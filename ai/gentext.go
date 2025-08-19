package ai

import (
  "context"
  "os"
  "strings"
  "encoding/json"

  "google.golang.org/genai"
)


func GenText(p Params) (string, error) {

  client, err := NewClient(p.ApiKey)
  if err != nil {
	  return "", err
  }

  jsonOut := false
  schemaData := strings.TrimSpace(p.SchemaData)
  schemaPath := p.SchemaPath
  var schemaBytes []byte
  if schemaData != "" {
	  schemaBytes = []byte(schemaData)
  } else if schemaPath != "" {
	  b, err := os.ReadFile(schemaPath)
	  if err != nil {
		  return "", err
	  }
	  jsonOut = true
	  schemaBytes = b
  }



  var config *genai.GenerateContentConfig

  sysPrompt := p.SysPrompt
  if sysPrompt != "" {
	  config = &genai.GenerateContentConfig{
		  SystemInstruction: genai.NewContentFromText(sysPrompt, genai.RoleUser),
	  }
  }


  var outSchema genai.Schema
  if jsonOut {
	  err := json.Unmarshal(schemaBytes, &outSchema)
	  if err != nil {
		  return "", err
	  }

	  if config == nil {
		  config = &genai.GenerateContentConfig{
			  ResponseMIMEType: "application/json",
			  ResponseSchema: &outSchema,
		  }
	  } else {
		  config.ResponseMIMEType = "application/json"
		  config.ResponseSchema = &outSchema
	  }
  }

  ctx := context.Background()

  model := "gemini-2.5-flash"
  modelTrimmed := strings.TrimSpace(p.Model)
  if modelTrimmed != "" {
	  model = modelTrimmed
  }

  result, err := client.Models.GenerateContent(
      ctx,
	  model,
      genai.Text(p.UserPrompt),
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
