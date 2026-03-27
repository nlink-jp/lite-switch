# lite-switch

シェルパイプライン向けの自然言語分類器です。
stdin からテキストを読み込み、OpenAI 互換の LLM を使って最も合致するタグを stdout に出力します。

[lite-series](https://github.com/nlink-jp/lite-series) の一部です。

## 特徴

- **パイプライン向け設計** — stdin を読んで1つのタグを stdout に書くだけ。対話 UI なし
- **柔軟なスイッチ定義** — バージョン管理できるシンプルな YAML ファイルで分類オプションを定義
- **幅広い LLM 対応** — ツールコールを使い、JSON・プレーンテキストへのフォールバックを備える
- **プロンプトインジェクション対策** — ユーザー入力はノンス付き XML タグで隔離される
- **リトライ機能** — 一時的なエラーとレートリミットに対して指数バックオフでリトライ

## インストール

```sh
git clone https://github.com/nlink-jp/lite-switch.git
cd lite-switch
make build
# バイナリ: bin/lite-switch
```

## 設定

### システム設定（API 接続情報）

```sh
mkdir -p ~/.config/lite-switch
cp config.example.toml ~/.config/lite-switch/config.toml
chmod 600 ~/.config/lite-switch/config.toml
```

```toml
# ~/.config/lite-switch/config.toml
[api]
base_url = "https://api.openai.com"
api_key  = "sk-..."

[model]
name = "gpt-4o-mini"
```

**優先順位（高い順）:** CLI フラグ → 環境変数 → 設定ファイル → コンパイル時デフォルト

| 環境変数                 | 説明              |
|-------------------------|-------------------|
| `LITE_SWITCH_API_KEY`   | API キー          |
| `LITE_SWITCH_BASE_URL`  | API ベース URL    |
| `LITE_SWITCH_MODEL`     | モデル名          |

### スイッチファイル（分類定義）

```sh
cp switches.example.yaml switches.yaml
```

```yaml
switches:
  - tag: weather
    description: 天気予報や気象情報に関する質問や話題
  - tag: default
    description: 上記のいずれにも当てはまらない話題
```

スイッチファイルはシステム設定とは別ファイルなので、プロジェクトと一緒にバージョン管理できます。

## 使い方

```sh
echo "明日は雨が降りますか？" | lite-switch
# → weather

echo "今何時ですか？" | lite-switch -switches my-switches.yaml
# → time
```

```
フラグ:
  -config   string   システム設定ファイルパス（デフォルト: ~/.config/lite-switch/config.toml）
  -switches string   スイッチ定義ファイルパス（デフォルト: switches.yaml）
  -version          バージョンを表示して終了
```

## ビルド

```sh
make build        # 現在のプラットフォーム → bin/lite-switch
make build-all    # 全5プラットフォーム  → dist/
make check        # vet + lint + test + build + govulncheck
```

## ドキュメント

- [セットアップガイド](docs/ja/setup.md)
- [設計概要](docs/ja/design/overview.md)
- [English README](README.md)
