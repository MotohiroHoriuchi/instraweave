package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MotohiroHoriuchi/instraweave/internal/diff"
	"github.com/MotohiroHoriuchi/instraweave/internal/fragment"
	"github.com/MotohiroHoriuchi/instraweave/internal/recipe"
	"github.com/MotohiroHoriuchi/instraweave/internal/validation"
	"github.com/spf13/cobra"
)

func init() {
	var recipePath string
	var dryRun bool
	var showDiff bool
	var skipValidation bool

	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate instructions file from a recipe",
		RunE: func(cmd *cobra.Command, args []string) error {
			if dryRun && showDiff {
				return fmt.Errorf("--dry-run and --diff are mutually exclusive")
			}

			if !skipValidation {
				result := validation.Validate(recipePath, "./fragments")
				if result.HasErrors() {
					cwd, _ := os.Getwd()
					for _, issue := range result.Issues {
						switch issue.Level {
						case validation.LevelError:
							fmt.Printf("✗ %s\n", issue.Message)
							if issue.Detail != "" {
								fmt.Printf("  → expected: %s\n", relPath(cwd, issue.Detail))
							}
						case validation.LevelWarning:
							fmt.Printf("⚠ %s\n", issue.Message)
						}
					}
					return fmt.Errorf("validation failed: %d error(s); run 'instraweave validate' for details or use --skip-validation to bypass",
						result.ErrorCount())
				}
			}

			r, err := recipe.Load(recipePath)
			if err != nil {
				return err
			}

			entries := make([]fragment.Entry, len(r.Fragments))
			for i, f := range r.Fragments {
				entries[i] = fragment.Entry{Name: f.Name, Dir: f.FragmentsDir}
			}

			content, err := fragment.ComposeEntries(entries)
			if err != nil {
				return err
			}

			if dryRun {
				printDryRun(r)
				return nil
			}

			if showDiff {
				return printDiff(r.Output, content)
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

	generateCmd.Flags().StringVarP(&recipePath, "recipe", "r", "./instraweave-recipe.yaml", "path to recipe file")
	generateCmd.Flags().BoolVar(&dryRun, "dry-run", false, "show inheritance chain and resolved fragments instead of writing file")
	generateCmd.Flags().BoolVar(&showDiff, "diff", false, "show diff between current file and generated content")
	generateCmd.Flags().BoolVar(&skipValidation, "skip-validation", false, "skip validation before generating")
	rootCmd.AddCommand(generateCmd)
}

func printDiff(outputPath, newContent string) error {
	cwd, _ := os.Getwd()
	oldName := outputPath + " (current)"
	newName := outputPath + " (generated)"

	existingData, err := os.ReadFile(outputPath)
	var oldContent string
	if err != nil {
		// File does not exist: show full content as new
		oldName = "/dev/null"
		oldContent = ""
	} else {
		oldContent = string(existingData)
	}

	d := diff.Unified(oldName, newName, oldContent, newContent)
	if d == "" {
		relOut := relPath(cwd, outputPath)
		fmt.Printf("No changes: %s is already up to date.\n", relOut)
		return nil
	}
	fmt.Print(d)
	return nil
}

func printDryRun(r *recipe.ResolvedRecipe) {
	cwd, _ := os.Getwd()

	fmt.Println("Inheritance chain:")
	for i, p := range r.Chain {
		rel := relPath(cwd, p)
		label := ""
		if i == 0 {
			label = "  (root)"
		}
		if i == len(r.Chain)-1 {
			label = "  (current)"
		}
		prefix := "  " + strings.Repeat("     ", i)
		if i == 0 {
			fmt.Printf("%s%s%s\n", prefix, rel, label)
		} else {
			fmt.Printf("%s└─ %s%s\n", prefix, rel, label)
		}
	}

	fmt.Println()
	fmt.Println("Resolved fragments:")

	nameWidth := 0
	for _, f := range r.Fragments {
		if len(f.Name) > nameWidth {
			nameWidth = len(f.Name)
		}
	}

	for _, f := range r.Fragments {
		fragFile := filepath.Join(f.FragmentsDir, f.Name+".md")
		relFragFile := relPath(cwd, fragFile)

		opTag := ""
		if f.Op == "add" || f.Op == "override" {
			opTag = fmt.Sprintf("  [%s]", f.Op)
		}

		fmt.Printf("  %-*s ← %s%s\n", nameWidth, f.Name, relFragFile, opTag)
	}

	fmt.Println()
	fmt.Printf("Output: %s\n", r.Output)
}

func relPath(base, target string) string {
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return target
	}
	return rel
}
