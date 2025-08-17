/*
Copyright © 2025 NAME HERE <!-- <EMAIL ADDRESS> -->

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/ajuala/ggemini/ai"

	"github.com/spf13/cobra"
)

var (
	gimageOutputFile string
	gimagePrintRawBytes bool
)

// genimageCmd represents the genimage command
var genimageCmd = &cobra.Command{
	Use:   "genimage",
	Short: "genimage makes API calls to Google Gemini for image based response. Use this to generate images",

	Run: func(cmd *cobra.Command, args []string) {
		if userPrompt == "" {
			userPrompt = readStdin()
		}

		if userPrompt == "" {
			eprint("Error: prompt is missing")
			os.Exit(1)
		}

		imageData, text, err := ai.GenImage(userPrompt, apiKey)

		if err != nil {
			eprint(err)
			os.Exit(1)
		}

		if imageData == nil || len(imageData) == 0 {
			eprint("Error: unknown error: Gemini returned an empty data")
			os.Exit(1)
		}

		if gimageOutputFile == "" || gimageOutputFile == "-" {
			if ! gimagePrintRawBytes {
				fmt.Println(b64encode(imageData))
			} else {
				os.Stdout.Write(imageData)
			}
		} else {
			err := os.WriteFile(gimageOutputFile, imageData, 0644)
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
	rootCmd.AddCommand(genimageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genimageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genimageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	genimageCmd.Flags().StringVarP(&gimageOutputFile, "output", "o", "", "output file for generatrd image. If omitted or set to\"-\", would print a base64 encoded image to the standard output (use the \"--raw\" flag to force printing raw bytes to stdout), otherwise writes the output to the specified PNG file")
	genimageCmd.Flags().BoolVarP(&gimagePrintRawBytes, "raw", "b", false, "forces sending raw bytes to the standard output")
}
