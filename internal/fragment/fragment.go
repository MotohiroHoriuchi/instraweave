package fragment

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Compose(fragmentsDir string, names []string) (string, error) {
	var parts []string

	for _, name := range names {
		path := filepath.Join(fragmentsDir, name+".md")
		data, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("failed to read fragment %q: %w", name, err)
		}
		parts = append(parts, strings.TrimSpace(string(data)))
	}

	return strings.Join(parts, "\n\n") + "\n", nil
}

func List(dir string) ([]string, error) {
	var fragments []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		name := strings.TrimSuffix(rel, ".md")
		fragments = append(fragments, name)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list fragments: %w", err)
	}

	return fragments, nil
}
