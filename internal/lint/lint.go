// Package lint provides fragment content quality checks.
package lint

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Level represents the severity of a lint rule.
type Level string

const (
	LevelError   Level = "error"
	LevelWarning Level = "warning"
	LevelOff     Level = "off"
)

// Issue is a single lint finding.
type Issue struct {
	Rule    string
	Level   Level
	Message string
	Line    int // 0 means no specific line
}

// RulePass records a rule that passed, with optional detail info.
type RulePass struct {
	Rule   string
	Detail string // e.g. "(842 chars)"
}

// FragmentResult holds all issues for one fragment.
type FragmentResult struct {
	Name    string
	Issues  []Issue
	Passing []RulePass
}

func (r *FragmentResult) ErrorCount() int {
	n := 0
	for _, i := range r.Issues {
		if i.Level == LevelError {
			n++
		}
	}
	return n
}

func (r *FragmentResult) WarningCount() int {
	n := 0
	for _, i := range r.Issues {
		if i.Level == LevelWarning {
			n++
		}
	}
	return n
}

// RuleConfig holds configuration for a single rule.
type RuleConfig struct {
	Level   Level                  `yaml:"level"`
	Options map[string]interface{} `yaml:"options"`
}

// Config holds the full lint configuration.
type Config struct {
	Rules map[string]RuleConfig `yaml:"rules"`
}

// defaultConfig returns the built-in default lint configuration.
func defaultConfig() Config {
	return Config{
		Rules: map[string]RuleConfig{
			"require-h1": {
				Level: LevelWarning,
			},
			"max-length": {
				Level: LevelWarning,
				Options: map[string]interface{}{
					"max": 2000,
				},
			},
			"no-trailing-whitespace": {
				Level: LevelWarning,
			},
		},
	}
}

// LoadConfig loads a lint config file, merging with defaults.
// If path doesn't exist, returns default config with no error.
func LoadConfig(path string) (Config, error) {
	cfg := defaultConfig()

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return cfg, fmt.Errorf("failed to read lint config %q: %w", path, err)
	}

	var fileCfg Config
	if err := yaml.Unmarshal(data, &fileCfg); err != nil {
		return cfg, fmt.Errorf("failed to parse lint config %q: %w", path, err)
	}

	// Merge file config into defaults
	for ruleID, fileRule := range fileCfg.Rules {
		existing, ok := cfg.Rules[ruleID]
		if !ok {
			cfg.Rules[ruleID] = fileRule
			continue
		}
		if fileRule.Level != "" {
			existing.Level = fileRule.Level
		}
		if fileRule.Options != nil {
			if existing.Options == nil {
				existing.Options = make(map[string]interface{})
			}
			for k, v := range fileRule.Options {
				existing.Options[k] = v
			}
		}
		cfg.Rules[ruleID] = existing
	}

	return cfg, nil
}

// LintFragment runs all enabled rules against a single fragment file.
func LintFragment(name, path string, cfg Config) (*FragmentResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read fragment %q: %w", path, err)
	}
	content := string(data)
	lines := strings.Split(content, "\n")

	result := &FragmentResult{Name: name}

	// Use a stable rule order for consistent output
	ruleOrder := []string{"require-h1", "max-length", "no-trailing-whitespace", "forbidden-words", "require-metadata"}
	seen := map[string]bool{}
	for _, id := range ruleOrder {
		seen[id] = true
	}
	// Append any custom rules not in the fixed order
	for id := range cfg.Rules {
		if !seen[id] {
			ruleOrder = append(ruleOrder, id)
		}
	}

	for _, ruleID := range ruleOrder {
		ruleCfg, ok := cfg.Rules[ruleID]
		if !ok || ruleCfg.Level == LevelOff {
			continue
		}
		issues, pass := applyRule(ruleID, ruleCfg, content, lines)
		if len(issues) > 0 {
			result.Issues = append(result.Issues, issues...)
		} else if pass != nil {
			result.Passing = append(result.Passing, *pass)
		}
	}

	sortIssues(result.Issues)
	return result, nil
}

func applyRule(ruleID string, cfg RuleConfig, content string, lines []string) ([]Issue, *RulePass) {
	switch ruleID {
	case "require-h1":
		return ruleRequireH1(cfg, lines)
	case "max-length":
		return ruleMaxLength(cfg, content)
	case "no-trailing-whitespace":
		return ruleNoTrailingWhitespace(cfg, lines)
	case "forbidden-words":
		return ruleForbiddenWords(cfg, lines)
	case "require-metadata":
		return ruleRequireMetadata(cfg, content)
	}
	return nil, nil
}

