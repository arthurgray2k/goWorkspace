package template

import (
	"fmt"
	"strings"

	"github.com/arthurgray2k/goWorkspace/internal/config"
)

// Template defines a project template preset.
type Template struct {
	Name        string
	Editor      string
	Terminal    string
	Multiplexer string
	Panes       []config.PaneConfig
}

var builtinRegistry = map[string]Template{
	"go": {
		Name:        "go",
		Editor:      "vscode",
		Terminal:    "ghostty",
		Multiplexer: "zellij",
		Panes: []config.PaneConfig{
			{Name: "shell", Command: ""},
			{Name: "test", Command: ""},
			{Name: "git", Command: "git status"},
		},
	},
	"dotnet": {
		Name:        "dotnet",
		Editor:      "vscode",
		Terminal:    "ghostty",
		Multiplexer: "zellij",
		Panes: []config.PaneConfig{
			{Name: "shell", Command: ""},
			{Name: "build", Command: ""},
			{Name: "git", Command: "git status"},
		},
	},
	"node": {
		Name:        "node",
		Editor:      "vscode",
		Terminal:    "ghostty",
		Multiplexer: "zellij",
		Panes: []config.PaneConfig{
			{Name: "shell", Command: ""},
			{Name: "dev", Command: ""},
			{Name: "git", Command: "git status"},
		},
	},
	"rust": {
		Name:        "rust",
		Editor:      "vscode",
		Terminal:    "ghostty",
		Multiplexer: "zellij",
		Panes: []config.PaneConfig{
			{Name: "shell", Command: ""},
			{Name: "check", Command: ""},
			{Name: "git", Command: "git status"},
		},
	},
	"generic": {
		Name:        "generic",
		Editor:      "vscode",
		Terminal:    "ghostty",
		Multiplexer: "zellij",
		Panes: []config.PaneConfig{
			{Name: "shell", Command: ""},
			{Name: "git", Command: "git status"},
		},
	},
}

// GetTemplate retrieves a template by name (case-insensitive).
func GetTemplate(name string) (Template, bool) {
	tmpl, ok := builtinRegistry[strings.ToLower(name)]
	return tmpl, ok
}

// ListBuiltinNames returns the names of all registered built-in templates.
func ListBuiltinNames() []string {
	return []string{"generic", "go", "dotnet", "node", "rust"}
}

// ApplyToWorkspaceConfig initializes or updates a WorkspaceConfig based on a template.
func ApplyToWorkspaceConfig(name string, tmplName string, dirName string) (*config.WorkspaceConfig, error) {
	tmpl, ok := GetTemplate(tmplName)
	if !ok {
		return nil, fmt.Errorf("unknown template %q (available: %s)", tmplName, strings.Join(ListBuiltinNames(), ", "))
	}

	if name == "" {
		name = dirName
	}

	return &config.WorkspaceConfig{
		Name:        name,
		Template:    tmpl.Name,
		Editor:      config.ToolConfig{Type: tmpl.Editor},
		Terminal:    config.ToolConfig{Type: tmpl.Terminal},
		Multiplexer: config.ToolConfig{Type: tmpl.Multiplexer},
		Panes:       append([]config.PaneConfig(nil), tmpl.Panes...),
	}, nil
}
