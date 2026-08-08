package exec

import (
	"os/exec"
	"testing"
)

func TestMockRunner(t *testing.T) {
	mock := NewMockRunner()
	mock.Available["code"] = true
	mock.Available["zellij"] = true

	path, err := mock.LookPath("code")
	if err != nil || path != "/usr/bin/code" {
		t.Fatalf("expected /usr/bin/code, got path=%q, err=%v", path, err)
	}

	_, err = mock.LookPath("ghostty")
	if err != exec.ErrNotFound {
		t.Fatalf("expected ErrNotFound for ghostty, got %v", err)
	}

	err = mock.Start("/tmp", "code", ".")
	if err != nil {
		t.Fatalf("expected start to succeed, got %v", err)
	}

	if len(mock.StartCalls) != 1 || mock.StartCalls[0].Name != "code" {
		t.Fatalf("expected 1 start call for code, got %#v", mock.StartCalls)
	}
}
