package workspace

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/arthurgray2k/goWorkspace/internal/config"
	"github.com/arthurgray2k/goWorkspace/internal/detect"
	"github.com/arthurgray2k/goWorkspace/internal/editor"
	"github.com/arthurgray2k/goWorkspace/internal/exec"
	"github.com/arthurgray2k/goWorkspace/internal/multiplexer"
	"github.com/arthurgray2k/goWorkspace/internal/registry"
	"github.com/arthurgray2k/goWorkspace/internal/template"
	"github.com/arthurgray2k/goWorkspace/internal/terminal"
)

// FindWorkspaceDir traverses up from startDir searching for .goworkspace.yaml.
func FindWorkspaceDir(startDir string) (string, error) {
	abs, err := filepath.Abs(startDir)
	if err != nil {
		abs = startDir
	}

	curr := abs
	for {
		configPath := filepath.Join(curr, config.WorkspaceConfigFilename)
		info, err := os.Stat(configPath)
		if err == nil && !info.IsDir() {
			return curr, nil
		}

		parent := filepath.Dir(curr)
		if parent == curr { // Reached filesystem root
			break
		}
		curr = parent
	}

	return "", config.ErrNoWorkspaceConfig
}

type InitOptions struct {
	Dir         string
	Name        string
	Template    string
	Editor      string
	Terminal    string
	Multiplexer string
	Yes         bool
	Reader      *bufio.Reader
}

func Init(runner exec.Runner, opts InitOptions) error {
	dir := opts.Dir
	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current working directory: %w", err)
		}
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		absDir = dir
	}

	dirName := filepath.Base(absDir)
	existingPath := filepath.Join(absDir, config.WorkspaceConfigFilename)

	if _, err := os.Stat(existingPath); err == nil {
		if !opts.Yes {
			if opts.Reader == nil {
				return fmt.Errorf(".goworkspace.yaml already exists in %s. Pass -y to overwrite", absDir)
			}
			fmt.Printf(".goworkspace.yaml already exists in %s. Overwrite? [y/N]: ", absDir)
			ans, _ := opts.Reader.ReadString('\n')
			ans = strings.ToLower(strings.TrimSpace(ans))
			if ans != "y" && ans != "yes" {
				fmt.Println("Initialization cancelled.")
				return nil
			}
		}
	}

	detectedTemplate := detect.DetectProjectType(absDir)
	globalCfg, _ := config.LoadGlobalConfig()

	chosenTemplate := opts.Template
	if chosenTemplate == "" {
		chosenTemplate = detectedTemplate
	}

	chosenEditor := opts.Editor
	if chosenEditor == "" {
		chosenEditor = globalCfg.Editor
	}

	chosenTerminal := opts.Terminal
	if chosenTerminal == "" {
		chosenTerminal = globalCfg.Terminal
	}

	chosenMux := opts.Multiplexer
	if chosenMux == "" {
		chosenMux = globalCfg.Multiplexer
	}

	projectName := opts.Name
	if projectName == "" {
		projectName = dirName
	}

	// Interactive prompting if reader is provided and non-interactive yes flag not set
	if !opts.Yes && opts.Reader != nil {
		fmt.Println("goWorkspace - Initialize Workspace")
		fmt.Printf("\nProject: %s\n\n", projectName)

		fmt.Printf("Template [%s]: ", chosenTemplate)
		tInput, _ := opts.Reader.ReadString('\n')
		tInput = strings.TrimSpace(tInput)
		if tInput != "" {
			chosenTemplate = tInput
		}

		fmt.Printf("Editor [%s]: ", chosenEditor)
		eInput, _ := opts.Reader.ReadString('\n')
		eInput = strings.TrimSpace(eInput)
		if eInput != "" {
			chosenEditor = eInput
		}

		fmt.Printf("Terminal [%s]: ", chosenTerminal)
		termInput, _ := opts.Reader.ReadString('\n')
		termInput = strings.TrimSpace(termInput)
		if termInput != "" {
			chosenTerminal = termInput
		}

		fmt.Printf("Multiplexer [%s]: ", chosenMux)
		mInput, _ := opts.Reader.ReadString('\n')
		mInput = strings.TrimSpace(mInput)
		if mInput != "" {
			chosenMux = mInput
		}

		fmt.Printf("\nCreate workspace? [Y/n]: ")
		confirm, _ := opts.Reader.ReadString('\n')
		confirm = strings.ToLower(strings.TrimSpace(confirm))
		if confirm == "n" || confirm == "no" {
			fmt.Println("Initialization cancelled.")
			return nil
		}
	}

	wsConfig, err := template.ApplyToWorkspaceConfig(projectName, chosenTemplate, dirName)
	if err != nil {
		return err
	}

	wsConfig.Editor.Type = chosenEditor
	wsConfig.Terminal.Type = chosenTerminal
	wsConfig.Multiplexer.Type = chosenMux

	err = config.SaveWorkspaceConfig(absDir, wsConfig)
	if err != nil {
		return fmt.Errorf("failed to save workspace configuration: %w", err)
	}

	err = registry.RegisterWorkspace(wsConfig.Name, absDir, wsConfig.Template, time.Now().Format(time.RFC3339))
	if err != nil {
		// Non-fatal warning if registry cannot be updated
		fmt.Printf("Warning: Could not register workspace in state file: %v\n", err)
	}

	fmt.Printf("Initialized workspace configuration in %s\n", filepath.Join(absDir, config.WorkspaceConfigFilename))
	return nil
}

