package fragment

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Entry is a resolved fragment with its source directory.
type Entry struct {
	Name string
	Dir  string
}

// Compose assembles fragments from a single directory (used for simple recipes without inheritance).
func Compose(fragmentsDir string, names []string) (string, error) {
	entries := make([]Entry, len(names))
	for i, name := range names {
		entries[i] = Entry{Name: name, Dir: fragmentsDir}
	}
	return ComposeEntries(entries)
}

// ComposeEntries assembles fragments where each entry may come from a different directory.
func ComposeEntries(entries []Entry) (string, error) {
	var parts []string

	for _, e := range entries {
		path := filepath.Join(e.Dir, e.Name+".md")
		data, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("fragment file not found: %q", path)
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
