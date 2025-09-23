/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// genembedCmd represents the genembed command
var genembedCmd = &cobra.Command{
	Use:   "genembed",
	Short: "Generate embeddings",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("genembed called")
	},
}

func init() {
	rootCmd.AddCommand(genembedCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genembedCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	genembedCmd.Flags().StringP("model", "m", "gemini-embedding-001", "Embedding model")
	genembedCmd.Flags().String("data-list", "", "JSON encoded list of strings")
	genembedCmd.Flags().StringSliceP("data", "d", nil, "Input data")
	genembedCmd.Flags().String("task-type", "", "Task type")
	genembedCmd.Flags().IntP("size", "s", 0, "Control size of output embedding vector")
	genembedCmd.Flags().Lookup("size").NoDefOptVal = ""
}
