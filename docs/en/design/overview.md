# Design Overview

## Purpose

lite-switch classifies free-form text into a predefined set of tags using an LLM.
It is designed as a Unix filter: reads stdin, writes one tag to stdout.

## Configuration

Values are resolved in priority order (highest first):

1. CLI flags (`-config`, `-switches`)
2. Environment variables (`LITE_SWITCH_API_KEY`, `LITE_SWITCH_BASE_URL`, `LITE_SWITCH_MODEL`)
3. Config file (`~/.config/lite-switch/config.toml`)
4. Compiled-in defaults

### Two-file design

| File | Format | Contents | Location |
|------|--------|----------|----------|
| `config.toml` | TOML | API credentials, model, timeouts | `~/.config/lite-switch/` |
| `switches.yaml` | YAML | Tag/description pairs | Project directory |

The files are kept separate because they have different security profiles and change rates.
The switches file contains no secrets and can be version-controlled freely.
YAML is used for the switches file because its list syntax is more readable than TOML
array-of-tables for this kind of data.

## Classification pipeline

```
stdin
  └─ readInput()           — size-limited read, trim whitespace
       └─ WrapUserInput()  — nonce-tagged XML wrapper (prompt injection protection)
            └─ Classify()  — builds request, calls LLM, extracts tag
                 └─ extractTag()  — 4-strategy fallback:
                      1. tool call result
                      2. JSON {"tag": "..."} in content
                      3. known tag string in content
                      4. last switch (catch-all default)
  └─ stdout: one tag
```

## LLM interaction

The classifier uses the function-calling / tool-use API (`tools` + `tool_choice: "required"`).
This forces the model to return a structured response and avoids free-text parsing in most cases.

For models that do not support tool calling, the fallback chain (JSON → text → default) ensures
a usable result is always returned.

## Endpoint URL handling

`base_url` accepts both `http://host` and `http://host/v1` forms.
The client appends `/chat/completions` or `/v1/chat/completions` as appropriate.

## Retry policy

Transient errors (network failures, HTTP 429, HTTP 5xx) are retried with exponential backoff.
Non-retryable errors (4xx other than 429) fail immediately.
