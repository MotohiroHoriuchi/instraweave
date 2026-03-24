---
mode: agent
description: Decompose an existing instructions file into instrweave fragments
---

Decompose an existing instructions markdown file into instrweave fragments.

Steps:
1. Identify the target markdown file (e.g. `CLAUDE.md`, `.github/copilot-instructions.md`).
2. Choose an appropriate header level for splitting (e.g. 2 for ## headers).
3. Run `instrweave decompose --file <path> --level <n> --dir ./fragments/custom/`
   - `--level`: header level used as split boundary (default: 2)
   - `--dir`: output directory for generated fragment files (default: ./fragments)
4. Review the generated fragment files.
5. Rename files or adjust content as necessary.
6. Add the new fragments to `instrweave-recipe.yaml`.
7. Run `instrweave generate --dry-run` to verify the result.
