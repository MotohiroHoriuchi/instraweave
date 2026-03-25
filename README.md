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

## Recipe Inheritance

Recipes can inherit from another recipe using the `extends` field. This allows you to share a base set of fragments across teams or projects, and customize them per layer.

### Basic Syntax

A **derived recipe** (one with `extends`) uses operations instead of plain fragment names:

```yaml
extends: ../base/recipe.yaml   # relative or absolute path

target: claude
fragments_dir: ./fragments

fragments:
  - add: standard/go           # append to the list
  - remove: standard/code-review  # remove from the list
  - override: standard/security   # replace with this recipe's version
```

A **root recipe** (no `extends`) lists fragments as plain names:

```yaml
target: claude
fragments_dir: ./fragments
fragments:
  - standard/security
  - standard/git-convention
  - standard/code-review
```

### Fragment Operations

| Operation | Syntax | Behavior |
|-----------|--------|----------|
| *(plain)* | `- category/name` | Root recipe only. Error if used in a derived recipe. |
| `add` | `- add: category/name` | Append to the list. Error if already present. |
| `remove` | `- remove: category/name` | Remove from the list. Error if not present. |
| `override` | `- override: category/name` | Replace the fragment's source with this recipe's `fragments_dir`. Error if not present. |

### Inheritance Chain

`extends` is resolved recursively. Operations are applied from root to derived (last wins):

```
org/recipe.yaml          ← root (plain fragments)
  └─ team/recipe.yaml    ← adds Go, removes code-review
       └─ project/recipe.yaml  ← overrides security, adds db-convention
```

Each fragment is read from the `fragments_dir` of the recipe that **last modified it**:

- Plain fragment in root → resolved from root's `fragments_dir`
- `add` → resolved from the recipe that added it
- `override` → resolved from the recipe that overrode it

`target` and `output` are also inherited; a derived recipe's value overrides its parent's.

### Directory Structure Example

```
org/
├── recipe.yaml
└── fragments/
    └── standard/
        ├── security.md
        ├── git-convention.md
        └── code-review.md

team-backend/
├── recipe.yaml          # extends: ../org/recipe.yaml
└── fragments/
    ├── standard/
    │   └── go.md
    └── custom/
        └── our-code-review.md

project-payment/
├── recipe.yaml          # extends: ../team-backend/recipe.yaml
└── fragments/
    └── standard/
        └── security.md  # overrides org's version
```

### dry-run Output

`instraweave generate --dry-run` shows the resolved inheritance chain and fragment sources:

```
Inheritance chain:
  org/recipe.yaml           (root)
       └─ team-backend/recipe.yaml
            └─ project-payment/recipe.yaml  (current)

Resolved fragments:
  standard/security        ← project-payment/fragments/standard/security.md  [override]
  standard/git-convention  ← org/fragments/standard/git-convention.md
  standard/go              ← team-backend/fragments/standard/go.md            [add]
  custom/our-code-review   ← team-backend/fragments/custom/our-code-review.md [add]

Output: CLAUDE.md
```

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

Creates a sample `instraweave-recipe.yaml` and `fragments/` directory with starter fragments (`fragments/standard/go.md` and `fragments/custom/my-project.md`).

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

The **use command** guides the agent to review and update the project's AI agent instructions:

- Lists available fragments and their contents.
- Reviews the current recipe file.
- Suggests which fragments to add or remove based on project needs.
- Edits the recipe and regenerates the instructions file.

The **decompose command** guides the agent to decompose existing documents into instraweave fragments:

- **Header-based splitting** (preferred): uses `instraweave decompose` when consistent headers exist.
- **Semantic splitting** (fallback): when headers are absent or sparse, the agent infers logical topic boundaries from meaning and creates fragments manually.
- **Verbatim constraint**: body text is always copied as-is — no rewrites, paraphrasing, or additions.

## Example

See the [`examples/fragments/`](examples/fragments/) directory for sample fragments.

## License

MIT