type OpenOptions struct {
	Dir         string
	Flags       config.ConfigFlags
	NoEditor    bool
	NoMux       bool
	Editor      string
	Multiplexer string
}

func Open(runner exec.Runner, opts OpenOptions) error {
	dir := opts.Dir
	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
	}

	wsDir, err := FindWorkspaceDir(dir)
	if err != nil {
		return fmt.Errorf("No %s found.\nRun:\n    gws init", config.WorkspaceConfigFilename)
	}

	wsConfig, err := config.LoadWorkspaceConfig(wsDir)
	if err != nil {
		return err
	}

	globalCfg, _ := config.LoadGlobalConfig()

	flags := opts.Flags
	if opts.NoEditor {
		flags.NoEditor = true
	} else if opts.Editor != "" {
		flags.Editor = opts.Editor
	}

	if opts.NoMux {
		flags.NoMultiplexer = true
	} else if opts.Multiplexer != "" {
		flags.Multiplexer = opts.Multiplexer
	}

	resolved := config.Resolve(wsConfig, globalCfg, flags, wsDir)

	// Register last opened
	_ = registry.RegisterWorkspace(resolved.Name, wsDir, resolved.Template, time.Now().Format(time.RFC3339))

	// 1. Launch multiplexer / terminal session if multiplexer enabled
	if resolved.Multiplexer != "none" {
		err := multiplexer.Launch(runner, resolved.Multiplexer, resolved.Name, wsDir, resolved.Panes)
		if err != nil {
			fmt.Printf("Multiplexer error: %v\n", err)
		}
	} else if resolved.Terminal != "none" && resolved.Terminal != "default" && resolved.Terminal != "system" {
		// Launch requested terminal emulator if multiplexer is disabled
		err := terminal.Launch(runner, resolved.Terminal, "", nil, wsDir)
		if err != nil {
			fmt.Printf("Terminal error: %v\n", err)
		}
	}

	// 2. Launch editor if configured and enabled
	if resolved.Editor != "none" {
		err := editor.Launch(runner, resolved.Editor, wsDir)
		if err != nil {
			fmt.Printf("Editor error: %v\n", err)
		}
	}

	return nil
}

func List(runner exec.Runner) error {
	reg, err := registry.LoadRegistry()
	if err != nil {
		return fmt.Errorf("failed to load workspace registry: %w", err)
	}

	if len(reg.Workspaces) == 0 {
		fmt.Println("No workspaces registered. Run 'gws init' inside a project directory.")
		return nil
	}

	fmt.Printf("%-20s %-45s %-10s\n", "NAME", "PATH", "STATUS")
	fmt.Println(strings.Repeat("-", 77))

	for _, entry := range reg.Workspaces {
		status := registry.CheckSessionStatus(runner, entry.Name)
		displayPath := entry.Path
		if home, err := os.UserHomeDir(); err == nil && strings.HasPrefix(entry.Path, home) {
			displayPath = "~" + entry.Path[len(home):]
		}
		fmt.Printf("%-20s %-45s %-10s\n", entry.Name, displayPath, status)
	}

	return nil
}

type StatusOptions struct {
	Dir string
}

