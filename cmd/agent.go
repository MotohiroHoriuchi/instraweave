package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/MotohiroHoriuchi/instraweave/internal/agentprompt"
	"github.com/MotohiroHoriuchi/instraweave/internal/target"
	"github.com/spf13/cobra"
)

func init() {
	var targetName string
	var force bool

	agentCmd := &cobra.Command{
		Use:   "agent",
		Short: "Install AI agent prompt/command files for using instrweave",
		RunE: func(cmd *cobra.Command, args []string) error {
			paths, err := target.AgentFiles(targetName)
			if err != nil {
				return err
			}

			files := map[string]string{
				paths.UseCommand:       agentprompt.Get(targetName, "use"),
				paths.DecomposeCommand: agentprompt.Get(targetName, "decompose"),
			}

			for path, content := range files {
				if !force {
					if _, err := os.Stat(path); err == nil {
						return fmt.Errorf("%s already exists (use --force to overwrite)", path)
					}
				}
				if dir := filepath.Dir(path); dir != "." {
					if err := os.MkdirAll(dir, 0o755); err != nil {
						return fmt.Errorf("failed to create directory: %w", err)
					}
				}
				if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
					return fmt.Errorf("failed to write %s: %w", path, err)
				}
				fmt.Printf("Created %s\n", path)
			}
			return nil
		},
	}

	agentCmd.Flags().StringVarP(&targetName, "target", "t", "", "target agent: claude or copilot (required)")
	agentCmd.Flags().BoolVar(&force, "force", false, "overwrite existing files")
	_ = agentCmd.MarkFlagRequired("target")
	rootCmd.AddCommand(agentCmd)
}
