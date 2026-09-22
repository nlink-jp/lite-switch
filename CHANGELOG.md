# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Fixed

- **`make verify-release` now fails closed.** Its last block chained unzip, the
  packaged binary's `--version` and `spctl` with `&&` and ended the whole chain
  in `|| true`, so a zip that did not unpack or a binary that did not run exited
  0 and the upload proceeded. Each step is now judged on its own, the packaged
  binary's `--version` must contain the tag being released, and only the
  informational `spctl` line may be ignored. Matches the org template
  (CONVENTIONS.md §Code Signing → Verifying a release).
- **The Linux archives no longer carry macOS file metadata.** macOS `tar` wrote
  each bundled file's extended attributes (`com.apple.provenance`, and a Dropbox
  attribute where the tree is synced) into the `.tar.gz` twice: as AppleDouble
  `._` members, which GNU tar extracts as stray `._<name>` files beside the real
  ones, and as `LIBARCHIVE.xattr.*` / `SCHILY.xattr.*` pax headers, which it
  reports as unknown keywords. `make package` now archives with
  `COPYFILE_DISABLE=1 tar --no-xattrs`; each setting stops one of the two.
  Archives already published still carry them; the files themselves are
  unaffected.

### Internal

- `make verify-release` also judges each Linux archive: no AppleDouble or other
  macOS metadata members — listed with `--options 'tar:!mac-ext'`, because a
  plain macOS listing folds `._` members away — no extended attributes as pax
  headers, and exactly the canonical binary, `README.md` and `LICENSE`, compared
  in the C locale.
- The Linux-archive check in `make verify-release` reads each archive's pax
  headers with Python's `tarfile` instead of grepping the decompressed stream,
  which also matched file text that names the keywords (a bundled CHANGELOG,
  for one).

## [0.2.0] - 2026-07-12

### Added

- **`LICENSE` file (MIT).** The repository previously had no license file;
  it is now MIT-licensed and the license is bundled in every release archive.

### Removed

- **darwin/amd64 (Intel) pre-built binary.** macOS releases now ship
  **arm64 only**, per the org-wide policy (darwin is Apple-Silicon only; no
  universal binaries). Intel Mac users can build from source.

### Changed

- **Linux release archives are now `.tar.gz`** (darwin/windows remain `.zip`),
  per `nlink-jp/.github` CONVENTIONS.md §Release Archive Standard. Archives
  now bundle `README.md` + `LICENSE` alongside the canonical binary.
- **darwin code-signature identifier** is now the canonical `lite-switch`
  (was `lite-switch-darwin-arm64`), set via `codesign -i` so it stays stable
  after the archived binary is renamed to its canonical name.

No change to the binary's behaviour — a packaging / build-config release.

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
