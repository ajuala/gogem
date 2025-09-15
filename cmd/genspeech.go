/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"
	"fmt"
	"strings"

	"github.com/ajuala/gogem/ai"
	"github.com/ajuala/gogem/utils"

	"github.com/spf13/cobra"
)

var (
	voice string
	speechOut string
	outputRawWav bool
	printVoices bool
)

// genspeechCmd represents the genspeech command
var genspeechCmd = &cobra.Command{
	Use:   "genspeech",
	Short: "Generate speech using Google Gemini text-to-speech model",

	Aliases: []string{"speechgen"},

Run: func(cmd *cobra.Command, args []string) {

	if printVoices {
		fmt.Println("Supported voices:")
		utils.PrintVoices()
		return
	}

	if userPrompt == "" || userPrompt  == "-" {
		userPrompt = readStdin()
	}

	if userPrompt == "" {
		eprint("Error: User prompt is empty")
		os.Exit(1)
	}

	voice = strings.TrimSpace(voice)
	if ! utils.HasVoice(voice) {
		eprint("unsupported voice name: use option \"--show-voices\" to print supported voice names")
		os.Exit(1)
	}

	temp, topK, topP := getTempTopKP()

	pcmData, err := ai.GenSpeech(userPrompt, sysPrompt, voice, model, apiKey, temp, topK, topP)
	if err != nil {
		eprint(err)
		os.Exit(1)
	}

	b, err := utils.ConvertPCMToWav(
		pcmData,
		1,
		24000,
		16,
	)


	if err != nil {
		eprint(err)
		os.Exit(1)
	}

	if speechOut == "" || speechOut == "-" {
		if outputRawWav {
			os.Stdout.Write(b)
		} else {
			encoded := b64encode(b)
			fmt.Println(encoded)
		}
	} else {
		err := os.WriteFile(speechOut, b, 0644)

		if err != nil {
			eprint(err)
			os.Exit(1)
		}
	}

},
}

func init() {
	rootCmd.AddCommand(genspeechCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genspeechCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genspeechCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	genspeechCmd.Flags().StringVarP(&speechOut, "output", "o", "", "Output file for WAV encoded audio. If omitted or set to \"-\", defaults to printing Base64 encoded output to stdout")
	genspeechCmd.Flags().StringVar(&voice, "voice", "Kore", "Voice for reading out text. Visit Google Gemini API website, or use option \"--show-voices\" to print out available voices")
	genspeechCmd.Flags().BoolVarP(&outputRawWav, "raw", "b", false, "Forces output of raw bytes to standard output. Has no effect when \"--output\" is set to values other than \"-\" or empty string")
	genspeechCmd.Flags().BoolVar(&printVoices, "show-voices", false, "Print names of Gemini voices. Note: the bracketted letters, (F) and (M), are not part of the names but indicate gender the voices most resembles")
}
