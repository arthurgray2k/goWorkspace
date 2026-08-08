package registry

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/arthurgray2k/goWorkspace/internal/exec"
	"gopkg.in/yaml.v3"
)

// WorkspaceEntry represents a known registered workspace.
type WorkspaceEntry struct {
	Name       string `yaml:"name"`
	Path       string `yaml:"path"`
	Template   string `yaml:"template"`
	LastOpened string `yaml:"last_opened,omitempty"`
}

// WorkspaceRegistry holds all registered workspaces.
type WorkspaceRegistry struct {
	Workspaces []WorkspaceEntry `yaml:"workspaces"`
}

func RegistryPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home dir: %w", err)
	}
	return filepath.Join(home, ".config", "goworkspace", "workspaces.yaml"), nil
}

func LoadRegistry() (*WorkspaceRegistry, error) {
	path, err := RegistryPath()
	if err != nil {
		return &WorkspaceRegistry{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &WorkspaceRegistry{}, nil
		}
		return nil, fmt.Errorf("error reading workspace registry: %w", err)
	}

	var reg WorkspaceRegistry
	if err := yaml.Unmarshal(data, &reg); err != nil {
		return nil, fmt.Errorf("error parsing workspace registry: %w", err)
	}

	return &reg, nil
}

func SaveRegistry(reg *WorkspaceRegistry) error {
	path, err := RegistryPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(reg)
	if err != nil {
		return fmt.Errorf("failed to marshal registry: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}

func RegisterWorkspace(name string, projectPath string, template string, lastOpened string) error {
	reg, err := LoadRegistry()
	if err != nil {
		reg = &WorkspaceRegistry{}
	}

	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		absPath = projectPath
	}

	updated := false
	for i, entry := range reg.Workspaces {
		if entry.Path == absPath {
			reg.Workspaces[i].Name = name
			reg.Workspaces[i].Template = template
			if lastOpened != "" {
				reg.Workspaces[i].LastOpened = lastOpened
			}
			updated = true
			break
		}
	}

	if !updated {
		reg.Workspaces = append(reg.Workspaces, WorkspaceEntry{
			Name:       name,
			Path:       absPath,
			Template:   template,
			LastOpened: lastOpened,
		})
	}

	return SaveRegistry(reg)
}

func UnregisterWorkspace(projectPath string) error {
	reg, err := LoadRegistry()
	if err != nil {
		return nil
	}

	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		absPath = projectPath
	}

	var remaining []WorkspaceEntry
	for _, entry := range reg.Workspaces {
		if entry.Path != absPath && entry.Name != projectPath {
			remaining = append(remaining, entry)
		}
	}

	reg.Workspaces = remaining
	return SaveRegistry(reg)
}

// CheckSessionStatus checks if a Zellij session with the specified name is currently running.
func CheckSessionStatus(runner exec.Runner, name string) string {
	if runner == nil {
		runner = exec.NewOSExecRunner()
	}

	path, err := runner.LookPath("zellij")
	if err != nil {
		return "stopped"
	}

	out, err := runner.CombinedOutput("", path, "list-sessions")
	if err != nil {
		return "stopped"
	}

	// Format of zellij list-sessions is e.g.: "goJPeek [CREATED ...]" or "goJPeek"
	lines := string(out)
	for _, line := range filepath.SplitList(lines) {
		if len(line) > 0 && line == name {
			return "running"
		}
	}

	if len(lines) > 0 {
		for _, l := range splitLines(lines) {
			l = stringsTrim(l)
			if l == name || hasWordPrefix(l, name) {
				return "running"
			}
		}
	}

	return "stopped"
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func stringsTrim(s string) string {
	for len(s) > 0 && (s[0] == ' ' || s[0] == '\t' || s[0] == '\r') {
		s = s[1:]
	}
	for len(s) > 0 && (s[len(s)-1] == ' ' || s[len(s)-1] == '\t' || s[len(s)-1] == '\r') {
		s = s[:len(s)-1]
	}
	return s
}

func hasWordPrefix(s string, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	if s[:len(prefix)] == prefix {
		if len(s) == len(prefix) || s[len(prefix)] == ' ' || s[len(prefix)] == '[' || s[len(prefix)] == '\t' {
			return true
		}
	}
	return false
}
