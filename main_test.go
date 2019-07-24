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
