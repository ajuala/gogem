/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"
	"os/exec"
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
	model string
	temperature float32
	topP float32
	topK float32
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gogem",
	Short: "gogem is a client for Google Gemini's API",
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
	rootCmd.PersistentFlags().StringVarP(&model, "model", "m", "", "Gemini AI model. Each command uses a different default. Make sure the model supports the task the command seeks to execute before setting this option")
	rootCmd.PersistentFlags().Float32Var(&temperature, "temp", 0, "Temperature value")
	rootCmd.PersistentFlags().Lookup("temp").NoOptDefVal = ""
	rootCmd.PersistentFlags().Float32Var(&topP, "topp", 0, "TopP value")
	rootCmd.PersistentFlags().Lookup("topp").NoOptDefVal = ""
	rootCmd.PersistentFlags().Float32Var(&topK, "topk", 0, "TopK value")
	rootCmd.PersistentFlags().Lookup("topk").NoOptDefVal = ""

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}


func eprint(msg any) (int, error) {
	return fmt.Fprintln(os.Stderr, msg)
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


// getTempTopKP returns pointers to teperature, topK, and topP if flags are set. Returns nils otherwise
func getTempTopKP() (tempVal *float32, topKVal *float32, topPVal *float32) {

	if rootCmd.Flags().Lookup("temp").Changed {
		tempVal = &temperature
	}

	if rootCmd.Flags().Lookup("topk").Changed {
		topKVal = &topK
	}

	if rootCmd.Flags().Lookup("topp").Changed {
		topPVal = &topP
	}

	return tempVal, topKVal, topPVal
}

func vimEditor(initText string, extension string, neovim bool) (string, error) {

    tmpFile, err := os.CreateTemp("", "tempfile_*." + strings.TrimLeft(extension, ". "))
    if err != nil {
		return "", err
    }
    defer os.Remove(tmpFile.Name()) // clean up

    // Write initial content or leave empty
    initialContent := []byte(initText)
    if _, err := tmpFile.Write(initialContent); err != nil {
		return "", err
    }
    tmpFile.Close() // close to flush content to disk

    // Launch Vim editor on the temp file
	editor := "vim"

	if neovim {
		editor = "nvim"
	}

    cmd := exec.Command(editor, tmpFile.Name())
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    if err := cmd.Run(); err != nil {
		return "", err
    }

    // Read the edited content back into memory using os.ReadFile
    editedContent, err := os.ReadFile(tmpFile.Name())
    if err != nil {
		return "", err
    }

	return string(editedContent), err
}
