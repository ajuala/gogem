package ai

import (
	"context"
	"errors"
	"os"
	"io"
	"strings"
	"encoding/json"

	"google.golang.org/genai"

	"github.com/gabriel-vasile/mimetype"
)


// func GenText(p Params) (string, error) {
func GenText(userPrompt, sysPrompt, model, schema, apiKey string, uploadFilePaths, inlineFilePaths []string,  uploadStdin, inlineStdin bool, temp, topK, topP *float32) (string, error) {

	client, err := NewClient(apiKey)
	if err != nil {
		return "", err
	}


	// START: Configuration
	var config *genai.GenerateContentConfig

	if sysPrompt != "" {
		config = &genai.GenerateContentConfig{
			SystemInstruction: genai.NewContentFromText(sysPrompt, genai.RoleUser),
		}
	}

	if schema != "" {

		var outSchema genai.Schema
		err := json.Unmarshal([]byte(schema), &outSchema)

		if err != nil {
			return "", err
		}

		if config != nil {
			config.ResponseMIMEType = "application/json"
			config.ResponseSchema = &outSchema
		} else {
			config = &genai.GenerateContentConfig{
				ResponseMIMEType: "application/json",
				ResponseSchema: &outSchema,
			}
		}
	}

	if temp != nil || topK != nil || topP != nil {
		if config == nil {
			config.Temperature = temp
			config.TopP = topP
			config.TopK = topK
		} else {
			config = &genai.GenerateContentConfig{
				Temperature: temp,
				TopP: topP,
				TopK: topK,
			}
		}
	}

	// END: Configuration

	var contents []*genai.Content

	if inlineStdin && uploadStdin {
		return "", errors.New("error: cannot read from STDIN twice")
	}

	if len(uploadFilePaths) == 0 && len(inlineFilePaths) == 0 && !inlineStdin && !uploadStdin{
		contents = genai.Text(userPrompt)
	} else {
		parts := []*genai.Part{
			genai.NewPartFromText(userPrompt),
		}

		for _, pth := range uploadFilePaths {
			mtype, err := mimetype.DetectFile(pth)

			if err != nil {
				return "", err
			}

			ctx := context.Background()
			uploadFile, err := client.Files.UploadFromPath(ctx, pth, &genai.UploadFileConfig{
				MIMEType: mtype.String(),
			})

			if err != nil {
				return "", err
			}

			parts = append(parts, genai.NewPartFromURI(uploadFile.URI, uploadFile.MIMEType))
		}


		for _, pth := range inlineFilePaths {
			mtype, err := mimetype.DetectFile(pth)


			if err != nil {
				return "", err
			}

			b, err := os.ReadFile(pth)

			if err != nil {
				return "", err
			}

			parts = append(parts, genai.NewPartFromBytes(b, mtype.String()))
		}

		if inlineStdin {
			b, err := io.ReadAll(os.Stdin)

			if err != nil {
				return "", err
			}

			mtype := mimetype.Detect(b)

			parts = append(parts, genai.NewPartFromBytes(b, mtype.String()))
		}

		if uploadStdin {
			ctx := context.Background()
			mtype, f, err := upload(ctx, client.Files, os.Stdin, nil)

			if err != nil {
				return "", err
			}

			parts = append(parts, genai.NewPartFromURI(f.URI, mtype))
		}

		contents = []*genai.Content{
			genai.NewContentFromParts(parts, genai.RoleUser),
		}

	}





	ctx := context.Background()

	model = strings.TrimSpace(model)
	if model == "" {
		model = "gemini-2.5-flash"
	}

	result, err := client.Models.GenerateContent(
		ctx,
		model,
		contents,
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

func upload(ctx context.Context, files *genai.Files, r io.Reader, config *genai.UploadFileConfig) (string, *genai.File, error) {

	type Response struct {
		f *genai.File
		err error
	}

	type MimeDetect struct {
		mtype string
		err error
	}

	respChan := make( chan *Response)

	mimeChan := make(chan *MimeDetect, 1)

	// 1. Create a pipe. 
	// Writing to pw makes the data available to read from pr.
	pr, pw := io.Pipe()

	// 2. Create a TeeReader.
	// Every read from 'tee' will read from os.Stdin AND write to 'pw'.
	tee := io.TeeReader(r, pw)

	go func(){
		f, err := files.Upload(ctx, tee, config)
		respChan<- &Response{
			f: f,
			err: err,
		}
	}()

	go func(){
		mtype, err := mimetype.DetectReader(pr)


		mimeChan<- &MimeDetect{
			mtype: mtype.String(),
			err: err,
		}

		io.Copy(io.Discard, pr)
	}()

	resp := <-respChan
	close(respChan)

	mime := <-mimeChan
	close(mimeChan)

	if resp.err != nil {
		return "", nil, resp.err
	}

	var mtype string

	if mime.err != nil {
		mtype = resp.f.MIMEType
	} else {
		mtype = mime.mtype
	}

	return mtype, resp.f, nil
}
