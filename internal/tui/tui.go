package tui

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"

	"github.com/arthurgray2k/goWorkspace/internal/exec"
	"github.com/arthurgray2k/goWorkspace/internal/registry"
	"github.com/arthurgray2k/goWorkspace/internal/workspace"
)

// Termios struct for POSIX raw terminal handling on Linux
type termios struct {
	Iflag  uint32
	Oflag  uint32
	Cflag  uint32
	Lflag  uint32
	Cc     [32]uint8
	Ispeed uint32
	Ospeed uint32
}

const (
	TCGETS = 0x5401
	TCSETS = 0x5402
)

func setRawMode(fd uintptr) (*termios, error) {
	var oldState termios
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(TCGETS), uintptr(unsafe.Pointer(&oldState)))
	if err != 0 {
		return nil, err
	}

	raw := oldState
	// Disable echo and canonical mode
	raw.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.ISIG | syscall.IEXTEN
	raw.Iflag &^= syscall.IXON | syscall.ICRNL
	raw.Cc[syscall.VMIN] = 1
	raw.Cc[syscall.VTIME] = 0

	_, _, err = syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(TCSETS), uintptr(unsafe.Pointer(&raw)))
	if err != 0 {
		return nil, err
	}

	return &oldState, nil
}

func restoreMode(fd uintptr, oldState *termios) {
	if oldState == nil {
		return
	}
	_, _, _ = syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(TCSETS), uintptr(unsafe.Pointer(oldState)))
}

// RunDashboard launches the interactive terminal user interface.
func RunDashboard(runner exec.Runner) error {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf("TUI requires an interactive terminal (/dev/tty): %w", err)
	}
	defer tty.Close()

	fd := tty.Fd()
	oldState, err := setRawMode(fd)
	if err != nil {
		return fmt.Errorf("failed to enable raw terminal mode: %w", err)
	}
	defer restoreMode(fd, oldState)

	// Hide cursor
	fmt.Fprint(tty, "\x1b[?25l")
	defer fmt.Fprint(tty, "\x1b[?25h\x1b[0m\n")

	selectedIndex := 0
	msg := ""

	for {
		reg, err := registry.LoadRegistry()
		if err != nil {
			reg = &registry.WorkspaceRegistry{}
		}

		workspaces := reg.Workspaces
		if selectedIndex >= len(workspaces) {
			selectedIndex = len(workspaces) - 1
		}
		if selectedIndex < 0 {
			selectedIndex = 0
		}

		// Clear screen & render
		fmt.Fprint(tty, "\x1b[2J\x1b[H")
		renderDashboard(tty, runner, workspaces, selectedIndex, msg)
		msg = ""

		// Read keypress
		buf := make([]byte, 3)
		n, err := tty.Read(buf)
		if err != nil || n == 0 {
			break
		}

		// Arrow Keys or Keybindings
		if n >= 3 && buf[0] == 0x1b && buf[1] == '[' {
			switch buf[2] {
			case 'A': // Up
				if selectedIndex > 0 {
					selectedIndex--
				}
			case 'B': // Down
				if selectedIndex < len(workspaces)-1 {
					selectedIndex++
				}
			}
			continue
		}

		key := buf[0]
		switch key {
		case 'q', 'Q', 0x1b, 0x03: // Quit, Esc, Ctrl+C
			return nil

		case 'k', 'K': // Up
			if selectedIndex > 0 {
				selectedIndex--
			}

		case 'j', 'J': // Down
			if selectedIndex < len(workspaces)-1 {
				selectedIndex++
			}

		case 'o', 'O': // Open workspace
			if len(workspaces) > 0 {
				target := workspaces[selectedIndex]
				restoreMode(fd, oldState)
				fmt.Fprint(tty, "\x1b[?25h\x1b[2J\x1b[H")
				fmt.Printf("Opening workspace '%s' (%s)...\n", target.Name, target.Path)
				_ = workspace.Open(runner, workspace.OpenOptions{Dir: target.Path})
				return nil
			}

		case 'x', 'X': // Stop active session
			if len(workspaces) > 0 {
				target := workspaces[selectedIndex]
				_ = workspace.Stop(runner, workspace.StopOptions{Dir: target.Path})
				msg = fmt.Sprintf("Stopped session '%s'", target.Name)
			}

		case 's', 'S': // Status
			if len(workspaces) > 0 {
				target := workspaces[selectedIndex]
				restoreMode(fd, oldState)
				fmt.Fprint(tty, "\x1b[?25h\x1b[2J\x1b[H")
				_ = workspace.Status(runner, workspace.StatusOptions{Dir: target.Path})
				fmt.Println("\nPress Enter to return to dashboard...")
				buf := make([]byte, 1)
				_, _ = tty.Read(buf)
				oldState, _ = setRawMode(fd)
				fmt.Fprint(tty, "\x1b[?25l")
			}
		}
	}

	return nil
}

func renderDashboard(tty *os.File, runner exec.Runner, workspaces []registry.WorkspaceEntry, selectedIndex int, msg string) {
	var sb strings.Builder

	sb.WriteString("\x1b[1;36mgoWorkspace (gws)\x1b[0m \x1b[90m-\x1b[0m \x1b[1mInteractive Dashboard\x1b[0m\n")
	sb.WriteString(strings.Repeat("─", 78) + "\n\n")

	if len(workspaces) == 0 {
		sb.WriteString("  \x1b[33mNo workspaces registered.\x1b[0m\n")
		sb.WriteString("  Run '\x1b[1mgws init\x1b[0m' inside a project directory to register a workspace.\n\n")
	} else {
		sb.WriteString(fmt.Sprintf("  \x1b[1m%-3s %-20s %-42s %-10s\x1b[0m\n", "", "NAME", "PATH", "STATUS"))
		sb.WriteString("  " + strings.Repeat("─", 75) + "\n")

		home, _ := os.UserHomeDir()
		for i, entry := range workspaces {
			status := registry.CheckSessionStatus(runner, entry.Name)
			statusFormatted := "\x1b[32m[RUNNING]\x1b[0m"
			if status != "running" {
				statusFormatted = "\x1b[90m[STOPPED]\x1b[0m"
			}

			displayPath := entry.Path
			if home != "" && strings.HasPrefix(entry.Path, home) {
				displayPath = "~" + entry.Path[len(home):]
			}
			if len(displayPath) > 40 {
				displayPath = displayPath[:37] + "..."
			}

			if i == selectedIndex {
				sb.WriteString(fmt.Sprintf("\x1b[7m> %-3s %-20s %-42s %-10s\x1b[0m\n", "", entry.Name, displayPath, status))
			} else {
				sb.WriteString(fmt.Sprintf("  %-3s %-20s %-42s %-10s\n", "", entry.Name, displayPath, statusFormatted))
			}
		}
		sb.WriteString("\n")
	}

	if msg != "" {
		sb.WriteString(fmt.Sprintf("  \x1b[33m%s\x1b[0m\n\n", msg))
	} else {
		sb.WriteString("\n")
	}

	sb.WriteString(strings.Repeat("─", 78) + "\n")
	sb.WriteString("\x1b[1mControls:\x1b[0m \x1b[36m[↑/↓/j/k]\x1b[0m Navigate  \x1b[36m[o]\x1b[0m Open  \x1b[36m[s]\x1b[0m Status  \x1b[36m[x]\x1b[0m Stop Session  \x1b[36m[q]\x1b[0m Quit\n")

	fmt.Fprint(tty, sb.String())
}
