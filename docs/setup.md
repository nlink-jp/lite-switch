# Setup Guide

## Prerequisites

- Go 1.22 or later
- An OpenAI-compatible LLM API (OpenAI, LM Studio, Ollama, etc.)

## Installation

```sh
git clone https://github.com/nlink-jp/lite-switch.git
cd lite-switch
make build
# Add bin/ to your PATH or copy bin/lite-switch somewhere on your PATH
```

## System config

1. Copy the example config:

   ```sh
   mkdir -p ~/.config/lite-switch
   cp config.example.toml ~/.config/lite-switch/config.toml
   chmod 600 ~/.config/lite-switch/config.toml
   ```

2. Edit `~/.config/lite-switch/config.toml`:

   ```toml
   [api]
   base_url = "https://api.openai.com"
   api_key  = "sk-..."

   [model]
   name = "gpt-4o-mini"
   ```

3. Alternatively, use environment variables:

   ```sh
   export LITE_SWITCH_BASE_URL="http://localhost:1234"
   export LITE_SWITCH_API_KEY="lm-studio"
   export LITE_SWITCH_MODEL="my-model"
   ```

## Switches file

Copy the example and customize:

```sh
cp switches.example.yaml switches.yaml
```

The switches file defines the classification options. It is safe to version-control alongside your project.

## Install Git hooks

```sh
make setup
```

This installs `pre-commit` (vet + lint) and `pre-push` (full check) hooks.
