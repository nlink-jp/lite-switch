# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [0.1.3] - 2026-05-23

### Added

- **Pre-built binary releases for the first time.** A new `package`
  target produces zipped binaries for darwin/amd64, darwin/arm64,
  linux/amd64, linux/arm64, and windows/amd64. Previously
  lite-switch was installed via `go install` only. Asset naming:
  `lite-switch-vX.Y.Z-<os>-<arch>.zip`.
- **Darwin builds are Developer ID signed and Apple-notarized.**
  `make package` runs `scripts/codesign-darwin.sh` per darwin
  binary and `scripts/notarize-darwin.sh` per darwin zip,
  following the org-wide convention in `nlink-jp/.github`
  CONVENTIONS.md §Code Signing. End users on macOS no longer
  need to bypass Gatekeeper with right-click → Open; local
  Dropbox-synced (FileProvider-managed) install paths no longer
  SIGKILL the binary on launch.

No behaviour change to the binary itself — feature-wise this is
identical to v0.1.2.

## [0.1.2] - 2026-03-31

### Fixed
- Skip config file permission check on Windows/NTFS (always reports 0666)

## [0.1.1] - 2026-03-27

### Security

- Added config file permission check: warns to stderr and suggests `chmod 600`
  when the config file is readable by group or others (`perm & 0077 != 0`).


## [0.1.0] - 2026-03-27

### Added

- Initial release.
- `lite-switch`: reads free-form text from stdin and writes the best-matching tag to stdout.
- Tool-calling classification with JSON and plain-text fallbacks for broad LLM compatibility.
- Nonce-wrapped user input to prevent prompt injection.
- Two-file configuration: `config.toml` (TOML, API settings) + `switches.yaml` (YAML, classification definitions).
- Environment variable overrides: `LITE_SWITCH_BASE_URL`, `LITE_SWITCH_API_KEY`, `LITE_SWITCH_MODEL`.
- Exponential backoff retry on transient errors and rate limiting.


[0.1.1]: https://github.com/nlink-jp/lite-switch/releases/tag/v0.1.1
[0.1.0]: https://github.com/nlink-jp/lite-switch/releases/tag/v0.1.0
