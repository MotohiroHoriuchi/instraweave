package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MotohiroHoriuchi/instraweave/internal/fragment"
	"github.com/spf13/cobra"
)

func init() {
	var dir string
	var verbose bool

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List available fragments",
		RunE: func(cmd *cobra.Command, args []string) error {
			fragments, err := fragment.List(dir)
			if err != nil {
				return err
			}

			if len(fragments) == 0 {
				fmt.Println("No fragments found.")
				return nil
			}

			for i, f := range fragments {
				if !verbose {
					fmt.Println(f)
					continue
				}

				data, err := os.ReadFile(filepath.Join(dir, f+".md"))
				if err != nil {
					return fmt.Errorf("failed to read fragment %q: %w", f, err)
				}

				if i > 0 {
					fmt.Println()
				}
				fmt.Printf("=== %s ===\n", f)
				fmt.Println(strings.TrimSpace(string(data)))
			}
			return nil
		},
	}

	listCmd.Flags().StringVarP(&dir, "dir", "d", "./fragments", "fragments directory")
	listCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "show fragment contents")
	rootCmd.AddCommand(listCmd)
}
