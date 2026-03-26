package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	inlint "github.com/MotohiroHoriuchi/instraweave/internal/lint"
	"github.com/spf13/cobra"
)

func init() {
	var fragmentsDir string
	var configPath string
	var strict bool

	lintCmd := &cobra.Command{
		Use:   "lint [fragment...]",
		Short: "Check fragment content quality",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := inlint.LoadConfig(configPath)
			if err != nil {
				return err
			}

			var results []*inlint.FragmentResult

			if len(args) > 0 {
				for _, name := range args {
					path := filepath.Join(fragmentsDir, name+".md")
					r, err := inlint.LintFragment(name, path, cfg)
					if err != nil {
						return err
					}
					results = append(results, r)
				}
			} else {
				results, err = inlint.LintDir(fragmentsDir, cfg)
				if err != nil {
					return err
				}
			}

			if len(results) == 0 {
				fmt.Println("No fragments found.")
				return nil
			}

			fmt.Printf("Linting %d fragment(s)...\n\n", len(results))

			totalErrors := 0
			totalWarnings := 0
			affectedFragments := 0

			for _, r := range results {
				fmt.Printf("%s\n", r.Name+".md")

				// Build a combined ordered list: passing rules first (in order), then issues
				// Issues are already sorted by line; passing rules are in rule order.
				// We interleave by printing passing first, then issues.
				for _, p := range r.Passing {
					detail := ""
					if p.Detail != "" {
						detail = " " + p.Detail
					}
					fmt.Printf("  ✓ %s%s\n", p.Rule, detail)
				}
				for _, issue := range r.Issues {
					switch issue.Level {
					case inlint.LevelError:
						totalErrors++
						fmt.Printf("  ✗ %s: %s\n", issue.Rule, issue.Message)
					case inlint.LevelWarning:
						totalWarnings++
						fmt.Printf("  ⚠ %s: %s\n", issue.Rule, issue.Message)
					}
				}

				if len(r.Issues) > 0 {
					affectedFragments++
				}
				fmt.Println()
			}

			fmt.Printf("Results: %d error(s), %d warning(s) in %d/%d fragment(s)\n",
				totalErrors, totalWarnings, affectedFragments, len(results))

			if totalErrors > 0 {
				os.Exit(1)
			}
			if totalWarnings > 0 && strict {
				os.Exit(1)
			}
			if totalWarnings > 0 {
				os.Exit(2)
			}
			return nil
		},
	}

	lintCmd.Flags().StringVarP(&fragmentsDir, "dir", "d", "./fragments", "path to fragments directory")
	lintCmd.Flags().StringVarP(&configPath, "config", "c", "./instraweave-lint.yaml", "path to lint config file")
	lintCmd.Flags().BoolVar(&strict, "strict", false, "exit with code 1 on warnings (default: exit 2)")
	rootCmd.AddCommand(lintCmd)
}
