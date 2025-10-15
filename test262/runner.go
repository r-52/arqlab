package test262

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Runner coordinates discovery and execution of Test262 compliance tests.
type Runner struct {
	// RootDir holds the path to the cloned test262 repository.
	RootDir string
	// OutDir is where harness artifacts such as filtered lists and reports are written.
	OutDir string
	// SkipAsync controls whether async/await tests are excluded.
	SkipAsync bool
}

// TestCase describes a single Test262 test file.
type TestCase struct {
	Path        string
	Description string
	Flags       []string
}

// Report aggregates the outcome of a single test run.
type Report struct {
	Total   int
	Passed  int
	Failed  int
	Skipped int
}

// NewRunner validates the file system layout and returns a configured Runner.
func NewRunner(rootDir, outDir string) (*Runner, error) {
	if rootDir == "" {
		return nil, errors.New("test262 root directory cannot be empty")
	}
	info, err := os.Stat(rootDir)
	if err != nil {
		return nil, fmt.Errorf("stat test262 root: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("test262 root %q is not a directory", rootDir)
	}

	if outDir == "" {
		outDir = filepath.Join(rootDir, "..", "out")
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return nil, fmt.Errorf("create output directory: %w", err)
	}

	return &Runner{RootDir: rootDir, OutDir: outDir, SkipAsync: true}, nil
}

// Discover walks the test262 repository and returns metadata for each test file.
func (r *Runner) Discover() ([]TestCase, error) {
	return nil, errors.New("test discovery not implemented yet")
}

// Run executes the provided test cases and returns a summarized report.
func (r *Runner) Run(cases []TestCase) (*Report, error) {
	if len(cases) == 0 {
		return &Report{}, nil
	}
	return nil, errors.New("test execution not implemented yet")
}
