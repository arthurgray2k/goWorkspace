package exec

import (
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

// Runner abstracts external process execution for testing and execution.
type Runner interface {
	LookPath(file string) (string, error)
	Run(dir string, name string, args ...string) error
	Start(dir string, name string, args ...string) error
	CombinedOutput(dir string, name string, args ...string) ([]byte, error)
}

// OSExecRunner is the default implementation of Runner using os/exec.
type OSExecRunner struct{}

func NewOSExecRunner() *OSExecRunner {
	return &OSExecRunner{}
}

func (r *OSExecRunner) LookPath(file string) (string, error) {
	return exec.LookPath(file)
}

func (r *OSExecRunner) Run(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (r *OSExecRunner) Start(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	return cmd.Start()
}

func (r *OSExecRunner) CombinedOutput(dir string, name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	return cmd.CombinedOutput()
}

// CallRecord records a single command execution invocation for testing.
type CallRecord struct {
	Dir  string
	Name string
	Args []string
}

// MockRunner is a mock implementation of Runner for unit testing.
type MockRunner struct {
	mu           sync.Mutex
	Available    map[string]bool
	RunCalls     []CallRecord
	StartCalls   []CallRecord
	OutputCalls  []CallRecord
	RunErr       map[string]error
	OutputReturn map[string][]byte
}

func NewMockRunner() *MockRunner {
	return &MockRunner{
		Available:    make(map[string]bool),
		RunErr:       make(map[string]error),
		OutputReturn: make(map[string][]byte),
	}
}

func (m *MockRunner) LookPath(file string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.Available[file] {
		return "/usr/bin/" + file, nil
	}
	return "", exec.ErrNotFound
}

func (m *MockRunner) Run(dir string, name string, args ...string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.RunCalls = append(m.RunCalls, CallRecord{Dir: dir, Name: name, Args: args})
	if err, ok := m.RunErr[name]; ok {
		return err
	}
	base := filepath.Base(name)
	if err, ok := m.RunErr[base]; ok {
		return err
	}
	return nil
}

func (m *MockRunner) Start(dir string, name string, args ...string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.StartCalls = append(m.StartCalls, CallRecord{Dir: dir, Name: name, Args: args})
	if err, ok := m.RunErr[name]; ok {
		return err
	}
	base := filepath.Base(name)
	if err, ok := m.RunErr[base]; ok {
		return err
	}
	return nil
}

func (m *MockRunner) CombinedOutput(dir string, name string, args ...string) ([]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.OutputCalls = append(m.OutputCalls, CallRecord{Dir: dir, Name: name, Args: args})
	if err, ok := m.RunErr[name]; ok {
		return m.OutputReturn[name], err
	}
	if ret, ok := m.OutputReturn[name]; ok {
		return ret, nil
	}
	base := filepath.Base(name)
	if err, ok := m.RunErr[base]; ok {
		return m.OutputReturn[base], err
	}
	if ret, ok := m.OutputReturn[base]; ok {
		return ret, nil
	}
	return nil, nil
}
