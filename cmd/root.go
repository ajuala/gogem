/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"
	"fmt"
	"bufio"
	"io"
	"strings"
	"encoding/base64"

	"github.com/spf13/cobra"

)

var (
	sysPrompt string
	userPrompt string
	apiKey string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ggemini",
	Short: "ggemini is a client for Google Gemini's API",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&sysPrompt, "sys", "s", "", "System prompt")
	rootCmd.PersistentFlags().StringVarP(&userPrompt, "prompt", "p", "", "Text prompt. (Default: reads from STDIN.)")
	rootCmd.PersistentFlags().StringVarP(&apiKey, "apikey", "k", "", "Google Gemini API key. (Default: uses the environment variable GEMINI_API_KEY)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


func eprint(msg any) (int, error) {
	return fmt.Fprint(os.Stderr, msg)
}

func b64encode(data []byte) string {
	// base64.StdEncoding provides the standard Base64 encoding.
	// EncodeToString is a convenient method to encode a byte slice
	// directly into a string.
	return base64.StdEncoding.EncodeToString(data)
}


func readStdin() string {
	reader := bufio.NewReader(os.Stdin)
	var result string

	for {
		line, err := reader.ReadString('\n')

		if err == io.EOF {
			if len(line) > 0 {
				result += line
			}
			
			break
		}

		if err != nil {
			return ""
		}

		result += line
	}

	return strings.TrimSpace(result)
}




