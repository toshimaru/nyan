package cmd

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	highlightedGoCode   = "[38;5;197mpackage[0m[38;5;231m"
	unhighlightedGoCode = "[38;5;231mpackage main[0m[38;5;231m"
)

func TestExecute(t *testing.T) {
	rootCmd.SetArgs([]string{"--help"})
	Execute()
	resetFlags()
}

func TestCommandExecute(t *testing.T) {
	err := rootCmd.Execute()

	assert.Nil(t, err)
}

func TestHelpCommand(t *testing.T) {
	var o, e bytes.Buffer
	rootCmd.SetArgs([]string{"--help"})
	rootCmd.SetOut(&o)
	rootCmd.SetErr(&e)
	err := rootCmd.Execute()
	resetFlags()

	assert.Nil(t, err)
	assert.Empty(t, e.String())
	assert.Contains(t, o.String(), rootCmd.Use)
	assert.Contains(t, o.String(), rootCmd.Long)
	assert.Contains(t, o.String(), rootCmd.Example)
}

func TestInvalidFilename(t *testing.T) {
	var o, e bytes.Buffer
	rootCmd.SetArgs([]string{"InvalidFilename"})
	rootCmd.SetOut(&o)
	rootCmd.SetErr(&e)
	err := rootCmd.Execute()

	assert.Error(t, err)
	assert.Empty(t, o.String())
	assert.Contains(t, e.String(), invalidFileErrorMsg())
}

func TestCmdExecute(t *testing.T) {
	var o, e bytes.Buffer
	isTerminalFunc = func(fd uintptr) bool { return true }
	rootCmd.SetArgs([]string{"testdata/dummy.go"})
	rootCmd.SetOut(&o)
	rootCmd.SetErr(&e)
	err := rootCmd.Execute()

	assert.Nil(t, err)
	assert.Empty(t, e.String())
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), highlightedGoCode)
}

func TestUnknownExtension(t *testing.T) {
	var o, e bytes.Buffer
	isTerminalFunc = func(fd uintptr) bool { return true }
	rootCmd.SetArgs([]string{"testdata/dummy.go.unknown"})
	rootCmd.SetOut(&o)
	rootCmd.SetErr(&e)
	err := rootCmd.Execute()

	assert.Nil(t, err)
	assert.Empty(t, e.String())
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), _unhighlightedGoCode())
}

func TestLanguageOption(t *testing.T) {
	var o, e bytes.Buffer
	isTerminalFunc = func(fd uintptr) bool { return true }
	rootCmd.SetArgs([]string{"--language", "go", "testdata/dummy.go.unknown"})
	rootCmd.SetOut(&o)
	rootCmd.SetErr(&e)
	err := rootCmd.Execute()
	resetStrings()

	assert.Nil(t, err)
	assert.Empty(t, e.String())
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), highlightedGoCode)
}

func TestInvlaidLanguageOption(t *testing.T) {
	var o, e bytes.Buffer
	isTerminalFunc = func(fd uintptr) bool { return true }
	rootCmd.SetArgs([]string{"--language", "invalid_lang", "testdata/dummy.go"})
	rootCmd.SetOut(&o)
	rootCmd.SetErr(&e)
	err := rootCmd.Execute()
	resetStrings()

	assert.Nil(t, err)
	assert.Empty(t, e.String())
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), _unhighlightedGoCode())
}

func TestMultipleFiles(t *testing.T) {
	var o, e bytes.Buffer
	isTerminalFunc = func(fd uintptr) bool { return true }
	rootCmd.SetArgs([]string{"testdata/dummy.go", "testdata/dummy.go.unknown"})
	rootCmd.SetOut(&o)
	rootCmd.SetErr(&e)
	err := rootCmd.Execute()

	assert.Nil(t, err)
	assert.Empty(t, e.String())
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), highlightedGoCode)
	assert.Contains(t, o.String(), _unhighlightedGoCode())
}

func TestMultipleFilesWithInvalidFileError(t *testing.T) {
	var o, e bytes.Buffer
	isTerminalFunc = func(fd uintptr) bool { return true }
	rootCmd.SetArgs([]string{"testdata/dummy.go", "InvalidFilename", "testdata/dummyfile"})
	rootCmd.SetOut(&o)
	rootCmd.SetErr(&e)
	err := rootCmd.Execute()

	assert.Error(t, err)
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), highlightedGoCode)
	assert.Contains(t, o.String(), "[38;5;231mThis is dummy.[0m")
	assert.Contains(t, e.String(), invalidFileErrorMsg())
}

func testThemes(t *testing.T) {
	var o, e bytes.Buffer
	isTerminalFunc = func(fd uintptr) bool { return true }
	rootCmd.SetOut(&o)
	rootCmd.SetErr(&e)

	t.Run("Valid Theme", func(t *testing.T) {
		rootCmd.SetArgs([]string{"testdata/dummy.go", "--theme", "vim"})
		err := rootCmd.Execute()
		resetStrings()

		assert.Nil(t, err)
		assert.Empty(t, e.String())
		assert.NotNil(t, o.String())
		assert.Contains(t, o.String(), "[38;5;164mpackage[0m")
	})

	t.Run("Inalid Theme", func(t *testing.T) {
		o.Reset()
		rootCmd.SetArgs([]string{"testdata/dummy.go", "--theme", "invalid"})
		err := rootCmd.Execute()
		resetStrings()

		assert.Nil(t, err)
		assert.Empty(t, e.String())
		assert.NotNil(t, o.String())
		assert.Contains(t, o.String(), "[1m[38;5;231mpackage")
	})
}

