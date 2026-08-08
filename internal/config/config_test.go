package config

import (
	"testing"
)

func TestWorkspaceConfigLoadSave(t *testing.T) {
	dir := t.TempDir()

	cfg := &WorkspaceConfig{
		Name:        "testProj",
		Template:    "go",
		Editor:      ToolConfig{Type: "vscode"},
		Terminal:    ToolConfig{Type: "ghostty"},
		Multiplexer: ToolConfig{Type: "zellij"},
		Panes: []PaneConfig{
			{Name: "shell", Command: ""},
			{Name: "git", Command: "git status"},
		},
		Env: map[string]string{"FOO": "BAR"},
	}

	err := SaveWorkspaceConfig(dir, cfg)
	if err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	loaded, err := LoadWorkspaceConfig(dir)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if loaded.Name != "testProj" || loaded.Template != "go" || loaded.Editor.Type != "vscode" {
		t.Fatalf("unexpected loaded config: %#v", loaded)
	}

	if len(loaded.Panes) != 2 || loaded.Env["FOO"] != "BAR" {
		t.Fatalf("unexpected loaded panes/env: %#v", loaded)
	}
}

func TestResolvePrecedence(t *testing.T) {
	dir := "/home/user/project"

	global := &GlobalConfig{
		Editor:          "vscode",
		Terminal:        "ghostty",
		Multiplexer:     "zellij",
		DefaultTemplate: "go",
	}

	ws := &WorkspaceConfig{
		Name:     "customProject",
		Template: "node",
		Editor:   ToolConfig{Type: "zed"},
	}

	// 1. Without flags, WS > Global
	res1 := Resolve(ws, global, ConfigFlags{}, dir)
	if res1.Name != "customProject" || res1.Template != "node" || res1.Editor != "zed" || res1.Terminal != "ghostty" {
		t.Fatalf("unexpected res1: %#v", res1)
	}

	// 2. With CLI flags, CLI > WS > Global
	res2 := Resolve(ws, global, ConfigFlags{NoEditor: true, Multiplexer: "none"}, dir)
	if res2.Editor != "none" || res2.Multiplexer != "none" {
		t.Fatalf("unexpected res2: %#v", res2)
	}
}
