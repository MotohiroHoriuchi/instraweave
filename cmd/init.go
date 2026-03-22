package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const sampleRecipe = `target: copilot
output: ""
fragments_dir: ./fragments
fragments:
  - standard/go
  - custom/my-project
`
// instrweave-recipe.yaml is the default recipe file name

const sampleStandardFragment = `# Go Coding Standards

- Follow standard Go conventions (Effective Go, Go Code Review Comments).
- Use gofmt / goimports for formatting.
- Handle errors explicitly; do not ignore returned errors.
`

const sampleCustomFragment = `# My Project Conventions

- Add your project-specific rules here.
`

func init() {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a sample recipe and fragments directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := os.Stat("instrweave-recipe.yaml"); err == nil {
				return fmt.Errorf("instrweave-recipe.yaml already exists")
			}

			if err := os.WriteFile("instrweave-recipe.yaml", []byte(sampleRecipe), 0o644); err != nil {
				return fmt.Errorf("failed to create recipe file: %w", err)
			}

			dirs := []string{"fragments/standard", "fragments/custom"}
			for _, d := range dirs {
				if err := os.MkdirAll(d, 0o755); err != nil {
					return fmt.Errorf("failed to create directory %s: %w", d, err)
				}
			}

			files := map[string]string{
				"fragments/standard/go.md":     sampleStandardFragment,
				"fragments/custom/my-project.md": sampleCustomFragment,
			}
			for path, content := range files {
				if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
					return fmt.Errorf("failed to create %s: %w", path, err)
				}
			}

			fmt.Println("Created instrweave-recipe.yaml and fragments/ directory.")
			return nil
		},
	}

	rootCmd.AddCommand(initCmd)
}
