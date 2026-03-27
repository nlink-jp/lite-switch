# セットアップガイド

## 前提条件

- Go 1.22 以降
- OpenAI 互換の LLM API（OpenAI、LM Studio、Ollama など）

## インストール

```sh
git clone https://github.com/nlink-jp/lite-switch.git
cd lite-switch
make build
# bin/ を PATH に追加するか、bin/lite-switch を PATH 上のディレクトリにコピーしてください
```

## システム設定

1. 設定ファイルをコピー:

   ```sh
   mkdir -p ~/.config/lite-switch
   cp config.example.toml ~/.config/lite-switch/config.toml
   chmod 600 ~/.config/lite-switch/config.toml
   ```

2. `~/.config/lite-switch/config.toml` を編集:

   ```toml
   [api]
   base_url = "https://api.openai.com"
   api_key  = "sk-..."

   [model]
   name = "gpt-4o-mini"
   ```

3. または環境変数で設定:

   ```sh
   export LITE_SWITCH_BASE_URL="http://localhost:1234"
   export LITE_SWITCH_API_KEY="lm-studio"
   export LITE_SWITCH_MODEL="my-model"
   ```

## スイッチファイル

サンプルをコピーしてカスタマイズ:

```sh
cp switches.example.yaml switches.yaml
```

スイッチファイルは分類オプションを定義します。プロジェクトと一緒にバージョン管理できます。

## Git フックのインストール

```sh
make setup
```

`pre-commit`（vet + lint）と `pre-push`（フルチェック）フックをインストールします。
