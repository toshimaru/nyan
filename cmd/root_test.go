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
	t.Cleanup(resetFlags)
	rootCmd.SetArgs([]string{"--help"})
	Execute()
}

func TestCommandExecute(t *testing.T) {
	err := rootCmd.Execute()

	assert.NoError(t, err)
}

func TestHelpCommand(t *testing.T) {
	t.Cleanup(resetFlags)
	var o, e bytes.Buffer
	rootCmd.SetArgs([]string{"--help"})
	rootCmd.SetOut(&o)
	rootCmd.SetErr(&e)
	err := rootCmd.Execute()

	assert.NoError(t, err)
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
	setupTerminalMock(t)
	var o, e bytes.Buffer
	rootCmd.SetArgs([]string{"testdata/dummy.go"})
	rootCmd.SetOut(&o)
	rootCmd.SetErr(&e)
	err := rootCmd.Execute()

	assert.NoError(t, err)
	assert.Empty(t, e.String())
	assert.NotEmpty(t, o.String())
	assert.Contains(t, o.String(), highlightedGoCode)
}

func TestUnknownExtension(t *testing.T) {
	setupTerminalMock(t)
	var o, e bytes.Buffer
	rootCmd.SetArgs([]string{"testdata/dummy.go.unknown"})
	rootCmd.SetOut(&o)
	rootCmd.SetErr(&e)
	err := rootCmd.Execute()

	assert.NoError(t, err)
	assert.Empty(t, e.String())
	assert.NotEmpty(t, o.String())
	assert.Contains(t, o.String(), _unhighlightedGoCode())
}

func TestLanguageOption(t *testing.T) {
	setupTerminalMockWithStrings(t)
	var o, e bytes.Buffer
	rootCmd.SetArgs([]string{"--language", "go", "testdata/dummy.go.unknown"})
	rootCmd.SetOut(&o)
	rootCmd.SetErr(&e)
	err := rootCmd.Execute()

	assert.NoError(t, err)
	assert.Empty(t, e.String())
	assert.NotEmpty(t, o.String())
	assert.Contains(t, o.String(), highlightedGoCode)
}

func TestInvalidLanguageOption(t *testing.T) {
	setupTerminalMockWithStrings(t)
	var o, e bytes.Buffer
	rootCmd.SetArgs([]string{"--language", "invalid_lang", "testdata/dummy.go"})
	rootCmd.SetOut(&o)
	rootCmd.SetErr(&e)
	err := rootCmd.Execute()

	assert.NoError(t, err)
	assert.Empty(t, e.String())
	assert.NotEmpty(t, o.String())
	assert.Contains(t, o.String(), _unhighlightedGoCode())
}

func TestMultipleFiles(t *testing.T) {
	setupTerminalMock(t)
	var o, e bytes.Buffer
	rootCmd.SetArgs([]string{"testdata/dummy.go", "testdata/dummy.go.unknown"})
	rootCmd.SetOut(&o)
	rootCmd.SetErr(&e)
	err := rootCmd.Execute()

	assert.NoError(t, err)
	assert.Empty(t, e.String())
	assert.NotEmpty(t, o.String())
	assert.Contains(t, o.String(), highlightedGoCode)
	assert.Contains(t, o.String(), _unhighlightedGoCode())
}

func TestMultipleFilesWithInvalidFileError(t *testing.T) {
	setupTerminalMock(t)
	var o, e bytes.Buffer
	rootCmd.SetArgs([]string{"testdata/dummy.go", "InvalidFilename", "testdata/dummyfile"})
	rootCmd.SetOut(&o)
	rootCmd.SetErr(&e)
	err := rootCmd.Execute()

	assert.Error(t, err)
	assert.NotEmpty(t, o.String())
	assert.Contains(t, o.String(), highlightedGoCode)
	assert.Contains(t, o.String(), "[38;5;231mThis is dummy.[0m")
	assert.Contains(t, e.String(), invalidFileErrorMsg())
}

func TestCompletionDisabled(t *testing.T) {
	var o, e bytes.Buffer
	rootCmd.SetArgs([]string{"completion"})
	rootCmd.SetOut(&o)
	rootCmd.SetErr(&e)
	err := rootCmd.Execute()

	assert.Error(t, err)
	assert.Contains(t, e.String(), "Error: open completion:")
	assert.Empty(t, o.String())
}

