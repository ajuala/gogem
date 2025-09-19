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

		// Read from stdin
		if userPrompt == "" || userPrompt == "-" {
			userPrompt = readStdin()
		}

		if userPrompt == "" {
			eprint("Error: cannot work with an empty prompt")
			os.Exit(1)
		}


		schemaData, err := getSchema()

		if err != nil {
			eprint(err)
			os.Exit(1)
		}

		temp, topK, topP := getTempTopKP()

		model := viper.GetString("gentext.model")
		result, err := ai.GenText(userPrompt, sysPrompt, model, schemaData, apiKey, temp, topK, topP)
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

	gentextCmd.Flags().StringP("model", "m", "gemini-2.5-flash", "Gemini AI model. Each command uses a different default. Make sure the model supports the task the command seeks to execute before setting this option")
	viper.BindPFlag("gentext.model", gentextCmd.Flags().Lookup("model"))
}
