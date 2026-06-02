/*
Copyright © 2025 NAME HERE <!-- <EMAIL ADDRESS> -->

*/
package cmd

import (
	"os"
	"fmt"
	"strings"

	"github.com/ajuala/gogem/ai"
	"github.com/ajuala/gogem/utils"
	"github.com/ajuala/gogem/types"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	voice string
	speechOut string
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

		if strings.TrimSpace(speechOut) == "" {
			cmd.Usage()
			eprint("\nerror: option `--output` is required. Invalid output filename")
			os.Exit(1)
		}

		// BEGIN: Read prompt from file
		if strings.HasPrefix(userPrompt, "@") {
			fname := userPrompt[1:]
			if fname != "" {
				b, err := os.ReadFile(fname)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error: cannot read %v\n", fname)
					os.Exit(1)
				}

				if len(b) == 0 {
					eprint("error: empty file")
					os.Exit(1)
				}

				userPrompt = string(b)

			} else {
				eprint("error: no filename following the '@' symbol")
				os.Exit(1)
			}
		}

		if promptPath != "" {
			b, err := os.ReadFile(promptPath)

			if err != nil {
				fmt.Fprintf(os.Stderr, "error: cannot access file %s\n", promptPath)
				os.Exit(1)
			}

			p := strings.TrimSpace(string(b))
			if p == "" || p == "-" {
				fmt.Fprint(os.Stderr, "error: file `%s` contains no usable prompt\n", promptPath)

				os.Exit(1)
			}

			userPrompt = p
		}

		// END: Read prompt from file


		if userPrompt == "" || userPrompt  == "-" {
			userPrompt = readStdin()
		}

		if userPrompt == "" {
			eprint("Error: User prompt is empty")
			os.Exit(1)
		}

		voice = strings.TrimSpace(voice)
		speakers := viper.GetStringSlice("genspeech.speaker")


		var parsedSpeakers []types.Speaker
		if len(speakers) > 0 {
			voice = ""
			pSpeakers, err := parseSpeakers(speakers)

			if err != nil {
				eprint(err)
				os.Exit(1)
			}

			parsedSpeakers = pSpeakers
		}

		if len(speakers) > 0 && len(speakers) != 2 {
			eprint(fmt.Sprintf("invalid number of speakers. Only 2 is currently supported %d given.", len(speakers)))
			os.Exit(1)
		}

		if len(voice) != 0 && !utils.HasVoice(voice) {
			eprint("unsupported voice name: use option \"--show-voices\" to print supported voice names")
			os.Exit(1)
		}

		if len(voice) == 0 && len(speakers) == 0 {
			voice = "Kore"
		}

		temp, topK, topP := getTempTopKP()

		model := viper.GetString("genspeech.model")
		pcmData, err := ai.GenSpeech(userPrompt, sysPrompt, voice,  model, apiKey, parsedSpeakers, temp, topK, topP)
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


		if speechOut == "-" {
				os.Stdout.Write(b)
		} else {
			err := os.WriteFile(speechOut, b, 0644)

			if err != nil {
				eprint(err)
				os.Exit(1)
			}
		}

	},
}

func parseSpeakers(speakers []string) ([]types.Speaker, error) {
	result := make([]types.Speaker, 0)

	for _, s := range speakers {
		splitVal := strings.SplitAfterN(s, ":", 2)
		if len(splitVal) != 2 {
			return nil, fmt.Errorf("%s is not a valid speaker format, Run the help command for detail")
		}

		splitVal[0] = strings.TrimSpace(strings.TrimSuffix(splitVal[0], ":"))
		splitVal[1] = strings.TrimSpace(splitVal[1])

		if !utils.HasVoice(splitVal[0]) {
			return nil, fmt.Errorf("%s is not a supported voice, Run 'gogem genspeech --show-voices' for supported voices.", splitVal[0])
		}

		result = append(result, types.Speaker{
			Voice: splitVal[0],
			Name: splitVal[1],
		})

	}

	return result, nil
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
	genspeechCmd.Flags().StringVarP(&speechOut, "output", "o", "", "Output file for WAV encoded audio(required). If set to \"-\" prints to the standard output")

	genspeechCmd.Flags().StringVar(&voice, "voice", "Kore", "Voice for reading out text. Visit Google Gemini API website, or use option \"--show-voices\" to print out available voices")
	genspeechCmd.Flags().BoolVar(&printVoices, "show-voices", false, "Print names of Gemini voices. Note: the bracketted letters, (F) and (M), are not part of the names but indicate gender the voices most resembles")

	genspeechCmd.Flags().StringP("model", "m", "gemini-2.5-flash-preview-tts", "Gemini text-to-speech  model variant")
	viper.BindPFlag("genspeech.model", genspeechCmd.Flags().Lookup("model"))

	genspeechCmd.Flags().StringSlice("speaker", nil, "Voice in the format <voice>:<speaker> where <voice> is a valid Gemini voice and <speaker> is the name used to reference the voice in your prompt. This option mutually excludes with `--voice`.\nNote: exactly 2 `--speaker` options is required.")
	viper.BindPFlag("genspeech.speaker", genspeechCmd.Flags().Lookup("speaker"))

	genspeechCmd.MarkFlagsMutuallyExclusive("voice", "speaker")
}
