package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainFunc(t *testing.T) {
	main()
}

func TestShell(t *testing.T) {
	t.Run("echo+pipe", func(t *testing.T) {
		cmd := exec.Command("bash", "-c", "echo pipetest | ./nyan")
		var o bytes.Buffer
		cmd.Stdout = &o
		err := cmd.Run()
		assert.Nil(t, err)
		assert.NotNil(t, o.String())
		assert.Equal(t, "pipetest\n", o.String())
	})

	t.Run("< StdInput", func(t *testing.T) {
		cmd := exec.Command("bash", "-c", "./nyan < testdata/dummyfile")
		var o bytes.Buffer
		cmd.Stdout = &o
		err := cmd.Run()
		assert.Nil(t, err)
		assert.NotNil(t, o.String())
		assert.Equal(t, "This is dummy.", o.String())
	})

	t.Run("direct input over echo+pipe", func(t *testing.T) {
		cmd := exec.Command("bash", "-c", "echo pipetest | ./nyan testdata/dummyfile")
		var o bytes.Buffer
		cmd.Stdout = &o
		err := cmd.Run()
		assert.Nil(t, err)
		assert.NotNil(t, o.String())
		assert.NotContains(t, o.String(), "pipetest")
		assert.Equal(t, "This is dummy.", o.String())
	})

	t.Run("echo+pipe & < StdInput", func(t *testing.T) {
		cmd := exec.Command("bash", "-c", "echo pipetest | ./nyan < testdata/dummyfile")
		var o bytes.Buffer
		cmd.Stdout = &o
		err := cmd.Run()
		assert.Nil(t, err)
		assert.NotNil(t, o.String())
		assert.Equal(t, "This is dummy.", o.String())
	})

	t.Run("`> file` out is not highlighted", func(t *testing.T) {
		outfile := "testdata/output"
		cmd := exec.Command("bash", "-c", "echo 'package main' | ./nyan > "+outfile)
		err := cmd.Run()
		data, err := ioutil.ReadFile(outfile)
		assert.Nil(t, err)
		assert.Equal(t, "package main\n", string(data))
	})
}
