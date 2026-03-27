# Project Structure

```
lite-switch/
├── main.go                      # Entry point: flag parsing, stdin reading, wiring
├── go.mod / go.sum
├── Makefile
├── config.example.toml          # System config template (TOML)
├── switches.example.yaml        # Switches definition template (YAML)
├── AGENTS.md                    # Context for AI coding agents
├── RULES.md                     # Project rules
├── CHANGELOG.md
├── README.md / README.ja.md
├── internal/
│   ├── config/
│   │   ├── config.go            # Config structs, Load(), defaults, env overrides
│   │   └── config_test.go
│   ├── llm/
│   │   ├── client.go            # HTTP client with retry and backoff
│   │   ├── prompt.go            # System prompt builder, input wrapper
│   │   ├── client_test.go
│   │   └── prompt_test.go
│   └── classifier/
│       ├── classifier.go        # Classify(), tool building, tag extraction
│       └── classifier_test.go
├── docs/
│   ├── design/
│   │   └── overview.md
│   ├── ja/                      # Japanese translations
│   │   ├── design/
│   │   │   └── overview.md
│   │   └── setup.md
│   ├── dependencies.md
│   ├── setup.md
│   └── structure.md
└── scripts/
    └── hooks/
        ├── pre-commit           # vet + lint
        └── pre-push             # full check
```
