Decompose an existing instructions markdown file into instrweave fragments.

Usage: /instrweave-decompose <file-path> [header-level] [output-dir]

Steps:
1. Run `instrweave decompose --file <file-path> --level <n> --dir ./fragments/custom/`
   - `--level`: header level used as split boundary (default: 2, i.e. ## headers)
   - `--dir`: output directory for generated fragment files (default: ./fragments)
2. Review the generated fragment files.
3. Rename files or adjust content as necessary.
4. Add the new fragments to `instrweave-recipe.yaml`.
5. Run `instrweave generate --dry-run` to verify the result.
