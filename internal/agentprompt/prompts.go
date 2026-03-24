package agentprompt

import _ "embed"

//go:embed claude-instrweave.md
var claudeUse string

//go:embed claude-instrweave-decompose.md
var claudeDecompose string

//go:embed copilot-instrweave.prompt.md
var copilotUse string

//go:embed copilot-instrweave-decompose.prompt.md
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
