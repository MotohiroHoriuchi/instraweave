Review and update this project's AI agent instructions using instrweave.

Steps:
1. Run `instrweave list --verbose` to see all available fragments and their contents.
2. Run `cat instrweave-recipe.yaml` to review the current recipe.
3. Consider the project's current needs and suggest which fragments to add or remove.
4. Edit `instrweave-recipe.yaml` to reflect the desired changes.
5. Run `instrweave generate --dry-run` to preview the composed output.
6. Run `instrweave generate` to write the instructions file.