func TestThemes(t *testing.T) {
	setupTerminalMock(t)
	var o, e bytes.Buffer
	rootCmd.SetOut(&o)
	rootCmd.SetErr(&e)

	t.Run("Valid Theme", func(t *testing.T) {
		t.Cleanup(resetStrings)
		rootCmd.SetArgs([]string{"testdata/dummy.go", "--theme", "vim"})
		err := rootCmd.Execute()

		assert.NoError(t, err)
		assert.Empty(t, e.String())
		assert.NotEmpty(t, o.String())
		assert.Contains(t, o.String(), "[38;5;164mpackage[0m")
	})

	t.Run("Invalid Theme", func(t *testing.T) {
		t.Cleanup(resetStrings)
		o.Reset()
		rootCmd.SetArgs([]string{"testdata/dummy.go", "--theme", "invalid"})
		err := rootCmd.Execute()

		assert.NoError(t, err)
		assert.Empty(t, e.String())
		assert.NotEmpty(t, o.String())
		assert.Contains(t, o.String(), "[1m[38;5;231mpackage")
	})
}

func TestSpecialFlags(t *testing.T) {
	var o, e bytes.Buffer
	rootCmd.SetOut(&o)
	rootCmd.SetErr(&e)

	t.Run("version Flag", func(t *testing.T) {
		t.Cleanup(resetFlags)
		rootCmd.SetArgs([]string{"--version"})
		err := rootCmd.Execute()

		assert.NoError(t, err)
		assert.Empty(t, e.String())
		assert.NotEmpty(t, o.String())
		assert.Contains(t, o.String(), "version ")
		assert.NotContains(t, o.String(), "Theme: ")
	})

	t.Run("listThemes Flag", func(t *testing.T) {
		t.Cleanup(resetFlags)
		o.Reset()
		rootCmd.SetArgs([]string{"--list-themes"})
		err := rootCmd.Execute()

		assert.NoError(t, err)
		assert.Empty(t, e.String())
		assert.NotEmpty(t, o.String())
		assert.Contains(t, o.String(), "Theme: ")
		assert.Contains(t, o.String(), "Sample Code in Go")
		assert.NotContains(t, o.String(), "version ")
	})

	t.Run("multiple flags", func(t *testing.T) {
		t.Cleanup(resetFlags)
		o.Reset()
		rootCmd.SetArgs([]string{"--version", "--list-themes"})
		err := rootCmd.Execute()

		assert.NoError(t, err)
		assert.Empty(t, e.String())
		assert.NotEmpty(t, o.String())
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

	assert.NoError(t, err)
	assert.Empty(t, e.String())
	assert.NotEmpty(t, o.String())
	assert.Contains(t, o.String(), "This is dummy.")
}

func TestFromStdIn(t *testing.T) {
	setupTerminalMock(t)
	i := bytes.NewBufferString("package main")
	var o, e bytes.Buffer
	rootCmd.SetArgs([]string{"--theme", "monokai"})
	rootCmd.SetIn(i)
	rootCmd.SetOut(&o)
	rootCmd.SetErr(&e)
	err := rootCmd.Execute()

	assert.NoError(t, err)
	assert.Empty(t, e.String())
	assert.NotEmpty(t, o.String())
	assert.Contains(t, o.String(), highlightedGoCode)
}

func TestFromStdInWithLanguageOption(t *testing.T) {
	setupTerminalMockWithStrings(t)
	i := bytes.NewBufferString("package main")
	var o, e bytes.Buffer
	rootCmd.SetArgs([]string{"--theme", "monokai", "--language", "go"})
	rootCmd.SetIn(i)
	rootCmd.SetOut(&o)
	rootCmd.SetErr(&e)
	err := rootCmd.Execute()

	assert.NoError(t, err)
	assert.Empty(t, e.String())
	assert.NotEmpty(t, o.String())
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

	assert.NoError(t, err)
	assert.Empty(t, e.String())
	assert.NotEmpty(t, o.String())
	assert.Contains(t, o.String(), "TestFromStdIn")
}

func TestNumberOption(t *testing.T) {
	var o, e bytes.Buffer
	rootCmd.SetArgs([]string{"-n", "testdata/dummy.go"})
	rootCmd.SetIn(nil)
	rootCmd.SetOut(&o)
	rootCmd.SetErr(&e)
	err := rootCmd.Execute()

	assert.NoError(t, err)

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

func setupTerminalMock(t *testing.T) {
	t.Helper()
	originalIsTerminalFunc := isTerminalFunc
	isTerminalFunc = func(fd uintptr) bool { return true }
	t.Cleanup(func() {
		isTerminalFunc = originalIsTerminalFunc
	})
}

func setupTerminalMockWithStrings(t *testing.T) {
	t.Helper()
	setupTerminalMock(t)
	t.Cleanup(resetStrings)
}

func resetFlags() {
	showVersion = false
	listThemes = false
	number = false
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

func TestDetectShebang(t *testing.T) {
	t.Run("Bash shebang", func(t *testing.T) {
		data := []byte("#!/bin/bash\necho 'hello'")
		lexer := detectShebang(data)
		assert.NotNil(t, lexer)
		assert.Contains(t, lexer.Config().Name, "Bash")
	})

	t.Run("Python3 with env", func(t *testing.T) {
		data := []byte("#!/usr/bin/env python3\nprint('hello')")
		lexer := detectShebang(data)
		assert.NotNil(t, lexer)
		assert.Contains(t, lexer.Config().Name, "Python")
	})

	t.Run("Ruby shebang", func(t *testing.T) {
		data := []byte("#!/usr/bin/ruby\nputs 'hello'")
		lexer := detectShebang(data)
		assert.NotNil(t, lexer)
		assert.Contains(t, lexer.Config().Name, "Ruby")
	})

	t.Run("Node with env", func(t *testing.T) {
		data := []byte("#!/usr/bin/env node\nconsole.log('hello')")
		lexer := detectShebang(data)
		assert.NotNil(t, lexer)
		assert.Contains(t, lexer.Config().Name, "JavaScript")
	})

	t.Run("Sh shebang", func(t *testing.T) {
		data := []byte("#!/bin/sh\necho 'hello'")
		lexer := detectShebang(data)
		assert.NotNil(t, lexer)
		assert.Contains(t, lexer.Config().Name, "Bash")
	})

	t.Run("Perl shebang", func(t *testing.T) {
		data := []byte("#!/usr/bin/perl\nprint 'hello'")
		lexer := detectShebang(data)
		assert.NotNil(t, lexer)
		assert.Contains(t, lexer.Config().Name, "Perl")
	})

	t.Run("No shebang", func(t *testing.T) {
		data := []byte("echo 'hello'")
		lexer := detectShebang(data)
		assert.Nil(t, lexer)
	})

	t.Run("Comment but not shebang", func(t *testing.T) {
		data := []byte("# This is a comment\necho 'hello'")
		lexer := detectShebang(data)
		assert.Nil(t, lexer)
	})

	t.Run("Unknown interpreter", func(t *testing.T) {
		data := []byte("#!/bin/unknowninterpreter\necho 'hello'")
		lexer := detectShebang(data)
		assert.Nil(t, lexer)
	})

	t.Run("Shebang with arguments", func(t *testing.T) {
		data := []byte("#!/bin/bash -e\necho 'hello'")
		lexer := detectShebang(data)
		assert.NotNil(t, lexer)
		assert.Contains(t, lexer.Config().Name, "Bash")
	})

	t.Run("Empty file", func(t *testing.T) {
		data := []byte("")
		lexer := detectShebang(data)
		assert.Nil(t, lexer)
	})
}

func TestShebangFileDetection(t *testing.T) {
	setupTerminalMock(t)

	t.Run("Bash script without extension", func(t *testing.T) {
		var o, e bytes.Buffer
		rootCmd.SetArgs([]string{"testdata/bashscript"})
		rootCmd.SetOut(&o)
		rootCmd.SetErr(&e)
		err := rootCmd.Execute()

		assert.NoError(t, err)
		assert.Empty(t, e.String())
		assert.NotEmpty(t, o.String())
		// Should be syntax highlighted as bash
		assert.Contains(t, o.String(), "[38;5;")
	})

	t.Run("Python script without extension", func(t *testing.T) {
		t.Cleanup(resetStrings)
		var o, e bytes.Buffer
		rootCmd.SetArgs([]string{"testdata/pythonscript"})
		rootCmd.SetOut(&o)
		rootCmd.SetErr(&e)
		err := rootCmd.Execute()

		assert.NoError(t, err)
		assert.Empty(t, e.String())
		assert.NotEmpty(t, o.String())
		// Should be syntax highlighted as python
		assert.Contains(t, o.String(), "[38;5;")
	})
}

func TestShebangFromStdin(t *testing.T) {
	setupTerminalMock(t)

	t.Run("Bash script from stdin", func(t *testing.T) {
		t.Cleanup(resetStrings)
		i := bytes.NewBufferString("#!/bin/bash\necho 'hello'")
		var o, e bytes.Buffer
		rootCmd.SetArgs([]string{})
		rootCmd.SetIn(i)
		rootCmd.SetOut(&o)
		rootCmd.SetErr(&e)
		err := rootCmd.Execute()

		assert.NoError(t, err)
		assert.Empty(t, e.String())
		assert.NotEmpty(t, o.String())
		// Should be syntax highlighted
		assert.Contains(t, o.String(), "[38;5;")
	})
}
