package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/MotohiroHoriuchi/instraweave/internal/fragment"
	"github.com/MotohiroHoriuchi/instraweave/internal/recipe"
	"github.com/spf13/cobra"
)

func init() {
	var recipePath string
	var dryRun bool

	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate instructions file from a recipe",
		RunE: func(cmd *cobra.Command, args []string) error {
			r, err := recipe.Load(recipePath)
			if err != nil {
				return err
			}

			content, err := fragment.Compose(r.FragmentsDir, r.Fragments)
			if err != nil {
				return err
			}

			if dryRun {
				fmt.Print(content)
				return nil
			}

			if dir := filepath.Dir(r.Output); dir != "." {
				if err := os.MkdirAll(dir, 0o755); err != nil {
					return fmt.Errorf("failed to create output directory: %w", err)
				}
			}

			if err := os.WriteFile(r.Output, []byte(content), 0o644); err != nil {
				return fmt.Errorf("failed to write output: %w", err)
			}

			fmt.Printf("Generated %s\n", r.Output)
			return nil
		},
	}

	generateCmd.Flags().StringVarP(&recipePath, "recipe", "r", "./instrweave-recipe.yaml", "path to recipe file")
	generateCmd.Flags().BoolVar(&dryRun, "dry-run", false, "print to stdout instead of writing file")
	rootCmd.AddCommand(generateCmd)
}
