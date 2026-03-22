# instrweave

再利用可能なルールフラグメントからAIエージェント用instructionsファイルを組み立てるCLIツール。

[English README](README.md)

## 概要

`instrweave` は、AIコーディングエージェント（GitHub Copilot、Claudeなど）のinstructionsを、カテゴリ別に整理されたMarkdownフラグメントとして管理できるツールです。YAMLレシピファイルで使用するフラグメントを選択し、1つのinstructionsファイルに結合します。

## インストール

```bash
go install github.com/motohirohoriuchi/instrweave@latest
```

ソースからビルドする場合:

```bash
git clone https://github.com/motohirohoriuchi/instrweave.git
cd instrweave
go build -o instrweave .
```

## クイックスタート

```bash
# 1. サンプルのレシピファイルとfragmentsディレクトリを生成
instrweave init

# 2. 利用可能なフラグメント一覧を表示
instrweave list

# 3. 結合結果をプレビュー
instrweave generate --dry-run

# 4. instructionsファイルを生成
instrweave generate
```

## レシピファイル

`instrweave` はYAMLレシピファイル（`instrweave-recipe.yaml`）で生成内容を定義します:

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

### `instrweave init`

カレントディレクトリにサンプルの `instrweave-recipe.yaml` と `fragments/` ディレクトリを生成します。

```bash
instrweave init
```

### `instrweave list`

指定ディレクトリ内の利用可能なフラグメント一覧を表示します。

```bash
instrweave list
instrweave list --dir ./my-fragments
```

| フラグ | 短縮形 | デフォルト | 説明 |
|-------|--------|-----------|------|
| `--dir` | `-d` | `./fragments` | フラグメントディレクトリ |
| `--verbose` | `-v` | `false` | フラグメントの内容を表示 |

### `instrweave show`

1つ以上のフラグメントの内容を表示します。AIエージェントがレシピを組み立てる前にフラグメントを確認する用途に適しています。

```bash
instrweave show standard/go
instrweave show standard/go standard/testing
instrweave show --all
instrweave show --all --dir ./my-fragments
```

| フラグ | 短縮形 | デフォルト | 説明 |
|-------|--------|-----------|------|
| `--dir` | `-d` | `./fragments` | フラグメントディレクトリ |
| `--all` | | `false` | 全フラグメントを表示 |

### `instrweave generate`

レシピファイルを読み込み、フラグメントを結合してinstructionsファイルを出力します。

```bash
instrweave generate
instrweave generate --recipe ./my-recipe.yaml
instrweave generate --dry-run
```

| フラグ | 短縮形 | デフォルト | 説明 |
|-------|--------|-----------|------|
| `--recipe` | `-r` | `./instrweave-recipe.yaml` | レシピファイルのパス |
| `--dry-run` | | `false` | ファイルに書き込まず標準出力に出力 |

## サンプル

[`examples/fragments/`](examples/fragments/) ディレクトリにサンプルフラグメントがあります。

## ライセンス

MIT
