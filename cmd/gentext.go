/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ajuala/gogem/ai"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	schema string
	schemaPath string
	gtextOutFile string
)

// gentextCmd represents the gentext command
var gentextCmd = &cobra.Command{
	Use:   "gentext",
	Short: "gentext makes API calls to Google Gemini for text based response. Use this to generate text output",
	Aliases: []string{"textgen"},
	Run: func(cmd *cobra.Command, args []string) {

		// Read prompt from file
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

		// Read from stdin
		if userPrompt == "" || userPrompt == "-" {
			userPrompt = readStdin()
		}

		if userPrompt == "" {
			eprint("Error: cannot work with an empty prompt")
			os.Exit(1)
		}


		uploadFromStdin := false
		inlineFromStdin := false

		uploadFiles := viper.GetStringSlice("gentext.upload")
		uploadFilesFiltered := make([]string, 0)

		for _, pth := range uploadFiles {

			if uploadFromStdin {
				eprint("error: multiple stdin read from option '--upload'.")
				os.Exit(1)
			}

			if pth == "-" {
				uploadFromStdin = true
			} else {
				uploadFilesFiltered = append(uploadFilesFiltered, pth)
			}

		}

		inlineFiles := viper.GetStringSlice("gentext.inline")
		inlineFilesFiltered := make([]string, 0)


		for _, pth := range inlineFiles {
			if inlineFromStdin || uploadFromStdin {
				if uploadFromStdin {
					eprint("error: cannot read multiple times from stdin. Option '--upload' already set to '-'")
				} else {
					eprint("error: cannot read multiple times from stdin. Option '--inline' set to '-' multiple times")
				}

				os.Exit(1)
			}

			if pth == "-" {
				inlineFromStdin = true
			} else {
				inlineFilesFiltered = append(inlineFilesFiltered, pth)
			}

		}

		if (userPrompt == "" || userPrompt == "-") && (uploadFromStdin || inlineFromStdin) {
			eprint("Only one of '--prompt', '--upload', or '--inline' can be read from stdin")
			os.Exit(1)
		}

		// Check fir file size
		maxFileSize := int64(20971520)
		totalFileSize := int64(0)

		for _, pth := range inlineFilesFiltered {
			stat, err := os.Stat(pth)

			if err != nil {
				eprint(err)
				os.Exit(1)
			}

			totalFileSize += stat.Size()

			if totalFileSize >= maxFileSize {
				eprint("error: inline data total size exceeds max file size")

				os.Exit(1)
			}
		}


		schemaData, err := getSchema()

		if err != nil {
			eprint(err)
			os.Exit(1)
		}

		temp, topK, topP := getTempTopKP()

		model := viper.GetString("gentext.model")
		result, err := ai.GenText(userPrompt, sysPrompt, model, schemaData, apiKey, uploadFilesFiltered, inlineFilesFiltered, false, false, temp, topK, topP)
		if err != nil {
			eprint(err)
			os.Exit(1)
		}

		if gtextOutFile == "" {
			fmt.Println(result)
		} else {
			err := os.WriteFile(gtextOutFile, []byte(result), 0644)

			if err != nil {
				eprint(err)
				os.Exit(1)
			}
		}

	},
}

func getSchema() (string, error) {
	data := strings.TrimSpace(schema)
	if data == "" {
		if schema!= "" {
			b, err := os.ReadFile(schemaPath)
			if err != nil {
				return "", err
			}

			return string(b), nil
		} else {
			return "", nil
		}
	} else {
		return data, nil
	}
}


func init() {
	rootCmd.AddCommand(gentextCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gentextCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gentextCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	gentextCmd.Flags().StringVarP(&gtextOutFile, "output", "o", "", "Output file. Prints to stanard output by default or if set to\"-\".")
	gentextCmd.Flags().StringVar(&schema, "schema", "", "JSON schema, to constrain output structure. Use either this option or \"--schema-path\", not both")
	gentextCmd.Flags().StringVar(&schemaPath, "schema-path", "", "Path to JSON schema file, to constrain output structure. Use either this option or \"--schema\", not both")
	gentextCmd.MarkFlagsMutuallyExclusive("schema", "schema-path")

	gentextCmd.Flags().StringP("model", "m", "gemini-2.5-flash", "Gemini AI model. Each command uses a different default. Make sure the model supports the task the command seeks to execute before setting this option")
	viper.BindPFlag("gentext.model", gentextCmd.Flags().Lookup("model"))


	gentextCmd.Flags().StringSliceP("upload", "u", nil, "File to upload to Gemini File API. Use '-' to read from stdin.")
	viper.BindPFlag("gentext.upload", gentextCmd.Flags().Lookup("upload"))


	gentextCmd.Flags().StringSliceP("inline", "i", nil, "File to embed inline with your request. Use '-' to read from stdin.")
	viper.BindPFlag("gentext.inline", gentextCmd.Flags().Lookup("inline"))
}
