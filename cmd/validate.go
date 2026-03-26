package cmd

import (
	"fmt"
	"os"

	"github.com/MotohiroHoriuchi/instraweave/internal/validation"
	"github.com/spf13/cobra"
)

func init() {
	var recipePath string
	var fragmentsDir string

	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate recipe and fragments for correctness",
		RunE: func(cmd *cobra.Command, args []string) error {
			result := validation.Validate(recipePath, fragmentsDir)
			cwd, _ := os.Getwd()

			for _, issue := range result.Issues {
				switch issue.Level {
				case validation.LevelOK:
					fmt.Printf("✓ %s\n", issue.Message)
				case validation.LevelError:
					fmt.Printf("✗ %s\n", issue.Message)
					if issue.Detail != "" {
						fmt.Printf("  → expected: %s\n", relPath(cwd, issue.Detail))
					}
				case validation.LevelWarning:
					fmt.Printf("⚠ %s\n", issue.Message)
				case validation.LevelInfo:
					fmt.Printf("ℹ %s\n", issue.Message)
				}
			}

			fmt.Println()
			if result.HasErrors() {
				fmt.Printf("validation failed: %d error(s), %d warning(s)\n",
					result.ErrorCount(), result.WarningCount())
				os.Exit(1)
			}
			if result.WarningCount() > 0 {
				fmt.Printf("validation passed with %d warning(s)\n", result.WarningCount())
			} else {
				fmt.Println("validation passed")
			}
			return nil
		},
	}

	validateCmd.Flags().StringVarP(&recipePath, "recipe", "r", "./instraweave-recipe.yaml", "path to recipe file")
	validateCmd.Flags().StringVarP(&fragmentsDir, "dir", "d", "./fragments", "path to fragments directory")
	rootCmd.AddCommand(validateCmd)
}
