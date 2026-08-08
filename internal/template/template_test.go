package template

import (
	"testing"
)

func TestGetTemplate(t *testing.T) {
	tmpl, ok := GetTemplate("go")
	if !ok || tmpl.Name != "go" {
		t.Fatalf("expected go template, got ok=%v, tmpl=%#v", ok, tmpl)
	}

	if len(tmpl.Panes) < 3 {
		t.Fatalf("expected at least 3 panes in go template, got %d", len(tmpl.Panes))
	}

	_, ok = GetTemplate("nonexistent")
	if ok {
		t.Fatalf("expected false for nonexistent template")
	}
}

func TestApplyToWorkspaceConfig(t *testing.T) {
	cfg, err := ApplyToWorkspaceConfig("myGoApp", "go", "myGoApp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Name != "myGoApp" || cfg.Template != "go" || cfg.Editor.Type != "vscode" {
		t.Fatalf("unexpected config: %#v", cfg)
	}
}
