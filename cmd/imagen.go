/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"google.golang.org/genai"
)

type GenPersonOption genai.PersonGeneration

func (o *GenPersonOption) String() string {
	return string(*o)
}

func (o *GenPersonOption) Set(s string) error {
	v := strings.ToUpper(s)
	switch v {
	case "ALLOW_ALL", "ALLOW_ADULT", "DONT_ALLOW":
		*o = GenPersonOption(v)
		return nil

	default:
		return fmt.Errorf("error: `%s` is not a valid option. Valid options include `DONT_ALLOW`, `ALLOW_ADULT`, `ALLOW_ALL`", s)
	}
}

func (o *GenPersonOption) Type() string {
	return "GenPersonOption"
}

// imagenCmd represents the imagen command
var imagenCmd = &cobra.Command{
	Use:   "imagen",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("imagen called")
	},
}

func init() {
	rootCmd.AddCommand(imagenCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// imagenCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// imagenCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
