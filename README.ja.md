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
| `claude` | `.claude/CLAUDE.md` |

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

カレントディレクトリにサンプルの `instraweave-recipe.yaml` と `fragments/` ディレクトリを生成します。

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

**decomposeコマンド**は、既存ドキュメント群をinstraweaveフラグメントに分解する手順をエージェントに提示します:

- **ヘッダー分割**（優先）: 一貫したヘッダーがある場合は `instraweave decompose` を使用。
- **セマンティック分割**（フォールバック）: ヘッダーがない・少ない場合は、意味からトピック境界を推論してフラグメントを手動作成。
- **原文保持の制約**: 本文テキストは必ず逐語コピー — 書き換え・言い換え・追記は禁止。

## サンプル

[`examples/fragments/`](examples/fragments/) ディレクトリにサンプルフラグメントがあります。

## ライセンス

MIT
