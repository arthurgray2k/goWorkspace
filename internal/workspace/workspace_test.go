package workspace

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/arthurgray2k/goWorkspace/internal/config"
	"github.com/arthurgray2k/goWorkspace/internal/exec"
)

func TestFindWorkspaceDir(t *testing.T) {
	parentDir := t.TempDir()
	childDir := filepath.Join(parentDir, "nested", "subfolder")
	err := os.MkdirAll(childDir, 0755)
	if err != nil {
		t.Fatalf("failed to create dirs: %v", err)
	}

	cfg := &config.WorkspaceConfig{Name: "parentWorkspace", Template: "go"}
	err = config.SaveWorkspaceConfig(parentDir, cfg)
	if err != nil {
		t.Fatalf("failed to save workspace config: %v", err)
	}

	foundDir, err := FindWorkspaceDir(childDir)
	if err != nil || foundDir != parentDir {
		t.Fatalf("expected foundDir=%s, got err=%v, dir=%s", parentDir, err, foundDir)
	}
}

func TestMultiTabStatus(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	projDir := filepath.Join(tmpHome, "multiTabProject")
	os.MkdirAll(projDir, 0755)

	cfg := &config.WorkspaceConfig{
		Name:     "multiTabProject",
		Template: "go",
		Tabs: []config.TabConfig{
			{
				Name: "Code",
				Panes: []config.PaneConfig{
					{Name: "shell", Command: ""},
					{Name: "git", Command: "git status", Dir: "./subfolder"},
				},
			},
			{
				Name: "Services",
				Panes: []config.PaneConfig{
					{Name: "server", Command: "go run main.go", Dir: "./cmd/server", Env: map[string]string{"PORT": "8080"}},
				},
			},
		},
	}
	config.SaveWorkspaceConfig(projDir, cfg)

	runner := exec.NewMockRunner()
	err := Status(runner, StatusOptions{Dir: projDir})
	if err != nil {
		t.Fatalf("Status failed for multi-tab project: %v", err)
	}
}
