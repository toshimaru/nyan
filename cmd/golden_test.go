package cmd

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/toshimaru/nyan/styles"
)

// updateGolden flag allows regenerating golden files.
// Usage: go test ./cmd -run TestGolden -update
var updateGolden = flag.Bool("update", false, "update golden files")

const goldenDir = "testdata/golden"

type goldenTestCase struct {
	name string
	args []string
}

// TestGoldenOutput tests syntax highlighting output against golden files.
func TestGoldenOutput(t *testing.T) {
	var tests []goldenTestCase

	// Add test cases for all themes
	for _, themeName := range styles.Names() {
		tests = append(tests, goldenTestCase{
			name: themeName,
			args: []string{"--theme", themeName, "testdata/dummy.go"},
		})
	}

	// Add line-numbered output test
	tests = append(tests, goldenTestCase{
		name: "monokai-numbered",
		args: []string{"--theme", "monokai", "--number", "testdata/dummy.go"},
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goldenPath := goldenFilePath(tt.name)
			actual := runNyanAndCapture(t, tt.args)
			compareOrUpdateGolden(t, goldenPath, actual)
		})
	}
}

// goldenFilePath returns the path to a golden file for the given name.
func goldenFilePath(name string) string {
	return filepath.Join(goldenDir, name+".golden")
}

// runNyanAndCapture executes nyan with the given arguments and captures output.
func runNyanAndCapture(t *testing.T, args []string) string {
	t.Helper()

	var o, e bytes.Buffer
	// Mock terminal to enable highlighting
	originalIsTerminalFunc := isTerminalFunc
	isTerminalFunc = func(fd uintptr) bool { return true }
	defer func() { isTerminalFunc = originalIsTerminalFunc }()

	rootCmd.SetArgs(args)
	rootCmd.SetOut(&o)
	rootCmd.SetErr(&e)
	rootCmd.SetIn(nil)
	err := rootCmd.Execute()

	require.NoError(t, err)
	require.Empty(t, e.String())

	// Reset flags after execution
	resetStrings()
	resetFlags()

	return o.String()
}

// compareOrUpdateGolden compares output against golden file, or updates it if -update flag is set.
func compareOrUpdateGolden(t *testing.T, goldenPath, actual string) {
	t.Helper()

	if *updateGolden {
		// Ensure directory exists
		dir := filepath.Dir(goldenPath)
		err := os.MkdirAll(dir, 0755)
		require.NoError(t, err, "failed to create golden directory")

		// Write the golden file
		err = os.WriteFile(goldenPath, []byte(actual), 0644)
		require.NoError(t, err, "failed to write golden file")
		t.Logf("Updated golden file: %s", goldenPath)
		return
	}

	// Read expected output from golden file
	expected, err := os.ReadFile(goldenPath)
	require.NoError(t, err, "failed to read golden file %s (run with -update to create)", goldenPath)

	// Compare
	assert.Equal(t, string(expected), actual,
		"Output does not match golden file %s\nRun 'go test ./cmd -run TestGolden -update' to update golden files", goldenPath)
}
