---
mode: agent
description: Decompose existing documents into instraweave fragments
---

Decompose existing documents into instraweave fragments.

**Constraint: never alter the original text. Copy content verbatim into fragment files — no rewrites, paraphrasing, or additions.**

Steps:
1. Identify target documents:
   - If a path is given, use it; otherwise discover all markdown files in the project
     (e.g. `find . -name "*.md" -not -path "*/node_modules/*" -not -path "*/.git/*"`)
   - Typical candidates: `.github/copilot-instructions.md`, `docs/*.md`, `README.md`, `*.md` in root

2. For each target file, choose a splitting strategy:
   a. **Header-based** (preferred): if the file has consistent headers, determine the
      appropriate level (default: 2, i.e. `##`; use 1 if only `#` sections exist) and run:
      `instraweave decompose --file <path> --level <n> --dir ./fragments/custom/`
   b. **Semantic** (fallback): if the file has no headers or headers are too sparse/deep,
      identify logical topic boundaries from meaning (e.g. a block describing one rule,
      one workflow, or one concept). For each boundary:
      - Insert a `##` header that names the topic (the header itself is new; body text is untouched)
      - Write the resulting section as a fragment file manually

3. Review all generated fragment files:
   - Verify each fragment's body matches the source verbatim
   - Remove or merge duplicate/redundant fragments
   - Rename files to follow the project's naming convention

4. Add the new fragments to `instraweave-recipe.yaml`.

5. Run `instraweave generate --dry-run` to verify the composed output.
