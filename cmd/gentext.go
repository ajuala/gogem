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
	schema string
)

// gentextCmd represents the gentext command
var gentextCmd = &cobra.Command{
	Use:   "gentext",
	Short: "gentext makes API calls to Google Gemini for text based response. Use this to generate text output",
	Run: func(cmd *cobra.Command, args []string) {

		// Read from stdin
		if userPrompt == "" || userPrompt == "-" {
			userPrompt = readStdin()
		}

		if userPrompt == "" {
			eprint("Error: cannot work with an empty prompt")
			os.Exit(1)
		}


		result, err := ai.GenText(ai.Params{
			UserPrompt: userPrompt,
			SysPrompt: sysPrompt,
			SchemaPath: schema,
			ApiKey: apiKey,
			Model: model,
		})
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

var (
	gtextOutFile string
)

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
	gentextCmd.Flags().StringVar(&schema, "schema", "", "Path to JSON schema file, to constrain output structure")
}
