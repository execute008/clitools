package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "clitools",
	Short: "A collection of CLI tools for various tasks",
	Long: `clitools is a collection of command-line utilities designed to help with
common development and optimization tasks. Currently includes image manipulation
tools for web optimization.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add subcommands here
	rootCmd.AddCommand(imageCmd)
}
