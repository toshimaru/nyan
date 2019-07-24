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
	rootCmd.SetOutput(o)
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
	rootCmd.SetOutput(o)
	err := rootCmd.Execute()

	assert.NotNil(t, err)
	assert.Contains(t, o.String(), "Error: open InvalidFilename: no such file or directory")
}

func TestExecute(t *testing.T) {
	o := bytes.NewBufferString("")
	rootCmd.SetArgs([]string{"testdata/dummy.go"})
	rootCmd.SetOutput(o)
	err := rootCmd.Execute()

	assert.Nil(t, err)
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), "[38;5;197mpackage[0m[38;5;231m")
}

func TestInvalidTheme(t *testing.T) {
	o := bytes.NewBufferString("")
	rootCmd.SetArgs([]string{"testdata/dummy.go", "-t", "invalid"})
	rootCmd.SetOutput(o)
	err := rootCmd.Execute()

	assert.Nil(t, err)
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), "[1m[38;5;231mpackage")
}

func TestValidTheme(t *testing.T) {
	o := bytes.NewBufferString("")
	rootCmd.SetArgs([]string{"testdata/dummy.go", "-t", "vim"})
	rootCmd.SetOutput(o)
	err := rootCmd.Execute()

	assert.Nil(t, err)
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), "[38;5;164mpackage[0m")
}

func TestVersionFlag(t *testing.T) {
	o := bytes.NewBufferString("")
	rootCmd.SetArgs([]string{"-v"})
	rootCmd.SetOutput(o)
	err := rootCmd.Execute()

	assert.Nil(t, err)
	assert.Contains(t, o.String(), "Version 0.0.0")
}
