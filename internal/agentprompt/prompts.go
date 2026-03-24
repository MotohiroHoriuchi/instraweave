package agentprompt

import _ "embed"

//go:embed claude-instraweave.md
var claudeUse string

//go:embed claude-instraweave-decompose.md
var claudeDecompose string

//go:embed copilot-instraweave.prompt.md
var copilotUse string

//go:embed copilot-instraweave-decompose.prompt.md
var copilotDecompose string

func Get(target, kind string) string {
	switch target + "/" + kind {
	case "claude/use":
		return claudeUse
	case "claude/decompose":
		return claudeDecompose
	case "copilot/use":
		return copilotUse
	case "copilot/decompose":
		return copilotDecompose
	}
	return ""
}
