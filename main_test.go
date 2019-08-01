package main

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandExecute(t *testing.T) {
	err := rootCmd.Execute()

	assert.Nil(t, err)
}

func TestHelpCommand(t *testing.T) {
	o := bytes.NewBufferString("")
	rootCmd.SetArgs([]string{"--help"})
	rootCmd.SetOut(o)
	err := rootCmd.Execute()
	resetFlags()

	assert.Nil(t, err)
	assert.Contains(t, o.String(), rootCmd.Use)
	assert.Contains(t, o.String(), rootCmd.Short)
	assert.Contains(t, o.String(), rootCmd.Long)
	assert.Contains(t, o.String(), rootCmd.Example)
}

func TestInvalidFilename(t *testing.T) {
	o := bytes.NewBufferString("")
	rootCmd.SetArgs([]string{"InvalidFilename"})
	rootCmd.SetOut(o)
	err := rootCmd.Execute()

	assert.NotNil(t, err)
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), "Error: open InvalidFilename: no such file or directory")
}

func TestExecute(t *testing.T) {
	o := bytes.NewBufferString("")
	isTerminalFunc = func(fd uintptr) bool { return true }
	rootCmd.SetArgs([]string{"testdata/dummy.go"})
	rootCmd.SetOut(o)
	err := rootCmd.Execute()

	assert.Nil(t, err)
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), "[38;5;197mpackage[0m[38;5;231m")
}

func TestInvalidTheme(t *testing.T) {
	o := bytes.NewBufferString("")
	isTerminalFunc = func(fd uintptr) bool { return true }
	rootCmd.SetArgs([]string{"testdata/dummy.go", "-t", "invalid"})
	rootCmd.SetOut(o)
	err := rootCmd.Execute()

	assert.Nil(t, err)
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), "[1m[38;5;231mpackage")
}

func TestValidTheme(t *testing.T) {
	o := bytes.NewBufferString("")
	isTerminalFunc = func(fd uintptr) bool { return true }
	rootCmd.SetArgs([]string{"testdata/dummy.go", "-t", "vim"})
	rootCmd.SetOut(o)
	err := rootCmd.Execute()
	resetTheme()

	assert.Nil(t, err)
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), "[38;5;164mpackage[0m")
}

func TestVersionFlag(t *testing.T) {
	o := bytes.NewBufferString("")
	rootCmd.SetArgs([]string{"-v"})
	rootCmd.SetOut(o)
	err := rootCmd.Execute()
	resetFlags()

	assert.Nil(t, err)
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), "Version 0.0.0")
}

func TestUnknownFile(t *testing.T) {
	o := bytes.NewBufferString("")
	rootCmd.SetArgs([]string{"testdata/dummyfile"})
	rootCmd.SetOut(o)
	err := rootCmd.Execute()

	assert.Nil(t, err)
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), "This is dummy.")
}

func TestFromStdIn(t *testing.T) {
	i := bytes.NewBufferString("package main")
	o := bytes.NewBufferString("")
	isTerminalFunc = func(fd uintptr) bool { return true }
	rootCmd.SetArgs([]string{})
	rootCmd.SetIn(i)
	rootCmd.SetOut(o)
	err := rootCmd.Execute()

	assert.Nil(t, err)
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), "[38;5;197mpackage[0m[38;5;231m")
}

func TestFromStdInWithDash(t *testing.T) {
	i := bytes.NewBufferString("TestFromStdIn")
	o := bytes.NewBufferString("")
	rootCmd.SetArgs([]string{"-"})
	rootCmd.SetIn(i)
	rootCmd.SetOut(o)
	err := rootCmd.Execute()

	assert.Nil(t, err)
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), "TestFromStdIn")
}

func TestShell(t *testing.T) {
	t.Run("echo+pipe", func(t *testing.T) {
		cmd := exec.Command("bash", "-c", "echo pipetest | ./nyan")
		var o bytes.Buffer
		cmd.Stdout = &o
		err := cmd.Run()
		assert.Nil(t, err)
		assert.NotNil(t, o.String())
		assert.Contains(t, o.String(), "pipetest")
	})

	t.Run("< StdInput", func(t *testing.T) {
		cmd := exec.Command("bash", "-c", "./nyan < testdata/dummyfile")
		var o bytes.Buffer
		cmd.Stdout = &o
		err := cmd.Run()
		assert.Nil(t, err)
		assert.NotNil(t, o.String())
		assert.Contains(t, o.String(), "This is dummy.")
	})

	t.Run("direct input over echo+pipe", func(t *testing.T) {
		cmd := exec.Command("bash", "-c", "echo pipetest | ./nyan testdata/dummyfile")
		var o bytes.Buffer
		cmd.Stdout = &o
		err := cmd.Run()
		assert.Nil(t, err)
		assert.NotNil(t, o.String())
		assert.NotContains(t, o.String(), "pipetest")
		assert.Contains(t, o.String(), "This is dummy.")
	})

	t.Run("echo+pipe & < StdInput", func(t *testing.T) {
		cmd := exec.Command("bash", "-c", "echo pipetest | ./nyan < testdata/dummyfile")
		var o bytes.Buffer
		cmd.Stdout = &o
		err := cmd.Run()
		assert.Nil(t, err)
		assert.NotNil(t, o.String())
		assert.Contains(t, o.String(), "This is dummy.")
	})
}

func resetFlags() {
	showVersion = false
	rootCmd.Flags().Set("help", "false")
}

func resetTheme() {
	theme = "monokai"
}
