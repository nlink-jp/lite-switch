# プロジェクト構成

```
lite-switch/
├── main.go                      # エントリポイント: フラグ解析、stdin 読み込み、結線
├── go.mod / go.sum
├── Makefile
├── config.example.toml          # システム設定のテンプレート（TOML）
├── switches.example.yaml        # スイッチ定義のテンプレート（YAML）
├── AGENTS.md                    # AI コーディングエージェント向けのコンテキスト
├── RULES.md                     # プロジェクトルール
├── CHANGELOG.md
├── README.md / README.ja.md
├── internal/
│   ├── config/
│   │   ├── config.go            # 設定構造体、Load()、デフォルト値、環境変数オーバーライド
│   │   └── config_test.go
│   ├── llm/
│   │   ├── client.go            # リトライとバックオフを備えた HTTP クライアント
│   │   ├── prompt.go            # システムプロンプトビルダー、入力ラッパー
│   │   ├── client_test.go
│   │   └── prompt_test.go
│   └── classifier/
│       ├── classifier.go        # Classify()、ツール構築、タグ抽出
│       └── classifier_test.go
├── docs/
│   ├── en/                      # 英語ドキュメント（一次言語、言語サフィックスなし）
│   │   ├── design/
│   │   │   └── overview.md
│   │   ├── dependencies.md
│   │   ├── setup.md
│   │   ├── structure.md
│   │   └── verification.md
│   └── ja/                      # 日本語訳（`.ja.md` サフィックス）
│       ├── design/
│       │   └── overview.ja.md
│       ├── dependencies.ja.md
│       ├── setup.ja.md
│       ├── structure.ja.md
│       └── verification.ja.md
└── scripts/
    └── hooks/
        ├── pre-commit           # vet + lint
        └── pre-push             # フルチェック
```
