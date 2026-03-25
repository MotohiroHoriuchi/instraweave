package recipe

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MotohiroHoriuchi/instraweave/internal/target"
	"gopkg.in/yaml.v3"
)

// FragmentEntry represents one item in a recipe's fragments list.
// It can be a plain name (root recipe only) or an operation (add/remove/override).
type FragmentEntry struct {
	Op   string // "plain", "add", "remove", "override"
	Name string
}

func (f *FragmentEntry) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {
	case yaml.ScalarNode:
		f.Op = "plain"
		f.Name = value.Value
		return nil
	case yaml.MappingNode:
		if len(value.Content) != 2 {
			return fmt.Errorf("fragment entry must have exactly one key")
		}
		key := value.Content[0].Value
		val := value.Content[1].Value
		switch key {
		case "add", "remove", "override":
			f.Op = key
			f.Name = val
		default:
			return fmt.Errorf("unknown fragment operation %q", key)
		}
		return nil
	}
	return fmt.Errorf("invalid fragment entry format")
}

// rawRecipe holds the parsed YAML of a single recipe file before inheritance resolution.
type rawRecipe struct {
	Extends      string          `yaml:"extends"`
	Target       string          `yaml:"target"`
	Output       string          `yaml:"output"`
	FragmentsDir string          `yaml:"fragments_dir"`
	Fragments    []FragmentEntry `yaml:"fragments"`
	absPath      string
	absDir       string
	absFragsDir  string
}

// ResolvedFragment is a fragment after full inheritance resolution.
type ResolvedFragment struct {
	Name         string // e.g. "standard/security"
	FragmentsDir string // absolute path to the fragments_dir owning this fragment
	Op           string // "plain", "add", "override" (for dry-run display)
	SourceRecipe string // absolute path to recipe that last modified this fragment
}

// ResolvedRecipe is the final result after all inheritance has been applied.
type ResolvedRecipe struct {
	Chain     []string           // absolute recipe paths from root to current
	Fragments []ResolvedFragment
	Target    string
	Output    string
}

func loadRaw(absPath string) (*rawRecipe, error) {
	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read recipe file %q: %w", absPath, err)
	}
	var r rawRecipe
	if err := yaml.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("failed to parse recipe file %q: %w", absPath, err)
	}
	r.absPath = absPath
	r.absDir = filepath.Dir(absPath)

	fragsDir := r.FragmentsDir
	if fragsDir == "" {
		fragsDir = "./fragments"
	}
	if !filepath.IsAbs(fragsDir) {
		fragsDir = filepath.Join(r.absDir, fragsDir)
	}
	r.absFragsDir = fragsDir
	return &r, nil
}

// buildChain builds the inheritance chain [root, ..., current] for the given path.
// ancestors contains the paths traversed so far (outermost to current's parent),
// used for circular reference detection.
func buildChain(absPath string, ancestors []string) ([]*rawRecipe, error) {
	for i, p := range ancestors {
		if p == absPath {
			cycle := append(ancestors[i:], absPath)
			return nil, fmt.Errorf("circular extends detected: %s", strings.Join(cycle, " → "))
		}
	}

	r, err := loadRaw(absPath)
	if err != nil {
		return nil, err
	}

	if r.Extends == "" {
		return []*rawRecipe{r}, nil
	}

	extendsPath := r.Extends
	if !filepath.IsAbs(extendsPath) {
		extendsPath = filepath.Join(r.absDir, extendsPath)
	}
	extendsPath = filepath.Clean(extendsPath)

	if _, err := os.Stat(extendsPath); err != nil {
		return nil, fmt.Errorf("extends not found: %q (referenced from %q)", r.Extends, absPath)
	}

	parentChain, err := buildChain(extendsPath, append(ancestors, absPath))
	if err != nil {
		return nil, err
	}

	return append(parentChain, r), nil
}

// Load loads a recipe file and resolves all inheritance, returning a ResolvedRecipe.
func Load(path string) (*ResolvedRecipe, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path %q: %w", path, err)
	}

	chain, err := buildChain(absPath, nil)
	if err != nil {
		return nil, err
	}

	root := chain[0]
	for _, f := range root.Fragments {
		if f.Op != "plain" {
			return nil, fmt.Errorf("recipe %q: operation %q not allowed in root recipe (extends not set)", root.absPath, f.Op)
		}
	}

	type fragState struct {
		Name         string
		FragmentsDir string
		Op           string
		SourceRecipe string
	}

	var frags []fragState
	for _, f := range root.Fragments {
		frags = append(frags, fragState{
			Name:         f.Name,
			FragmentsDir: root.absFragsDir,
			Op:           "plain",
			SourceRecipe: root.absPath,
		})
	}

	for _, r := range chain[1:] {
		for _, entry := range r.Fragments {
			switch entry.Op {
			case "plain":
				return nil, fmt.Errorf("recipe %q: bare fragment %q not allowed in derived recipe; use \"add: %s\"", r.absPath, entry.Name, entry.Name)
			case "add":
				for _, f := range frags {
					if f.Name == entry.Name {
						return nil, fmt.Errorf("recipe %q: add conflict: fragment %q already exists", r.absPath, entry.Name)
					}
				}
				frags = append(frags, fragState{
					Name:         entry.Name,
					FragmentsDir: r.absFragsDir,
					Op:           "add",
					SourceRecipe: r.absPath,
				})
			case "remove":
				found := false
				for i, f := range frags {
					if f.Name == entry.Name {
						frags = append(frags[:i], frags[i+1:]...)
						found = true
						_ = f
						break
					}
				}
				if !found {
					return nil, fmt.Errorf("recipe %q: remove not found: fragment %q does not exist", r.absPath, entry.Name)
				}
			case "override":
				found := false
				for i, f := range frags {
					if f.Name == entry.Name {
						frags[i].FragmentsDir = r.absFragsDir
						frags[i].Op = "override"
						frags[i].SourceRecipe = r.absPath
						found = true
						_ = f
						break
					}
				}
				if !found {
					return nil, fmt.Errorf("recipe %q: override not found: fragment %q does not exist", r.absPath, entry.Name)
				}
			}
		}
	}

	// Resolve target and output: last non-empty value in chain wins.
	resolvedTarget := ""
	resolvedOutput := ""
	for _, r := range chain {
		if r.Target != "" {
			resolvedTarget = r.Target
		}
		if r.Output != "" {
			resolvedOutput = r.Output
		}
	}

	if resolvedTarget == "" {
		return nil, fmt.Errorf("recipe: target is required")
	}
	if _, err := target.DefaultOutputPath(resolvedTarget); err != nil {
		return nil, fmt.Errorf("recipe: %w", err)
	}
	if len(frags) == 0 {
		return nil, fmt.Errorf("recipe: at least one fragment is required")
	}
	if resolvedOutput == "" {
		resolvedOutput, _ = target.DefaultOutputPath(resolvedTarget)
	}

	chainPaths := make([]string, len(chain))
	for i, r := range chain {
		chainPaths[i] = r.absPath
	}

	resolved := make([]ResolvedFragment, len(frags))
	for i, f := range frags {
		resolved[i] = ResolvedFragment{
			Name:         f.Name,
			FragmentsDir: f.FragmentsDir,
			Op:           f.Op,
			SourceRecipe: f.SourceRecipe,
		}
	}

	return &ResolvedRecipe{
		Chain:     chainPaths,
		Fragments: resolved,
		Target:    resolvedTarget,
		Output:    resolvedOutput,
	}, nil
}
