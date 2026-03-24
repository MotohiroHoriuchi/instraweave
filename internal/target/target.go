package target

import "fmt"

var defaultOutputPaths = map[string]string{
	"copilot": ".github/copilot-instructions.md",
	"claude":  "CLAUDE.md",
}

type AgentPaths struct {
	UseCommand        string
	DecomposeCommand  string
}

var agentFilePaths = map[string]AgentPaths{
	"claude": {
		UseCommand:       ".claude/commands/instrweave.md",
		DecomposeCommand: ".claude/commands/instrweave-decompose.md",
	},
	"copilot": {
		UseCommand:       ".github/prompts/instrweave.prompt.md",
		DecomposeCommand: ".github/prompts/instrweave-decompose.prompt.md",
	},
}

func DefaultOutputPath(target string) (string, error) {
	path, ok := defaultOutputPaths[target]
	if !ok {
		return "", fmt.Errorf("unknown target: %q (supported: copilot, claude)", target)
	}
	return path, nil
}

func AgentFiles(target string) (AgentPaths, error) {
	paths, ok := agentFilePaths[target]
	if !ok {
		return AgentPaths{}, fmt.Errorf("unknown target: %q (supported: copilot, claude)", target)
	}
	return paths, nil
}

func SupportedTargets() []string {
	return []string{"copilot", "claude"}
}
