package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/arthurgray2k/goWorkspace/internal/exec"
	"github.com/arthurgray2k/goWorkspace/internal/workspace"
)

const version = "0.1.0"

func printGlobalHelp() {
	helpText := `goWorkspace (gws) - Lightweight Linux Developer Workspace Manager

Usage:
  gws <command> [target] [options]

Commands:
  init      Initialize workspace configuration (.goworkspace.yaml) in current or target project
  open      Open/attach workspace tools (multiplexer, editor, terminal)
  list      List all registered workspaces and their current status
  status    Show status of a workspace
  resume    Restore/open workspace session
  remove    Remove workspace configuration (NEVER deletes project files)
  config    View or update global configuration (~/.config/goworkspace/config.yaml)

Global Options:
  -h, --help     Show help information
  -v, --version  Show version information

Examples:
  gws init
  gws open
  gws open ~/projects/my-repo
  gws open my-repo
  gws open --no-editor
  gws open --no-multiplexer
  gws list
  gws status
  gws remove

Run 'gws <command> --help' for details on a specific command.
`
	fmt.Print(helpText)
}

func main() {
	if len(os.Args) < 2 {
		printGlobalHelp()
		os.Exit(0)
	}

	subcommand := os.Args[1]
	runner := exec.NewOSExecRunner()

	switch subcommand {
	case "-h", "--help", "help":
		printGlobalHelp()

	case "-v", "--version", "version":
		fmt.Printf("gws version %s\n", version)

	case "init":
		fs := flag.NewFlagSet("init", flag.ExitOnError)
		tmpl := fs.String("template", "", "Project template (go, dotnet, node, rust, generic)")
		ed := fs.String("editor", "", "Editor type (vscode, zed, none)")
		term := fs.String("terminal", "", "Terminal emulator (ghostty, default, none)")
		mux := fs.String("multiplexer", "", "Multiplexer type (zellij, none)")
		name := fs.String("name", "", "Workspace name")
		yes := fs.Bool("y", false, "Accept default configuration without interactive prompts")
		fs.BoolVar(yes, "yes", false, "Accept default configuration without interactive prompts")

		fs.Usage = func() {
			fmt.Print(`gws init - Initialize Workspace

Usage:
  gws init [target_dir] [options]

Options:
  --template <type>     Set project template (go, dotnet, node, rust, generic)
  --editor <type>       Set editor (vscode, zed, none)
  --terminal <type>     Set terminal emulator (ghostty, default, none)
  --multiplexer <type>  Set multiplexer (zellij, none)
  --name <name>         Set project workspace name
  -y, --yes             Non-interactive mode accepting defaults
`)
		}

		_ = fs.Parse(os.Args[2:])
		target := ""
		if fs.NArg() > 0 {
			target = fs.Arg(0)
		}

		opts := workspace.InitOptions{
			Dir:         target,
			Template:    *tmpl,
			Editor:      *ed,
			Terminal:    *term,
			Multiplexer: *mux,
			Name:        *name,
			Yes:         *yes,
		}

		if !*yes {
			opts.Reader = bufio.NewReader(os.Stdin)
		}

		if err := workspace.Init(runner, opts); err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing workspace: %v\n", err)
			os.Exit(1)
		}

	case "open":
		fs := flag.NewFlagSet("open", flag.ExitOnError)
		ed := fs.String("editor", "", "Override editor (vscode, zed, none)")
		noEd := fs.Bool("no-editor", false, "Disable launching editor")
		mux := fs.String("multiplexer", "", "Override multiplexer (zellij, none)")
		noMux := fs.Bool("no-multiplexer", false, "Disable multiplexer")

		fs.Usage = func() {
			fmt.Print(`gws open - Open Workspace

Usage:
  gws open [path_or_workspace_name] [options]

Options:
  --editor <type>       Override editor (vscode, zed, none)
  --no-editor           Do not open configured editor
  --multiplexer <type>   Override multiplexer (zellij, none)
  --no-multiplexer     Do not launch multiplexer session

Examples:
  gws open
  gws open ~/projects/myApp
  gws open myApp
`)
		}

		_ = fs.Parse(os.Args[2:])
		target := ""
		if fs.NArg() > 0 {
			target = fs.Arg(0)
		}

		opts := workspace.OpenOptions{
			Dir:         target,
			Editor:      *ed,
			NoEditor:    *noEd,
			Multiplexer: *mux,
			NoMux:       *noMux,
		}

		if err := workspace.Open(runner, opts); err != nil {
			fmt.Fprintf(os.Stderr, "Error opening workspace: %v\n", err)
			os.Exit(1)
		}

	case "list":
		fs := flag.NewFlagSet("list", flag.ExitOnError)
		fs.Usage = func() {
			fmt.Print(`gws list - List Workspaces

Usage:
  gws list
`)
		}
		_ = fs.Parse(os.Args[2:])

		if err := workspace.List(runner); err != nil {
			fmt.Fprintf(os.Stderr, "Error listing workspaces: %v\n", err)
			os.Exit(1)
		}

	case "status":
		fs := flag.NewFlagSet("status", flag.ExitOnError)
		fs.Usage = func() {
			fmt.Print(`gws status - Workspace Status

Usage:
  gws status [path_or_workspace_name]
`)
		}
		_ = fs.Parse(os.Args[2:])
		target := ""
		if fs.NArg() > 0 {
			target = fs.Arg(0)
		}

		if err := workspace.Status(runner, workspace.StatusOptions{Dir: target}); err != nil {
			fmt.Fprintf(os.Stderr, "Error getting workspace status: %v\n", err)
			os.Exit(1)
		}

	case "resume":
		fs := flag.NewFlagSet("resume", flag.ExitOnError)
		ed := fs.String("editor", "", "Override editor")
		noEd := fs.Bool("no-editor", false, "Disable editor")
		mux := fs.String("multiplexer", "", "Override multiplexer")
		noMux := fs.Bool("no-multiplexer", false, "Disable multiplexer")

		fs.Usage = func() {
			fmt.Print(`gws resume - Restore Workspace

Usage:
  gws resume [path_or_workspace_name] [options]
`)
		}
		_ = fs.Parse(os.Args[2:])
		target := ""
		if fs.NArg() > 0 {
			target = fs.Arg(0)
		}

		opts := workspace.OpenOptions{
			Dir:         target,
			Editor:      *ed,
			NoEditor:    *noEd,
			Multiplexer: *mux,
			NoMux:       *noMux,
		}

		if err := workspace.Resume(runner, opts); err != nil {
			fmt.Fprintf(os.Stderr, "Error resuming workspace: %v\n", err)
			os.Exit(1)
		}

	case "remove":
		fs := flag.NewFlagSet("remove", flag.ExitOnError)
		yes := fs.Bool("y", false, "Confirm removal without prompt")
		fs.BoolVar(yes, "yes", false, "Confirm removal without prompt")

		fs.Usage = func() {
			fmt.Print(`gws remove - Remove Workspace Configuration

Usage:
  gws remove [path_or_workspace_name] [options]

Options:
  -y, --yes  Confirm configuration removal without interactive confirmation prompt
`)
		}
		_ = fs.Parse(os.Args[2:])
		target := ""
		if fs.NArg() > 0 {
			target = fs.Arg(0)
		}

		opts := workspace.RemoveOptions{
			Dir: target,
			Yes: *yes,
		}

		if !*yes {
			opts.Reader = bufio.NewReader(os.Stdin)
		}

		if err := workspace.Remove(runner, opts); err != nil {
			fmt.Fprintf(os.Stderr, "Error removing workspace: %v\n", err)
			os.Exit(1)
		}

	case "config":
		fs := flag.NewFlagSet("config", flag.ExitOnError)
		setEditor := fs.String("set-editor", "", "Set global default editor")
		setTerm := fs.String("set-terminal", "", "Set global default terminal")
		setMux := fs.String("set-multiplexer", "", "Set global default multiplexer")
		setTmpl := fs.String("set-template", "", "Set global default template")

		fs.Usage = func() {
			fmt.Print(`gws config - Global Configuration

Usage:
  gws config [options]

Options:
  --set-editor <editor>          Set default global editor (vscode, zed, none)
  --set-terminal <terminal>      Set default global terminal (ghostty, default, none)
  --set-multiplexer <mux>        Set default global multiplexer (zellij, none)
  --set-template <template>      Set default global template (go, dotnet, node, rust, generic)
`)
		}
		_ = fs.Parse(os.Args[2:])

		opts := workspace.ConfigOptions{
			SetEditor:   *setEditor,
			SetTerm:     *setTerm,
			SetMux:      *setMux,
			SetTemplate: *setTmpl,
		}

		if err := workspace.ConfigCmd(opts); err != nil {
			fmt.Fprintf(os.Stderr, "Error managing global config: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\nRun 'gws --help' for usage.\n", subcommand)
		os.Exit(1)
	}
}
