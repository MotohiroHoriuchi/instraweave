package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "instrweave",
	Short: "instrweave: build AI agent instructions from fragments",
}

func Execute() error {
	return rootCmd.Execute()
}
