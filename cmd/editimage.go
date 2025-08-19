/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/ajuala/gogem/ai"

	"github.com/spf13/cobra"
)

var (
	eimageInputFile string
	eimageOutputFile string
	eimagePrintRawBytes bool
)


// editimageCmd represents the editimage command
var editimageCmd = &cobra.Command{
	Use:   "editimage",
	Short: "editimage lets you manipulate imagez with Google Gemini using text prompts",

	Run: func(cmd *cobra.Command, args []string) {
		if userPrompt == "" {
			userPrompt = readStdin()
		}

		if userPrompt == "" {
			eprint("Error: prompt is missing")
			os.Exit(1)
		}

		inFilePath := eimageInputFile

		if inFilePath == "" {
			eprint("Error: you need to specify image to edit")
			os.Exit(1)
		}

		imageData, text, err := ai.EditImage(ai.Params{
			UserPrompt: userPrompt,
			FilePath: inFilePath,
			ApiKey: apiKey,
			Model: model,
		})

		if err != nil {
			eprint(err)
			os.Exit(1)
		}

		if imageData == nil || len(imageData) == 0 {
			eprint("Error: unknown error: Gemini returned an empty data")
			os.Exit(1)
		}

		if eimageOutputFile == "" || eimageOutputFile == "-" {
			if ! eimagePrintRawBytes {
				fmt.Println(b64encode(imageData))
			} else {
				os.Stdout.Write(imageData)
			}
		} else {
			err := os.WriteFile(eimageOutputFile, imageData, 0644)
			if err != nil {
				eprint(err)
				os.Exit(1)
			}
		}

		if text != "" {
			eprint(text)
		}
	},
}

func init() {
	rootCmd.AddCommand(editimageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// editimageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// editimageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	editimageCmd.Flags().StringVarP(&eimageOutputFile, "output", "o", "", "output file for edited image. If omitted or set to\"-\", would print a base64 encoded image to the standard output (use the \"--raw\" flag to force printing raw bytes to stdout), otherwise writes the output to the specified PNG file")
	editimageCmd.Flags().BoolVarP(&eimagePrintRawBytes, "raw", "b", false, "forces sending raw bytes to the standard output")
	editimageCmd.Flags().StringVarP(&eimageInputFile, "infile", "i", "", "Image file to edit")
	editimageCmd.MarkFlagRequired("infile")
}