func ruleRequireH1(cfg RuleConfig, lines []string) ([]Issue, *RulePass) {
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "# ") || trimmed == "#" {
			return nil, &RulePass{Rule: "require-h1"}
		}
		break
	}
	return []Issue{{Rule: "require-h1", Level: cfg.Level, Message: "no H1 header found"}}, nil
}

func ruleMaxLength(cfg RuleConfig, content string) ([]Issue, *RulePass) {
	maxChars := 2000
	if cfg.Options != nil {
		if v, ok := cfg.Options["max"]; ok {
			switch val := v.(type) {
			case int:
				maxChars = val
			case float64:
				maxChars = int(val)
			}
		}
	}
	chars := len([]rune(content))
	if chars > maxChars {
		return []Issue{{
			Rule:    "max-length",
			Level:   cfg.Level,
			Message: fmt.Sprintf("%d chars (limit: %d)", chars, maxChars),
		}}, nil
	}
	return nil, &RulePass{Rule: "max-length", Detail: fmt.Sprintf("(%d chars)", chars)}
}

func ruleNoTrailingWhitespace(cfg RuleConfig, lines []string) ([]Issue, *RulePass) {
	var issues []Issue
	for i, line := range lines {
		if line != strings.TrimRight(line, " \t") {
			issues = append(issues, Issue{
				Rule:    "no-trailing-whitespace",
				Level:   cfg.Level,
				Message: fmt.Sprintf("trailing whitespace at line %d", i+1),
				Line:    i + 1,
			})
		}
	}
	if len(issues) > 0 {
		return issues, nil
	}
	return nil, &RulePass{Rule: "no-trailing-whitespace"}
}

func ruleForbiddenWords(cfg RuleConfig, lines []string) ([]Issue, *RulePass) {
	if cfg.Options == nil {
		return nil, &RulePass{Rule: "forbidden-words"}
	}
	raw, ok := cfg.Options["words"]
	if !ok {
		return nil, &RulePass{Rule: "forbidden-words"}
	}

	var words []string
	switch v := raw.(type) {
	case []interface{}:
		for _, w := range v {
			if s, ok := w.(string); ok {
				words = append(words, s)
			}
		}
	case []string:
		words = v
	}

	var issues []Issue
	for i, line := range lines {
		for _, word := range words {
			if strings.Contains(line, word) {
				issues = append(issues, Issue{
					Rule:    "forbidden-words",
					Level:   cfg.Level,
					Message: fmt.Sprintf("%q found at line %d", word, i+1),
					Line:    i + 1,
				})
			}
		}
	}
	if len(issues) > 0 {
		return issues, nil
	}
	return nil, &RulePass{Rule: "forbidden-words"}
}

func ruleRequireMetadata(cfg RuleConfig, content string) ([]Issue, *RulePass) {
	if cfg.Options == nil {
		return nil, &RulePass{Rule: "require-metadata"}
	}
	raw, ok := cfg.Options["required-keys"]
	if !ok {
		return nil, &RulePass{Rule: "require-metadata"}
	}

	var requiredKeys []string
	switch v := raw.(type) {
	case []interface{}:
		for _, k := range v {
			if s, ok := k.(string); ok {
				requiredKeys = append(requiredKeys, s)
			}
		}
	case []string:
		requiredKeys = v
	}

	if !strings.HasPrefix(strings.TrimSpace(content), "---") {
		return []Issue{{Rule: "require-metadata", Level: cfg.Level, Message: "no YAML front matter found"}}, nil
	}

	var issues []Issue
	for _, key := range requiredKeys {
		if !strings.Contains(content, key+":") {
			issues = append(issues, Issue{
				Rule:    "require-metadata",
				Level:   cfg.Level,
				Message: fmt.Sprintf("required metadata key %q not found", key),
			})
		}
	}
	if len(issues) > 0 {
		return issues, nil
	}
	return nil, &RulePass{Rule: "require-metadata"}
}

func sortIssues(issues []Issue) {
	// Simple insertion sort by line number
	for i := 1; i < len(issues); i++ {
		for j := i; j > 0 && issues[j].Line < issues[j-1].Line; j-- {
			issues[j], issues[j-1] = issues[j-1], issues[j]
		}
	}
}

// LintDir runs lint on all fragments in a directory.
func LintDir(dir string, cfg Config) ([]*FragmentResult, error) {
	var results []*FragmentResult

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
		result, err := LintFragment(name, path, cfg)
		if err != nil {
			return err
		}
		results = append(results, result)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to lint fragments: %w", err)
	}
	return results, nil
}
