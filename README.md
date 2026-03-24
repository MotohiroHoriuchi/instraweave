# instraweave

A CLI tool that composes AI agent instruction files from reusable rule fragments.

[日本語版 README](README.ja.md)

## Overview

`instraweave` lets you manage AI coding agent instructions (e.g., GitHub Copilot, Claude) as modular Markdown fragments organized by category. Define a YAML recipe to select which fragments to include, and `instraweave` assembles them into a single instructions file.

## Installation

```bash
go install github.com/MotohiroHoriuchi/instraweave@latest
```

Or build from source:

```bash
git clone https://github.com/MotohiroHoriuchi/instraweave.git
cd instraweave
go build -o instraweave .
```

## Quick Start

```bash
# 1. Initialize a sample recipe and fragments directory
instraweave init

# 2. List available fragments
instraweave list

# 3. Preview the composed output
instraweave generate --dry-run

# 4. Generate the instructions file
instraweave generate
```

## Recipe File

`instraweave` uses a YAML recipe file (`instraweave-recipe.yaml`) to define what to generate:

```yaml
target: copilot              # copilot | claude
output: ""                   # Leave empty to use target's default path
fragments_dir: ./fragments   # Directory containing fragment files (default: ./fragments)
fragments:
  - standard/go
  - standard/testing
  - standard/security
  - custom/our-api-convention
```

### Supported Targets

| Target | Default Output Path |
|--------|-------------------|
| `copilot` | `.github/copilot-instructions.md` |
| `claude` | `CLAUDE.md` |

## Fragment Structure

Fragments are Markdown files organized in subdirectories:

```
fragments/
├── standard/          # Shared, reusable rules
│   ├── go.md
│   ├── testing.md
│   └── security.md
└── custom/            # Project-specific rules
    └── our-api-convention.md
```

Fragment names in the recipe correspond to file paths under `fragments_dir`, without the `.md` extension.

## Commands

### `instraweave init`

Creates a sample `instraweave-recipe.yaml` and `fragments/` directory in the current directory.

```bash
instraweave init
```

### `instraweave list`

Lists all available fragments in the specified directory.

```bash
instraweave list
instraweave list --dir ./my-fragments
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--dir` | `-d` | `./fragments` | Fragments directory |
| `--verbose` | `-v` | `false` | Show fragment contents |

### `instraweave show`

Shows the content of one or more fragments. Useful for AI agents to inspect fragments before building a recipe.

```bash
instraweave show standard/go
instraweave show standard/go standard/testing
instraweave show --all
instraweave show --all --dir ./my-fragments
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--dir` | `-d` | `./fragments` | Fragments directory |
| `--all` | | `false` | Show all fragments |

### `instraweave generate`

Reads the recipe file and composes fragments into an instructions file.

```bash
instraweave generate
instraweave generate --recipe ./my-recipe.yaml
instraweave generate --dry-run
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--recipe` | `-r` | `./instraweave-recipe.yaml` | Path to recipe file |
| `--dry-run` | | `false` | Print to stdout instead of writing to file |

### `instraweave decompose`

Splits a single Markdown file into fragment files by header level.

```bash
instraweave decompose --file CLAUDE.md
instraweave decompose --file docs/guide.md --level 1 --dir ./fragments/custom/
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--file` | `-f` | *(required)* | Markdown file to decompose |
| `--level` | `-l` | `2` | Header level used as split boundary (1–6) |
| `--dir` | `-d` | `./fragments` | Output directory for fragment files |

### `instraweave agent`

Installs AI agent prompt/command files so your agent can manage instraweave directly.

```bash
instraweave agent --target claude
instraweave agent --target copilot
instraweave agent --target claude --force   # overwrite existing files
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--target` | `-t` | *(required)* | Agent target: `claude` or `copilot` |
| `--force` | | `false` | Overwrite existing files |

**Installed files per target:**

| Target | Use command | Decompose command |
|--------|-------------|-------------------|
| `claude` | `.claude/commands/instraweave.md` | `.claude/commands/instraweave-decompose.md` |
| `copilot` | `.github/prompts/instraweave.prompt.md` | `.github/prompts/instraweave-decompose.prompt.md` |

The **decompose command** guides the agent to decompose existing documents into instraweave fragments:

- **Header-based splitting** (preferred): uses `instraweave decompose` when consistent headers exist.
- **Semantic splitting** (fallback): when headers are absent or sparse, the agent infers logical topic boundaries from meaning and creates fragments manually.
- **Verbatim constraint**: body text is always copied as-is — no rewrites, paraphrasing, or additions.

## Example

See the [`examples/fragments/`](examples/fragments/) directory for sample fragments.

## License

MIT
