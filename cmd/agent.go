package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/MotohiroHoriuchi/instraweave/internal/target"
	"github.com/spf13/cobra"
)

// claudeUseCommand is the slash command placed at .claude/commands/instrweave.md
const claudeUseCommand = `Review and update this project's AI agent instructions using instrweave.

Steps:
1. Run ` + "`instrweave list --verbose`" + ` to see all available fragments and their contents.
2. Run ` + "`cat instrweave-recipe.yaml`" + ` to review the current recipe.
3. Consider the project's current needs and suggest which fragments to add or remove.
4. Edit ` + "`instrweave-recipe.yaml`" + ` to reflect the desired changes.
5. Run ` + "`instrweave generate --dry-run`" + ` to preview the composed output.
6. Run ` + "`instrweave generate`" + ` to write the instructions file.
`

// claudeDecomposeCommand is the slash command placed at .claude/commands/instrweave-decompose.md
const claudeDecomposeCommand = `Decompose an existing instructions markdown file into instrweave fragments.

Usage: /instrweave-decompose <file-path> [header-level] [output-dir]

Steps:
1. Run ` + "`instrweave decompose --file <file-path> --level <n> --dir ./fragments/custom/`" + `
   - ` + "`--level`" + `: header level used as split boundary (default: 2, i.e. ## headers)
   - ` + "`--dir`" + `: output directory for generated fragment files (default: ./fragments)
2. Review the generated fragment files.
3. Rename files or adjust content as necessary.
4. Add the new fragments to ` + "`instrweave-recipe.yaml`" + `.
5. Run ` + "`instrweave generate --dry-run`" + ` to verify the result.
`

// copilotUsePrompt is placed at .github/prompts/instrweave.prompt.md
const copilotUsePrompt = `---
mode: agent
description: Review and update AI agent instructions using instrweave
---

Review and update this project's AI agent instructions using instrweave.

Steps:
1. Run ` + "`instrweave list --verbose`" + ` to see all available fragments and their contents.
2. Run ` + "`cat instrweave-recipe.yaml`" + ` to review the current recipe.
3. Consider the project's current needs and suggest which fragments to add or remove.
4. Edit ` + "`instrweave-recipe.yaml`" + ` to reflect the desired changes.
5. Run ` + "`instrweave generate --dry-run`" + ` to preview the composed output.
6. Run ` + "`instrweave generate`" + ` to write the instructions file.
`

// copilotDecomposePrompt is placed at .github/prompts/instrweave-decompose.prompt.md
const copilotDecomposePrompt = `---
mode: agent
description: Decompose an existing instructions file into instrweave fragments
---

Decompose an existing instructions markdown file into instrweave fragments.

Steps:
1. Identify the target markdown file (e.g. ` + "`CLAUDE.md`" + `, ` + "`.github/copilot-instructions.md`" + `).
2. Choose an appropriate header level for splitting (e.g. 2 for ## headers).
3. Run ` + "`instrweave decompose --file <path> --level <n> --dir ./fragments/custom/`" + `
   - ` + "`--level`" + `: header level used as split boundary (default: 2)
   - ` + "`--dir`" + `: output directory for generated fragment files (default: ./fragments)
4. Review the generated fragment files.
5. Rename files or adjust content as necessary.
6. Add the new fragments to ` + "`instrweave-recipe.yaml`" + `.
7. Run ` + "`instrweave generate --dry-run`" + ` to verify the result.
`

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
				paths.UseCommand:       contentFor(targetName, "use"),
				paths.DecomposeCommand: contentFor(targetName, "decompose"),
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

func contentFor(t, kind string) string {
	switch t + "/" + kind {
	case "claude/use":
		return claudeUseCommand
	case "claude/decompose":
		return claudeDecomposeCommand
	case "copilot/use":
		return copilotUsePrompt
	case "copilot/decompose":
		return copilotDecomposePrompt
	}
	return ""
}