func TestSpecialFlags(t *testing.T) {
	var o, e bytes.Buffer
	rootCmd.SetOut(&o)
	rootCmd.SetErr(&e)

	t.Run("version Flag", func(t *testing.T) {
		rootCmd.SetArgs([]string{"--version"})
		err := rootCmd.Execute()
		resetFlags()

		assert.Nil(t, err)
		assert.Empty(t, e.String())
		assert.NotNil(t, o.String())
		assert.Contains(t, o.String(), "version ")
		assert.NotContains(t, o.String(), "Theme: ")
	})

	t.Run("listThemes Flag", func(t *testing.T) {
		o.Reset()
		rootCmd.SetArgs([]string{"--list-themes"})
		err := rootCmd.Execute()
		resetFlags()

		assert.Nil(t, err)
		assert.Empty(t, e.String())
		assert.NotNil(t, o.String())
		assert.Contains(t, o.String(), "Theme: ")
		assert.Contains(t, o.String(), "Sample Code in Go")
		assert.NotContains(t, o.String(), "version ")
	})

	t.Run("multiple flags", func(t *testing.T) {
		o.Reset()
		rootCmd.SetArgs([]string{"--version", "--list-themes"})
		err := rootCmd.Execute()
		resetFlags()

		assert.Nil(t, err)
		assert.Empty(t, e.String())
		assert.NotNil(t, o.String())
		assert.Contains(t, o.String(), "version ")
		assert.NotContains(t, o.String(), "Theme: ")
	})
}

func TestUnknownFile(t *testing.T) {
	var o, e bytes.Buffer
	rootCmd.SetArgs([]string{"testdata/dummyfile"})
	rootCmd.SetOut(&o)
	rootCmd.SetErr(&e)
	err := rootCmd.Execute()

	assert.Nil(t, err)
	assert.Empty(t, e.String())
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), "This is dummy.")
}

func TestFromStdIn(t *testing.T) {
	i := bytes.NewBufferString("package main")
	var o, e bytes.Buffer
	isTerminalFunc = func(fd uintptr) bool { return true }
	rootCmd.SetArgs([]string{"--theme", "monokai"})
	rootCmd.SetIn(i)
	rootCmd.SetOut(&o)
	rootCmd.SetErr(&e)
	err := rootCmd.Execute()

	assert.Nil(t, err)
	assert.Empty(t, e.String())
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), highlightedGoCode)
}

func TestFromStdInWithLanguageOption(t *testing.T) {
	i := bytes.NewBufferString("package main")
	var o, e bytes.Buffer
	isTerminalFunc = func(fd uintptr) bool { return true }
	rootCmd.SetArgs([]string{"--theme", "monokai", "--language", "go"})
	rootCmd.SetIn(i)
	rootCmd.SetOut(&o)
	rootCmd.SetErr(&e)
	err := rootCmd.Execute()

	assert.Nil(t, err)
	assert.Empty(t, e.String())
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), highlightedGoCode)
}

func TestFromStdInWithDash(t *testing.T) {
	i := bytes.NewBufferString("TestFromStdIn")
	var o, e bytes.Buffer
	rootCmd.SetArgs([]string{"-"})
	rootCmd.SetIn(i)
	rootCmd.SetOut(&o)
	rootCmd.SetErr(&e)
	err := rootCmd.Execute()

	assert.Nil(t, err)
	assert.Empty(t, e.String())
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), "TestFromStdIn")
}

func TestNumberOption(t *testing.T) {
	var o, e bytes.Buffer
	rootCmd.SetArgs([]string{"-n", "testdata/dummy.go"})
	rootCmd.SetIn(nil)
	rootCmd.SetOut(&o)
	rootCmd.SetErr(&e)
	err := rootCmd.Execute()

	assert.Nil(t, err)

	// Line number check at the beginning of a line.
	lines := strings.Split(o.String(), "\n")
	for i, line := range lines {
		want := fmt.Sprintf("%6d\t", i+1)
		if !strings.HasPrefix(line, want) {
			t.Logf("want: %s got: %s", want, line)
		}
	}

	// EOF line feed check.
	lastLine := lines[len(lines)-1]
	if strings.Contains(lastLine, "\n") {
		t.Fatal("The EOF has an unnecessary line break")
	}
}

func resetFlags() {
	showVersion = false
	listThemes = false
	rootCmd.Flags().Set("help", "false")
}

func resetStrings() {
	language = ""
	theme = "monokai"
}

func invalidFileErrorMsg() string {
	if runtime.GOOS == "windows" {
		return "Error: open InvalidFilename: The system cannot find the file specified."
	}
	return "Error: open InvalidFilename: no such file or directory"
}

func _unhighlightedGoCode() string {
	if runtime.GOOS == "windows" {
		return "package main"
	}
	return unhighlightedGoCode
}
