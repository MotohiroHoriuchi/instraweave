package recipe

import (
	"fmt"
	"os"

	"github.com/MotohiroHoriuchi/instraweave/internal/target"
	"gopkg.in/yaml.v3"
)

type Recipe struct {
	Target       string   `yaml:"target"`
	Output       string   `yaml:"output"`
	FragmentsDir string   `yaml:"fragments_dir"`
	Fragments    []string `yaml:"fragments"`
}

func Load(path string) (*Recipe, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read recipe file: %w", err)
	}

	var r Recipe
	if err := yaml.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("failed to parse recipe file: %w", err)
	}

	if err := r.validate(); err != nil {
		return nil, err
	}

	r.applyDefaults()
	return &r, nil
}

func (r *Recipe) validate() error {
	if r.Target == "" {
		return fmt.Errorf("recipe: target is required")
	}
	if _, err := target.DefaultOutputPath(r.Target); err != nil {
		return fmt.Errorf("recipe: %w", err)
	}
	if len(r.Fragments) == 0 {
		return fmt.Errorf("recipe: at least one fragment is required")
	}
	return nil
}

func (r *Recipe) applyDefaults() {
	if r.FragmentsDir == "" {
		r.FragmentsDir = "./fragments"
	}
	if r.Output == "" {
		r.Output, _ = target.DefaultOutputPath(r.Target)
	}
}
