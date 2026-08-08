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

### Multi-Tab, Subdirectory, & Environment Variables Example:

```yaml
name: my-complex-app
template: go
editor:
  type: vscode
terminal:
  type: ghostty
multiplexer:
  type: zellij
env:
  GLOBAL_ENV: production
tabs:
  - name: Code
    panes:
      - name: shell
        command: ""
      - name: git
        command: git status

  - name: Services
    panes:
      - name: server
        dir: ./cmd/server       # Opens pane inside ./cmd/server
        command: go run main.go
        env:
          PORT: "8080"
          DB_URL: "postgres://localhost:5432/mydb"

  - name: Frontend
    panes:
      - name: web
        dir: ./web              # Opens pane inside ./web
        command: npm run dev
        env:
          VITE_PORT: "3000"
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

### 1. `gws tui` (or `gws`)
Launches the interactive terminal dashboard for navigating, opening, inspecting, and stopping workspaces visually.

```bash
# Launch interactive TUI dashboard
gws tui
# or simply
gws
```

*Keybindings in TUI:*
- **`↑` / `↓`** or **`k` / `j`**: Navigate workspace list
- **`o`**: Open selected workspace
- **`s`**: Inspect status details of selected workspace
- **`x`**: Stop active session of selected workspace
- **`q`** or **`Esc`**: Exit dashboard

### 2. `gws init`
Initializes workspace configuration in the current directory. Auto-detects project type and offers interactive defaults.

```bash
# Interactive initialization
gws init

# Non-interactive mode with flag overrides
gws init -y --template go --editor vscode --multiplexer zellij
```

### 2. `gws open`
Opens the project workspace. You can run it inside the project directory, or pass a path / registered workspace name from anywhere in your terminal.

```bash
# Open workspace in current directory
gws open

# Open workspace by path from anywhere
gws open ~/golang_toolshed/goJPeek

# Open workspace by registered project name from anywhere
gws open goJPeek

# Inject dynamic transient panes without modifying .goworkspace.yaml
gws open --pane "db:docker compose up" --pane "logs:tail -f app.log"

# Open without launching VS Code
gws open --no-editor

# Open without launching Zellij multiplexer
gws open --no-multiplexer
```

### 3. `gws stop`
Stops an active workspace session. Can be run inside the project directory or by workspace name/path from anywhere.

```bash
# Stop current workspace session
gws stop

# Stop workspace session by name from anywhere
gws stop goJPeek
```

### 4. `gws list`
Lists all registered workspaces and their active session status.

```bash
gws list
```

### 5. `gws status`
Displays process tree inspection status for each pane, configured tools, and session state.

```bash
gws status
```
*Example Output:*
```text
Workspace:   goJPeek
Path:        /home/user/golang_toolshed/goJPeek
Template:    go
Editor:      vscode
Terminal:    ghostty
Multiplexer: zellij

Session:     running

PANE            STATUS       COMMAND
------------------------------------------------------------
shell           running      (interactive shell)
test            running      (interactive shell)
git             completed    git status
db              running      docker compose up
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
