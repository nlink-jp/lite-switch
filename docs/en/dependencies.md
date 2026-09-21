# Dependencies

This document records the rationale for each third-party dependency (see RULES.md §18).

## github.com/BurntSushi/toml

- **Purpose**: Parse the TOML system config file (`config.toml`).
- **Why not in-house**: TOML parsing requires non-trivial grammar handling; a battle-tested library reduces risk.
- **License**: MIT

## gopkg.in/yaml.v3

- **Purpose**: Parse the YAML switches definition file (`switches.yaml`).
- **Why YAML for switches**: The switches file is user-facing classification data (a list of tag/description pairs). YAML's list syntax (`- tag: foo`) is more readable and concise than TOML's array-of-tables (`[[switches]]`) for this use case.
- **Why not in-house**: Same rationale as TOML — grammar complexity warrants a library.
- **License**: MIT and Apache 2.0
