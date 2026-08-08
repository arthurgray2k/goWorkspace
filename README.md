# goWorkspace (`gws`)

`goWorkspace` (`gws`) is a lightweight Linux developer workspace manager written in Go.

> "Describe my preferred development environment once, then recreate it with one command."

`gws` orchestrates your existing developer tools (such as **Zellij**, **VS Code**, **Zed**, and **Ghostty**) into reproducible, project-specific workspace configurations without replacing any of them.

---

## Key Features

- **Portable Workspace Configuration**: Managed via a clean, human-editable `.goworkspace.yaml` file inside your project.
- **Deterministic Project Detection**: Auto-detects Go (`go.mod`), Node (`package.json`), Rust (`Cargo.toml`), and .NET (`*.sln`/`*.csproj`) project structures.
- **Zellij Layout Orchestration**: Generates custom, non-destructive Zellij KDL layouts for tabs and panes without modifying global Zellij configuration.
- **Nested-Session Guard**: Automatically detects active Zellij environments (`$ZELLIJ`) to prevent nested sessions.
- **Editor & Terminal Integration**: Detached execution of VS Code (`code`), Zed (`zed`), or Ghostty (`ghostty`) via standard `PATH` lookups.
- **Desktop-Environment Agnostic**: Works headlessly, over SSH, X11, or Wayland across any Linux desktop environment.
- **Zero Data Loss Guarantee**: `gws remove` removes workspace configuration/registration but **NEVER** touches project files.

---

## Installation & Build

Requires Go 1.26 or later.

```bash
# Clone repository
git clone git@github.com:arthurgray2k/goWorkspace.git
cd goWorkspace

# Build binary
go build -o gws ./cmd/gws

# Optional: Install to PATH
sudo mv gws /usr/local/bin/
```

---

## Quick Start

```bash
cd ~/golang_toolshed/goJPeek

# Initialize workspace configuration (auto-detects Go template)
gws init

# Open editor, terminal, and multiplexer panes
gws open

# Check workspace session status
gws status

# Stop active workspace session
gws stop

# List all registered workspaces
gws list
```

---

## Documentation

- See [USAGE.md](USAGE.md) for detailed workflow examples, configuration format specification, and CLI flags.

---

## License

See [LICENSE](LICENSE) for license details.
