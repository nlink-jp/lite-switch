# lite-switch

A natural language classifier for shell pipelines.
Reads free-form text from stdin and writes the best-matching tag to stdout,
using an OpenAI-compatible LLM.

Part of [lite-series](https://github.com/nlink-jp/lite-series).

## Features

- **Pipeline-friendly** — reads stdin, writes a single tag to stdout; no interactive UI
- **Configurable switches** — define classification options in a plain YAML file you can version-control
- **Broad LLM compatibility** — uses tool calling with JSON and plain-text fallbacks
- **Prompt-injection protection** — user input is isolated in a nonce-tagged XML wrapper
- **Retry logic** — exponential backoff on transient errors and rate limiting

## Installation

```sh
git clone https://github.com/nlink-jp/lite-switch.git
cd lite-switch
make build
# binary: dist/lite-switch
```

## Configuration

### System config (API settings)

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

**Priority order (highest first):** CLI flags → environment variables → config file → compiled-in defaults

| Environment variable    | Description        |
|-------------------------|--------------------|
| `LITE_SWITCH_API_KEY`   | API key            |
| `LITE_SWITCH_BASE_URL`  | API base URL       |
| `LITE_SWITCH_MODEL`     | Model name         |

### Switches file (classification definitions)

```sh
cp switches.example.yaml switches.yaml
```

```yaml
switches:
  - tag: weather
    description: Questions or topics about weather forecasts
  - tag: default
    description: Anything that does not match the above switches
```

The switches file is separate from the system config so it can be version-controlled alongside your project.

## Usage

```sh
echo "Will it rain tomorrow?" | lite-switch
# → weather

echo "What time is it?" | lite-switch -switches my-switches.yaml
# → time
```

```
Flags:
  -config   string   system config file path (default: ~/.config/lite-switch/config.toml)
  -switches string   switches definition file path (default: switches.yaml)
  -version          print version and exit
```

## Building

```sh
make build        # current platform → dist/lite-switch
make build-all    # all 5 platforms  → dist/
make check        # vet + lint + test + build + govulncheck
```

## Documentation

- [Setup guide](docs/setup.md)
- [Design overview](docs/design/overview.md)
- [日本語 README](README.ja.md)
