package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const WorkspaceConfigFilename = ".goworkspace.yaml"

// ToolConfig specifies editor, terminal, or multiplexer settings.
type ToolConfig struct {
	Type string `yaml:"type"`
}

// PaneConfig describes a multiplexer pane configuration.
type PaneConfig struct {
	Name    string `yaml:"name"`
	Command string `yaml:"command"`
}

// WorkspaceConfig represents the .goworkspace.yaml file structure.
type WorkspaceConfig struct {
	Name        string            `yaml:"name"`
	Template    string            `yaml:"template"`
	Editor      ToolConfig        `yaml:"editor"`
	Terminal    ToolConfig        `yaml:"terminal"`
	Multiplexer ToolConfig        `yaml:"multiplexer"`
	Panes       []PaneConfig      `yaml:"panes"`
	Env         map[string]string `yaml:"env,omitempty"`
}

// GlobalConfig represents ~/.config/goworkspace/config.yaml.
type GlobalConfig struct {
	Editor          string `yaml:"editor,omitempty"`
	Terminal        string `yaml:"terminal,omitempty"`
	Multiplexer     string `yaml:"multiplexer,omitempty"`
	DefaultTemplate string `yaml:"default_template,omitempty"`
}

// ConfigFlags captures temporary command-line overrides.
type ConfigFlags struct {
	Editor        string
	NoEditor      bool
	Terminal      string
	Multiplexer   string
	NoMultiplexer bool
	Template      string
}

// ResolvedConfig contains the final combined configuration after applying precedence rules.
type ResolvedConfig struct {
	Name        string
	ProjectDir  string
	Template    string
	Editor      string
	Terminal    string
	Multiplexer string
	Panes       []PaneConfig
	Env         map[string]string
}

// DefaultGlobalConfig returns fallback defaults.
func DefaultGlobalConfig() *GlobalConfig {
	return &GlobalConfig{
		Editor:          "vscode",
		Terminal:        "ghostty",
		Multiplexer:     "zellij",
		DefaultTemplate: "go",
	}
}

// GlobalConfigPath returns the absolute path to ~/.config/goworkspace/config.yaml.
func GlobalConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home dir: %w", err)
	}
	return filepath.Join(home, ".config", "goworkspace", "config.yaml"), nil
}

// LoadGlobalConfig loads global configuration or returns default if file does not exist.
func LoadGlobalConfig() (*GlobalConfig, error) {
	path, err := GlobalConfigPath()
	if err != nil {
		return DefaultGlobalConfig(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultGlobalConfig(), nil
		}
		return nil, fmt.Errorf("error reading global config %s: %w", path, err)
	}

	cfg := DefaultGlobalConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("error parsing global config %s: %w", path, err)
	}

	return cfg, nil
}

// SaveGlobalConfig saves global configuration to ~/.config/goworkspace/config.yaml.
func SaveGlobalConfig(cfg *GlobalConfig) error {
	path, err := GlobalConfigPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal global config: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}

// LoadWorkspaceConfig reads .goworkspace.yaml from target directory.
func LoadWorkspaceConfig(dir string) (*WorkspaceConfig, error) {
	path := filepath.Join(dir, WorkspaceConfigFilename)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNoWorkspaceConfig
		}
		return nil, fmt.Errorf("error reading %s: %w", path, err)
	}

	var cfg WorkspaceConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("error parsing %s: %w", path, err)
	}

	return &cfg, nil
}

// SaveWorkspaceConfig writes .goworkspace.yaml to target directory.
func SaveWorkspaceConfig(dir string, cfg *WorkspaceConfig) error {
	path := filepath.Join(dir, WorkspaceConfigFilename)
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal workspace config: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

var ErrNoWorkspaceConfig = errors.New("no .goworkspace.yaml found")

// Resolve evaluates configuration according to precedence:
// CLI Flags > Workspace Config > Global Config > Default Fallbacks
func Resolve(ws *WorkspaceConfig, global *GlobalConfig, flags ConfigFlags, dir string) ResolvedConfig {
	if global == nil {
		global = DefaultGlobalConfig()
	}

	res := ResolvedConfig{
		ProjectDir: dir,
		Env:        make(map[string]string),
	}

	if ws != nil {
		res.Name = ws.Name
		res.Template = ws.Template
		res.Editor = ws.Editor.Type
		res.Terminal = ws.Terminal.Type
		res.Multiplexer = ws.Multiplexer.Type
		res.Panes = append(res.Panes, ws.Panes...)
		for k, v := range ws.Env {
			res.Env[k] = v
		}
	}

	if res.Name == "" {
		res.Name = filepath.Base(dir)
	}

	// Apply global defaults if field is empty
	if res.Template == "" {
		res.Template = global.DefaultTemplate
	}
	if res.Editor == "" {
		res.Editor = global.Editor
	}
	if res.Terminal == "" {
		res.Terminal = global.Terminal
	}
	if res.Multiplexer == "" {
		res.Multiplexer = global.Multiplexer
	}

	// Apply final fallback defaults if still empty
	if res.Template == "" {
		res.Template = "generic"
	}
	if res.Editor == "" {
		res.Editor = "vscode"
	}
	if res.Terminal == "" {
		res.Terminal = "ghostty"
	}
	if res.Multiplexer == "" {
		res.Multiplexer = "zellij"
	}

	// Apply CLI Flags overrides (highest priority)
	if flags.Template != "" {
		res.Template = flags.Template
	}
	if flags.NoEditor {
		res.Editor = "none"
	} else if flags.Editor != "" {
		res.Editor = flags.Editor
	}

	if flags.Terminal != "" {
		res.Terminal = flags.Terminal
	}

	if flags.NoMultiplexer {
		res.Multiplexer = "none"
	} else if flags.Multiplexer != "" {
		res.Multiplexer = flags.Multiplexer
	}

	return res
}
