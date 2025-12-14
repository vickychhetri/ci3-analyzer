/*
Copyright © 2025 Vicky Chhetri <vickychhetri4@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ci3-analyzer",
	Short: "CI3 HMVC code analyzer",
	Long: `CI3 Analyzer is a CLI tool to scan
			CodeIgniter 3 HMVC projects and generate
			HTML documentation with modules, classes, and methods.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
