# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [0.1.0] - 2026-03-27

### Added

- Initial release.
- `lite-switch`: reads free-form text from stdin and writes the best-matching tag to stdout.
- Tool-calling classification with JSON and plain-text fallbacks for broad LLM compatibility.
- Nonce-wrapped user input to prevent prompt injection.
- Two-file configuration: `config.toml` (TOML, API settings) + `switches.yaml` (YAML, classification definitions).
- Environment variable overrides: `LITE_SWITCH_BASE_URL`, `LITE_SWITCH_API_KEY`, `LITE_SWITCH_MODEL`.
- Exponential backoff retry on transient errors and rate limiting.


[0.1.0]: https://github.com/nlink-jp/lite-switch/releases/tag/v0.1.0
