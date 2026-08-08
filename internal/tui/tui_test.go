package tui

import (
	"os"
	"testing"

	"github.com/arthurgray2k/goWorkspace/internal/exec"
	"github.com/arthurgray2k/goWorkspace/internal/registry"
)

func TestRenderDashboard(t *testing.T) {
	mock := exec.NewMockRunner()
	workspaces := []registry.WorkspaceEntry{
		{Name: "goJPeek", Path: "/home/user/goJPeek", Template: "go"},
		{Name: "goMini", Path: "/home/user/goMini", Template: "go"},
	}

	tmpFile, err := os.CreateTemp("", "tui_test")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	renderDashboard(tmpFile, mock, workspaces, 0, "Test Message")

	info, err := tmpFile.Stat()
	if err != nil || info.Size() == 0 {
		t.Fatalf("expected non-empty TUI output, got size=%d, err=%v", info.Size(), err)
	}
}
