package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "instraweave",
	Short: "instraweave: build AI agent instructions from fragments",
}

func Execute() error {
	return rootCmd.Execute()
}
