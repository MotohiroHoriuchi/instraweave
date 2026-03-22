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
	var all bool

	showCmd := &cobra.Command{
		Use:   "show [fragment...]",
		Short: "Show the content of one or more fragments",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !all && len(args) == 0 {
				return fmt.Errorf("specify fragment name(s) or use --all")
			}

			if all {
				fragments, err := fragment.List(dir)
				if err != nil {
					return err
				}
				if len(fragments) == 0 {
					fmt.Println("No fragments found.")
					return nil
				}
				args = fragments
			}

			for i, name := range args {
				path := filepath.Join(dir, name+".md")
				data, err := os.ReadFile(path)
				if err != nil {
					return fmt.Errorf("failed to read fragment %q: %w", name, err)
				}

				if i > 0 {
					fmt.Println()
				}
				fmt.Printf("=== %s ===\n", name)
				fmt.Println(strings.TrimSpace(string(data)))
			}
			return nil
		},
	}

	showCmd.Flags().StringVarP(&dir, "dir", "d", "./fragments", "fragments directory")
	showCmd.Flags().BoolVar(&all, "all", false, "show all fragments")
	rootCmd.AddCommand(showCmd)
}
