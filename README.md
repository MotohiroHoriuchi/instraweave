# instrweave

A CLI tool that composes AI agent instruction files from reusable rule fragments.

[日本語版 README](README.ja.md)

## Overview

`instrweave` lets you manage AI coding agent instructions (e.g., GitHub Copilot, Claude) as modular Markdown fragments organized by category. Define a YAML recipe to select which fragments to include, and `instrweave` assembles them into a single instructions file.

## Installation

```bash
go install github.com/motohirohoriuchi/instrweave@latest
```

Or build from source:

```bash
git clone https://github.com/motohirohoriuchi/instrweave.git
cd instrweave
go build -o instrweave .
```

## Quick Start

```bash
# 1. Initialize a sample recipe and fragments directory
instrweave init

# 2. List available fragments
instrweave list

# 3. Preview the composed output
instrweave generate --dry-run

# 4. Generate the instructions file
instrweave generate
```

## Recipe File

`instrweave` uses a YAML recipe file (`instrweave-recipe.yaml`) to define what to generate:

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

### `instrweave init`

Creates a sample `instrweave-recipe.yaml` and `fragments/` directory in the current directory.

```bash
instrweave init
```

### `instrweave list`

Lists all available fragments in the specified directory.

```bash
instrweave list
instrweave list --dir ./my-fragments
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--dir` | `-d` | `./fragments` | Fragments directory |
| `--verbose` | `-v` | `false` | Show fragment contents |

### `instrweave show`

Shows the content of one or more fragments. Useful for AI agents to inspect fragments before building a recipe.

```bash
instrweave show standard/go
instrweave show standard/go standard/testing
instrweave show --all
instrweave show --all --dir ./my-fragments
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--dir` | `-d` | `./fragments` | Fragments directory |
| `--all` | | `false` | Show all fragments |

### `instrweave generate`

Reads the recipe file and composes fragments into an instructions file.

```bash
instrweave generate
instrweave generate --recipe ./my-recipe.yaml
instrweave generate --dry-run
```

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--recipe` | `-r` | `./instrweave-recipe.yaml` | Path to recipe file |
| `--dry-run` | | `false` | Print to stdout instead of writing to file |

## Example

See the [`examples/fragments/`](examples/fragments/) directory for sample fragments.

## License

MIT
