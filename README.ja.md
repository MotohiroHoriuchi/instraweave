# instraweave

再利用可能なルールフラグメントからAIエージェント用instructionsファイルを組み立てるCLIツール。

[English README](README.md)

## 概要

`instraweave` は、AIコーディングエージェント（GitHub Copilot、Claudeなど）のinstructionsを、カテゴリ別に整理されたMarkdownフラグメントとして管理できるツールです。YAMLレシピファイルで使用するフラグメントを選択し、1つのinstructionsファイルに結合します。

## インストール

```bash
go install github.com/MotohiroHoriuchi/instraweave@latest
```

ソースからビルドする場合:

```bash
git clone https://github.com/MotohiroHoriuchi/instraweave.git
cd instraweave
go build -o instraweave .
```

## クイックスタート

```bash
# 1. サンプルのレシピファイルとfragmentsディレクトリを生成
instraweave init

# 2. 利用可能なフラグメント一覧を表示
instraweave list

# 3. 結合結果をプレビュー
instraweave generate --dry-run

# 4. instructionsファイルを生成
instraweave generate
```

## レシピファイル

`instraweave` はYAMLレシピファイル（`instraweave-recipe.yaml`）で生成内容を定義します:

```yaml
target: copilot              # copilot | claude
output: ""                   # 空の場合、ターゲットのデフォルトパスを使用
fragments_dir: ./fragments   # フラグメント格納ディレクトリ（デフォルト: ./fragments）
fragments:
  - standard/go
  - standard/testing
  - standard/security
  - custom/our-api-convention
```

### 対応ターゲット

| ターゲット | デフォルト出力先 |
|-----------|----------------|
| `copilot` | `.github/copilot-instructions.md` |
| `claude` | `CLAUDE.md` |

## レシピ継承

`extends` フィールドを使うと、別のレシピを基底として継承できます。チームやプロジェクト間でベースとなるフラグメントセットを共有しながら、各レイヤーでカスタマイズできます。

### 基本構文

**派生レシピ**（`extends` あり）はプレーン名の代わりにオペレーションを使います:

```yaml
extends: ../base/recipe.yaml   # 相対パスまたは絶対パス

target: claude
fragments_dir: ./fragments

fragments:
  - add: standard/go             # リスト末尾に追加
  - remove: standard/code-review # リストから削除
  - override: standard/security  # このレシピのバージョンで上書き
```

**ルートレシピ**（`extends` なし）はフラグメント名をプレーンに列挙します:

```yaml
target: claude
fragments_dir: ./fragments
fragments:
  - standard/security
  - standard/git-convention
  - standard/code-review
```

### フラグメントオペレーション

| オペレーション | 構文 | 動作 |
|-------------|------|------|
| *(省略)* | `- category/name` | ルートレシピのみ使用可。派生レシピで使うとエラー。 |
| `add` | `- add: category/name` | リスト末尾に追加。すでに存在する場合はエラー。 |
| `remove` | `- remove: category/name` | リストから削除。存在しない場合はエラー。 |
| `override` | `- override: category/name` | このレシピの `fragments_dir` で解決するよう変更。存在しない場合はエラー。 |

### 継承チェーン

`extends` は再帰的に解決されます。オペレーションはルートから派生方向へ順に適用されます（後が勝つ）:

```
org/recipe.yaml          ← ルート（プレーン指定）
  └─ team/recipe.yaml    ← Go追加、code-review削除
       └─ project/recipe.yaml  ← security上書き、db-convention追加
```

各fragmentは**最後に操作したレシピの `fragments_dir`** から読み込まれます:

- ルートのプレーン指定 → ルートの `fragments_dir`
- `add` したfragment → そのaddを書いたレシピの `fragments_dir`
- `override` したfragment → そのoverrideを書いたレシピの `fragments_dir`

`target` と `output` も継承され、派生レシピに記述があれば親の値を上書きします。

### ディレクトリ構成例

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
        └── security.md  # org版をoverride
```

### dry-run 出力

`instraweave generate --dry-run` を実行すると、継承チェーンと各fragmentの解決元が表示されます:

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

## フラグメント構成

フラグメントはサブディレクトリに整理されたMarkdownファイルです:

```
fragments/
├── standard/          # 共通・再利用可能なルール
│   ├── go.md
│   ├── testing.md
│   └── security.md
└── custom/            # プロジェクト固有のルール
    └── our-api-convention.md
