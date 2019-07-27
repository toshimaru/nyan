package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

var o bytes.Buffer

func TestMain(m *testing.M) {
	o := bytes.NewBufferString("")
	rootCmd.SetOut(o)
	resetFlags()
	m.Run()
}

func TestCommandExecute(t *testing.T) {
	err := rootCmd.Execute()

	assert.Nil(t, err)
}

func TestHelpCommand(t *testing.T) {
	err := rootCmd.Execute()

	assert.Nil(t, err)
	assert.Contains(t, o.String(), rootCmd.Use)
	assert.Contains(t, o.String(), rootCmd.Short)
	assert.Contains(t, o.String(), rootCmd.Long)
	assert.Contains(t, o.String(), rootCmd.Example)
}

func TestInvalidFilename(t *testing.T) {
	rootCmd.SetArgs([]string{"InvalidFilename"})
	err := rootCmd.Execute()

	assert.NotNil(t, err)
	assert.Contains(t, o.String(), "Error: open InvalidFilename: no such file or directory")
}

func TestExecute(t *testing.T) {
	rootCmd.SetArgs([]string{"testdata/dummy.go"})
	err := rootCmd.Execute()

	assert.Nil(t, err)
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), "[38;5;197mpackage[0m[38;5;231m")
}

func TestInvalidTheme(t *testing.T) {
	rootCmd.SetArgs([]string{"testdata/dummy.go", "-t", "invalid"})
	err := rootCmd.Execute()

	assert.Nil(t, err)
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), "[1m[38;5;231mpackage")
}

func TestValidTheme(t *testing.T) {
	rootCmd.SetArgs([]string{"testdata/dummy.go", "-t", "vim"})
	err := rootCmd.Execute()

	assert.Nil(t, err)
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), "[38;5;164mpackage[0m")
}

func TestVersionFlag(t *testing.T) {
	rootCmd.SetArgs([]string{"-v"})
	err := rootCmd.Execute()

	assert.Nil(t, err)
	assert.Contains(t, o.String(), "Version 0.0.0")
}

func TestUnknownFile(t *testing.T) {
	rootCmd.SetArgs([]string{"testdata/dummyfile"})
	err := rootCmd.Execute()

	assert.Nil(t, err)
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), "This is dummy.")
}

func resetFlags() {
	showVersion = false
	theme = "monokai"
}
