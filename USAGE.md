# goWorkspace (`gws`) - Usage Guide

This guide covers common workflows, CLI commands, configuration options, and project templates for `gws`.

---

## Workspace Configuration (`.goworkspace.yaml`)

Each project can contain a `.goworkspace.yaml` configuration file at its root.

### Example Configuration:

```yaml
name: goJPeek
template: go
editor:
  type: vscode
terminal:
  type: ghostty
multiplexer:
  type: zellij
panes:
  - name: shell
    command: ""
  - name: test
    command: ""
  - name: git
    command: git status
env:
  PROJECT_ENV: development
```

---

## Global Configuration (`~/.config/goworkspace/config.yaml`)

Global user defaults are stored at `~/.config/goworkspace/config.yaml`.

```yaml
editor: vscode
terminal: ghostty
multiplexer: zellij
default_template: go
```

You can view or update global options using `gws config`:

```bash
# View global configuration
gws config

# Set default global editor to Zed
gws config --set-editor zed

# Set default global multiplexer to zellij
gws config --set-multiplexer zellij
```

---

## Configuration Precedence

When opening a workspace, settings are evaluated in order:

1. **CLI Overrides** (`--editor`, `--no-editor`, `--multiplexer`, `--no-multiplexer`)
2. **Project Workspace Configuration** (`.goworkspace.yaml`)
3. **Global User Configuration** (`~/.config/goworkspace/config.yaml`)
4. **Built-in Fallbacks** (`vscode`, `ghostty`, `zellij`, `go`)

---

## CLI Reference

### 1. `gws init`
Initializes workspace configuration in the current directory. Auto-detects project type and offers interactive defaults.

```bash
# Interactive initialization
gws init

# Non-interactive mode with flag overrides
gws init -y --template go --editor vscode --multiplexer zellij
```

### 2. `gws open`
Opens the project workspace. Traverses parent directories to locate `.goworkspace.yaml`, creates or attaches the Zellij multiplexer session, generates layout panes, and launches the configured editor.

```bash
# Open current workspace
gws open

# Open without launching VS Code
gws open --no-editor

# Open without launching Zellij multiplexer
gws open --no-multiplexer
```

### 3. `gws list`
Lists all registered workspaces and their active session status.

```bash
gws list
```

### 4. `gws status`
Displays status details, configured tools, session state, and defined panes for the current workspace.

```bash
gws status
```

### 5. `gws resume`
Restores/opens the active workspace session (Phase 1 alias for `gws open`).

```bash
gws resume
```

### 6. `gws remove`
Removes workspace configuration file (`.goworkspace.yaml`) and unregisters workspace from state tracking. **Does not delete project source files.**

```bash
# Interactive removal prompt
gws remove

# Non-interactive removal
gws remove -y
```

---

## Built-In Templates

- **`go`**: Editor `vscode`, Panes: `shell` (`""`), `test` (`""`), `git` (`git status`). Auto-detected via `go.mod`.
- **`dotnet`**: Editor `vscode`, Panes: `shell` (`""`), `build` (`""`), `git` (`git status`). Auto-detected via `*.sln` / `*.csproj`.
- **`node`**: Editor `vscode`, Panes: `shell` (`""`), `dev` (`""`), `git` (`git status`). Auto-detected via `package.json`.
- **`rust`**: Editor `vscode`, Panes: `shell` (`""`), `check` (`""`), `git` (`git status`). Auto-detected via `Cargo.toml`.
- **`generic`**: Editor `vscode`, Panes: `shell` (`""`), `git` (`git status`).
