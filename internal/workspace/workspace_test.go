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

	// Save .goworkspace.yaml in parentDir
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

func TestInitAndOpenWorkspace(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	projDir := filepath.Join(tmpHome, "goJPeek")
	err := os.MkdirAll(projDir, 0755)
	if err != nil {
		t.Fatalf("failed to create projDir: %v", err)
	}
	os.WriteFile(filepath.Join(projDir, "go.mod"), []byte("module goJPeek"), 0644)

	runner := exec.NewMockRunner()
	runner.Available["code"] = true
	runner.Available["zellij"] = true

	// 1. gws init --yes
	err = Init(runner, InitOptions{
		Dir: projDir,
		Yes: true,
	})
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	cfg, err := config.LoadWorkspaceConfig(projDir)
	if err != nil || cfg.Name != "goJPeek" || cfg.Template != "go" {
		t.Fatalf("unexpected init config: err=%v, cfg=%#v", err, cfg)
	}

	// 2. gws open
	t.Setenv("ZELLIJ", "")
	t.Setenv("ZELLIJ_SESSION_NAME", "")
	err = Open(runner, OpenOptions{
		Dir: projDir,
	})
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	// Verify editor was started and zellij was run
	if len(runner.StartCalls) != 1 || runner.StartCalls[0].Name != "/usr/bin/code" {
		t.Fatalf("expected code start call, got %#v", runner.StartCalls)
	}
	if len(runner.RunCalls) != 1 || runner.RunCalls[0].Name != "/usr/bin/zellij" {
		t.Fatalf("expected zellij run call, got %#v", runner.RunCalls)
	}

	// 3. gws remove --yes
	err = Remove(runner, RemoveOptions{
		Dir: projDir,
		Yes: true,
	})
	if err != nil {
		t.Fatalf("Remove failed: %v", err)
	}

	// Ensure project directory still exists!
	if _, err := os.Stat(projDir); os.IsNotExist(err) {
		t.Fatalf("Remove DELETED project directory! Safety violation!")
	}

	// Ensure .goworkspace.yaml was removed
	if _, err := os.Stat(filepath.Join(projDir, config.WorkspaceConfigFilename)); !os.IsNotExist(err) {
		t.Fatalf(".goworkspace.yaml was not removed")
	}
}

func TestStatus(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	projDir := filepath.Join(tmpHome, "goJPeek")
	os.MkdirAll(projDir, 0755)
	cfg := &config.WorkspaceConfig{Name: "goJPeek", Template: "go"}
	config.SaveWorkspaceConfig(projDir, cfg)

	runner := exec.NewMockRunner()
	err := Status(runner, StatusOptions{Dir: projDir})
	if err != nil {
		t.Fatalf("Status failed: %v", err)
	}
}
