package registry

import (
	"path/filepath"
	"testing"

	"github.com/arthurgray2k/goWorkspace/internal/exec"
)

func TestRegistryOperations(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	err := RegisterWorkspace("goJPeek", filepath.Join(tmpHome, "goJPeek"), "go", "2026-08-08T23:30:00Z")
	if err != nil {
		t.Fatalf("failed to register workspace: %v", err)
	}

	reg, err := LoadRegistry()
	if err != nil || len(reg.Workspaces) != 1 {
		t.Fatalf("expected 1 workspace entry, got err=%v, count=%d", err, len(reg.Workspaces))
	}

	if reg.Workspaces[0].Name != "goJPeek" || reg.Workspaces[0].Template != "go" {
		t.Fatalf("unexpected entry: %#v", reg.Workspaces[0])
	}

	err = UnregisterWorkspace(filepath.Join(tmpHome, "goJPeek"))
	if err != nil {
		t.Fatalf("failed to unregister: %v", err)
	}

	reg2, err := LoadRegistry()
	if err != nil || len(reg2.Workspaces) != 0 {
		t.Fatalf("expected 0 workspaces after unregister, got %d", len(reg2.Workspaces))
	}
}

func TestCheckSessionStatusWithANSI(t *testing.T) {
	mock := exec.NewMockRunner()
	mock.Available["zellij"] = true

	// Simulate actual zellij list-sessions output with ANSI escape codes
	mock.OutputReturn["zellij"] = []byte("\x1b[32;1mbrave-foxglove\x1b[m [Created \x1b[35;1m10h\x1b[m ago] (\x1b[31;1mEXITED\x1b[m - attach)\n\x1b[32;1mgoJPeek\x1b[m [Created \x1b[35;1m23m\x1b[m ago]\n")

	statusRunning := CheckSessionStatus(mock, "goJPeek")
	if statusRunning != "running" {
		t.Fatalf("expected running status for goJPeek, got %s", statusRunning)
	}

	statusExited := CheckSessionStatus(mock, "brave-foxglove")
	if statusExited != "stopped" {
		t.Fatalf("expected stopped status for brave-foxglove, got %s", statusExited)
	}

	statusStopped := CheckSessionStatus(mock, "nonExistent")
	if statusStopped != "stopped" {
		t.Fatalf("expected stopped status for nonExistent, got %s", statusStopped)
	}
}
