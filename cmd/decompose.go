package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/spf13/cobra"
)

func init() {
	var filePath string
	var level int
	var outputDir string

	decomposeCmd := &cobra.Command{
		Use:   "decompose",
		Short: "Decompose a markdown file into fragment files by header level",
		RunE: func(cmd *cobra.Command, args []string) error {
			if level < 1 || level > 6 {
				return fmt.Errorf("--level must be between 1 and 6")
			}

			data, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}

			sections, err := splitByHeader(string(data), level)
			if err != nil {
				return err
			}

			if len(sections) == 0 {
				return fmt.Errorf("no level-%d headers found in %s", level, filePath)
			}

			if err := os.MkdirAll(outputDir, 0o755); err != nil {
				return fmt.Errorf("failed to create output directory: %w", err)
			}

			for _, s := range sections {
				name := slugify(s.title) + ".md"
				dest := filepath.Join(outputDir, name)
				if err := os.WriteFile(dest, []byte(s.content), 0o644); err != nil {
					return fmt.Errorf("failed to write %s: %w", dest, err)
				}
				fmt.Printf("Created %s\n", dest)
			}
			return nil
		},
	}

	decomposeCmd.Flags().StringVarP(&filePath, "file", "f", "", "markdown file to decompose (required)")
	decomposeCmd.Flags().IntVarP(&level, "level", "l", 2, "header level used as split boundary (1-6)")
	decomposeCmd.Flags().StringVarP(&outputDir, "dir", "d", "./fragments", "output directory for fragment files")
	_ = decomposeCmd.MarkFlagRequired("file")
	rootCmd.AddCommand(decomposeCmd)
}

type section struct {
	title   string
	content string
}

func splitByHeader(text string, level int) ([]section, error) {
	prefix := strings.Repeat("#", level) + " "
	lines := strings.Split(text, "\n")

	var sections []section
	var current *section

	for _, line := range lines {
		if isHeader(line, level) {
			if current != nil {
				current.content = strings.TrimRight(current.content, "\n") + "\n"
				sections = append(sections, *current)
			}
			title := strings.TrimPrefix(line, prefix)
			current = &section{
				title:   strings.TrimSpace(title),
				content: line + "\n",
			}
		} else if current != nil {
			current.content += line + "\n"
		}
	}
	if current != nil {
		current.content = strings.TrimRight(current.content, "\n") + "\n"
		sections = append(sections, *current)
	}

	return sections, nil
}

// isHeader returns true only if line is exactly at the given header level.
// e.g. level=2 matches "## Foo" but not "### Foo".
func isHeader(line string, level int) bool {
	prefix := strings.Repeat("#", level) + " "
	if !strings.HasPrefix(line, prefix) {
		return false
	}
	// Ensure it's not a deeper header (e.g. "### " when level=2 would not match
	// because "## " != "### "[0:3], but be explicit for level 1)
	if level < 6 && strings.HasPrefix(line, strings.Repeat("#", level+1)) {
		return false
	}
	return true
}

func slugify(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case unicode.IsSpace(r) || r == '-' || r == '_':
			b.WriteRune('-')
		}
	}
	result := b.String()
	for strings.Contains(result, "--") {
		result = strings.ReplaceAll(result, "--", "-")
	}
	return strings.Trim(result, "-")
}