func Status(runner exec.Runner, opts StatusOptions) error {
	dir := opts.Dir
	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
	}

	wsDir, err := FindWorkspaceDir(dir)
	if err != nil {
		return fmt.Errorf("No %s found in %s or parent directories.", config.WorkspaceConfigFilename, dir)
	}

	wsConfig, err := config.LoadWorkspaceConfig(wsDir)
	if err != nil {
		return err
	}

	globalCfg, _ := config.LoadGlobalConfig()
	resolved := config.Resolve(wsConfig, globalCfg, config.ConfigFlags{}, wsDir)
	sessionStatus := registry.CheckSessionStatus(runner, resolved.Name)

	fmt.Printf("Workspace:   %s\n", resolved.Name)
	fmt.Printf("Path:        %s\n", wsDir)
	fmt.Printf("Template:    %s\n", resolved.Template)
	fmt.Printf("Editor:      %s\n", resolved.Editor)
	fmt.Printf("Terminal:    %s\n", resolved.Terminal)
	fmt.Printf("Multiplexer: %s\n", resolved.Multiplexer)
	fmt.Printf("\nSession:     %s\n\n", sessionStatus)

	fmt.Println("PANES")
	if len(resolved.Panes) == 0 {
		fmt.Println("  (no panes defined)")
	} else {
		for _, pane := range resolved.Panes {
			cmdStr := pane.Command
			if cmdStr == "" {
				cmdStr = "(interactive shell)"
			}
			fmt.Printf("  %-12s %s\n", pane.Name, cmdStr)
		}
	}

	return nil
}

func Resume(runner exec.Runner, opts OpenOptions) error {
	return Open(runner, opts)
}

type RemoveOptions struct {
	Dir    string
	Yes    bool
	Reader *bufio.Reader
}

func Remove(runner exec.Runner, opts RemoveOptions) error {
	dir := opts.Dir
	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
	}

	wsDir, err := FindWorkspaceDir(dir)
	if err != nil {
		return fmt.Errorf("No %s found to remove in %s.", config.WorkspaceConfigFilename, dir)
	}

	configPath := filepath.Join(wsDir, config.WorkspaceConfigFilename)

	if !opts.Yes {
		if opts.Reader == nil {
			return fmt.Errorf("Confirmation required to remove %s. Pass -y to confirm", configPath)
		}
		fmt.Printf("Are you sure you want to remove workspace configuration at %s? (Project files will NOT be deleted) [y/N]: ", configPath)
		ans, _ := opts.Reader.ReadString('\n')
		ans = strings.ToLower(strings.TrimSpace(ans))
		if ans != "y" && ans != "yes" {
			fmt.Println("Removal cancelled.")
			return nil
		}
	}

	err = os.Remove(configPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete configuration file: %w", err)
	}

	_ = registry.UnregisterWorkspace(wsDir)

	fmt.Printf("Removed workspace configuration %s.\n", configPath)
	return nil
}

type ConfigOptions struct {
	Show        bool
	SetEditor   string
	SetTerm     string
	SetMux      string
	SetTemplate string
}

func ConfigCmd(opts ConfigOptions) error {
	cfg, err := config.LoadGlobalConfig()
	if err != nil {
		cfg = config.DefaultGlobalConfig()
	}

	modified := false
	if opts.SetEditor != "" {
		cfg.Editor = opts.SetEditor
		modified = true
	}
	if opts.SetTerm != "" {
		cfg.Terminal = opts.SetTerm
		modified = true
	}
	if opts.SetMux != "" {
		cfg.Multiplexer = opts.SetMux
		modified = true
	}
	if opts.SetTemplate != "" {
		cfg.DefaultTemplate = opts.SetTemplate
		modified = true
	}

	if modified {
		err := config.SaveGlobalConfig(cfg)
		if err != nil {
			return fmt.Errorf("failed to save global configuration: %w", err)
		}
		fmt.Println("Global configuration updated successfully.")
	}

	path, _ := config.GlobalConfigPath()
	fmt.Printf("Global Configuration (%s):\n", path)
	fmt.Printf("  editor:           %s\n", cfg.Editor)
	fmt.Printf("  terminal:         %s\n", cfg.Terminal)
	fmt.Printf("  multiplexer:      %s\n", cfg.Multiplexer)
	fmt.Printf("  default_template: %s\n", cfg.DefaultTemplate)

	return nil
}