```

レシピ内のフラグメント名は、`fragments_dir` 配下のファイルパス（`.md` 拡張子なし）に対応します。

## コマンド

### `instraweave init`

カレントディレクトリにサンプルの `instraweave-recipe.yaml` と `fragments/` ディレクトリ（`fragments/standard/go.md` および `fragments/custom/my-project.md` を含む）を生成します。

```bash
instraweave init
```

### `instraweave list`

指定ディレクトリ内の利用可能なフラグメント一覧を表示します。

```bash
instraweave list
instraweave list --dir ./my-fragments
```

| フラグ | 短縮形 | デフォルト | 説明 |
|-------|--------|-----------|------|
| `--dir` | `-d` | `./fragments` | フラグメントディレクトリ |
| `--verbose` | `-v` | `false` | フラグメントの内容を表示 |

### `instraweave show`

1つ以上のフラグメントの内容を表示します。AIエージェントがレシピを組み立てる前にフラグメントを確認する用途に適しています。

```bash
instraweave show standard/go
instraweave show standard/go standard/testing
instraweave show --all
instraweave show --all --dir ./my-fragments
```

| フラグ | 短縮形 | デフォルト | 説明 |
|-------|--------|-----------|------|
| `--dir` | `-d` | `./fragments` | フラグメントディレクトリ |
| `--all` | | `false` | 全フラグメントを表示 |

### `instraweave generate`

レシピファイルを読み込み、フラグメントを結合してinstructionsファイルを出力します。

```bash
instraweave generate
instraweave generate --recipe ./my-recipe.yaml
instraweave generate --dry-run
```

| フラグ | 短縮形 | デフォルト | 説明 |
|-------|--------|-----------|------|
| `--recipe` | `-r` | `./instraweave-recipe.yaml` | レシピファイルのパス |
| `--dry-run` | | `false` | ファイルに書き込まず標準出力に出力 |

### `instraweave decompose`

Markdownファイルをヘッダーレベルで分割してフラグメントファイルを生成します。

```bash
instraweave decompose --file CLAUDE.md
instraweave decompose --file docs/guide.md --level 1 --dir ./fragments/custom/
```

| フラグ | 短縮形 | デフォルト | 説明 |
|-------|--------|-----------|------|
| `--file` | `-f` | *(必須)* | 分解対象のMarkdownファイル |
| `--level` | `-l` | `2` | 分割に使うヘッダーレベル（1〜6） |
| `--dir` | `-d` | `./fragments` | フラグメントファイルの出力ディレクトリ |

### `instraweave agent`

AIエージェントが instraweave を直接操作できるように、プロンプト/コマンドファイルをインストールします。

```bash
instraweave agent --target claude
instraweave agent --target copilot
instraweave agent --target claude --force   # 既存ファイルを上書き
```

| フラグ | 短縮形 | デフォルト | 説明 |
|-------|--------|-----------|------|
| `--target` | `-t` | *(必須)* | エージェントの種類: `claude` または `copilot` |
| `--force` | | `false` | 既存ファイルを上書きする |

**インストールされるファイル:**

| ターゲット | useコマンド | decomposeコマンド |
|-----------|------------|-----------------|
| `claude` | `.claude/commands/instraweave.md` | `.claude/commands/instraweave-decompose.md` |
| `copilot` | `.github/prompts/instraweave.prompt.md` | `.github/prompts/instraweave-decompose.prompt.md` |

**useコマンド**は、プロジェクトのAIエージェント用instructionsを確認・更新する手順をエージェントに提示します:

- 利用可能なフラグメントとその内容を一覧表示する。
- 現在のレシピファイルを確認する。
- プロジェクトの状況に合わせて追加・削除するフラグメントを提案する。
- レシピを更新し、instructionsファイルを再生成する。

**decomposeコマンド**は、既存ドキュメント群をinstraweaveフラグメントに分解する手順をエージェントに提示します:

- **ヘッダー分割**（優先）: 一貫したヘッダーがある場合は `instraweave decompose` を使用。
- **セマンティック分割**（フォールバック）: ヘッダーがない・少ない場合は、意味からトピック境界を推論してフラグメントを手動作成。
- **原文保持の制約**: 本文テキストは必ず逐語コピー — 書き換え・言い換え・追記は禁止。

## サンプル

[`examples/fragments/`](examples/fragments/) ディレクトリにサンプルフラグメントがあります。

## ライセンス

MIT
