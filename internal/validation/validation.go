package validation

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/MotohiroHoriuchi/instraweave/internal/fragment"
	"github.com/MotohiroHoriuchi/instraweave/internal/recipe"
	"github.com/MotohiroHoriuchi/instraweave/internal/target"
	"gopkg.in/yaml.v3"
)

type Level int

const (
	LevelError   Level = iota
	LevelWarning Level = iota
	LevelInfo    Level = iota
	LevelOK      Level = iota
)

type Issue struct {
	Level   Level
	Message string
	Detail  string // optional detail (absolute path; caller may relativize)
}

type Result struct {
	Issues []Issue
}

func (r *Result) ErrorCount() int {
	n := 0
	for _, i := range r.Issues {
		if i.Level == LevelError {
			n++
		}
	}
	return n
}

func (r *Result) WarningCount() int {
	n := 0
	for _, i := range r.Issues {
		if i.Level == LevelWarning {
			n++
		}
	}
	return n
}

func (r *Result) HasErrors() bool { return r.ErrorCount() > 0 }

func (r *Result) add(level Level, msg string, detail ...string) {
	d := ""
	if len(detail) > 0 {
		d = detail[0]
	}
	r.Issues = append(r.Issues, Issue{Level: level, Message: msg, Detail: d})
}

// Validate checks the recipe and fragments for correctness.
// fragmentsDir is used to detect unreferenced fragments; if empty, the check is skipped.
func Validate(recipePath, fragmentsDir string) *Result {
	result := &Result{}

	// 1. Recipe syntax
	data, err := os.ReadFile(recipePath)
	if err != nil {
		result.add(LevelError, "recipe not found: "+err.Error())
		return result
	}

	var rawData struct {
		Target string `yaml:"target"`
	}
	if err := yaml.Unmarshal(data, &rawData); err != nil {
		result.add(LevelError, "recipe syntax error: "+err.Error())
		return result
	}
	result.add(LevelOK, "recipe syntax OK")

	// 2. Target value
	supported := false
	for _, t := range target.SupportedTargets() {
		if t == rawData.Target {
			supported = true
			break
		}
	}
	if !supported {
		result.add(LevelError, fmt.Sprintf("unsupported target: %q (supported: %s)",
			rawData.Target, strings.Join(target.SupportedTargets(), ", ")))
	} else {
		result.add(LevelOK, fmt.Sprintf("target %q is supported", rawData.Target))
	}

	// 3. Load and resolve recipe (validates inheritance, operations, etc.)
	r, err := recipe.Load(recipePath)
	if err != nil {
		result.add(LevelError, err.Error())
		return result
	}

	// 4. Check each referenced fragment
	referencedPaths := map[string]bool{}
	for _, f := range r.Fragments {
		fragPath := filepath.Join(f.FragmentsDir, f.Name+".md")
		referencedPaths[fragPath] = true

		fragData, err := os.ReadFile(fragPath)
		if err != nil {
			result.add(LevelError, "fragment not found: "+f.Name, fragPath)
			continue
		}
		if !utf8.Valid(fragData) {
			result.add(LevelError, "fragment not UTF-8: "+f.Name)
			continue
		}
		if strings.TrimSpace(string(fragData)) == "" {
			result.add(LevelWarning, "fragment is empty: "+f.Name)
		}
	}

	// 5. Unreferenced fragments
	if fragmentsDir != "" {
		absFragsDir, err := filepath.Abs(fragmentsDir)
		if err == nil {
			allFrags, err := fragment.List(absFragsDir)
			if err == nil {
				var unreferenced []string
				for _, name := range allFrags {
					fragPath := filepath.Join(absFragsDir, name+".md")
					if !referencedPaths[fragPath] {
						unreferenced = append(unreferenced, name)
					}
				}
				if len(unreferenced) > 0 {
					result.add(LevelInfo, "unreferenced fragments: "+strings.Join(unreferenced, ", "))
				}
			}
		}
	}

	return result
}
