package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandExecute(t *testing.T) {
	err := rootCmd.Execute()

	assert.Nil(t, err)
}

func TestHelpCommand(t *testing.T) {
	o := bytes.NewBufferString("")
	rootCmd.SetOut(o)
	err := rootCmd.Execute()

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
	assert.Contains(t, o.String(), "Error: open InvalidFilename: no such file or directory")
}

func TestExecute(t *testing.T) {
	o := bytes.NewBufferString("")
	rootCmd.SetArgs([]string{"testdata/dummy.go"})
	rootCmd.SetOut(o)
	err := rootCmd.Execute()

	assert.Nil(t, err)
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), "[38;5;197mpackage[0m[38;5;231m")
}

func TestInvalidTheme(t *testing.T) {
	o := bytes.NewBufferString("")
	rootCmd.SetArgs([]string{"testdata/dummy.go", "-t", "invalid"})
	rootCmd.SetOut(o)
	err := rootCmd.Execute()

	assert.Nil(t, err)
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), "[1m[38;5;231mpackage")
}

func TestValidTheme(t *testing.T) {
	o := bytes.NewBufferString("")
	rootCmd.SetArgs([]string{"testdata/dummy.go", "-t", "vim"})
	rootCmd.SetOut(o)
	err := rootCmd.Execute()

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
	o := bytes.NewBufferString("")
	i := bytes.NewBufferString("TestFromStdIn")
	rootCmd.SetArgs([]string{"-"})
	rootCmd.SetOut(o)
	rootCmd.SetIn(i)
	err := rootCmd.Execute()

	assert.Nil(t, err)
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), "TestFromStdIn")
}

func resetFlags() {
	showVersion = false
}
