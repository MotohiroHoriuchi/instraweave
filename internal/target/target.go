package target

import "fmt"

var defaultOutputPaths = map[string]string{
	"copilot": ".github/copilot-instructions.md",
	"claude":  "CLAUDE.md",
}

func DefaultOutputPath(target string) (string, error) {
	path, ok := defaultOutputPaths[target]
	if !ok {
		return "", fmt.Errorf("unknown target: %q (supported: copilot, claude)", target)
	}
	return path, nil
}

func SupportedTargets() []string {
	return []string{"copilot", "claude"}
}
