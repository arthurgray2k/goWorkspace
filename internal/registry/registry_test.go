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

func TestCheckSessionStatus(t *testing.T) {
	mock := exec.NewMockRunner()
	mock.Available["zellij"] = true
	mock.OutputReturn["zellij"] = []byte("goJPeek [CREATED 2 HOURS AGO]\nsomeOtherSession")

	status := CheckSessionStatus(mock, "goJPeek")
	if status != "running" {
		t.Fatalf("expected running status, got %s", status)
	}

	statusStopped := CheckSessionStatus(mock, "nonExistent")
	if statusStopped != "stopped" {
		t.Fatalf("expected stopped status, got %s", statusStopped)
	}
}
