package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	highlightedGoCode   = "[38;5;197mpackage[0m[38;5;231m"
	unhighlightedGoCode = "[38;5;231mpackage main[0m[38;5;231m"
)

func TestMain(t *testing.T) {
	rootCmd.SetArgs([]string{"--help"})
	main()
	resetFlags()
}

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
	assert.Contains(t, o.String(), invalidFileErrorMsg())
}

func TestExecute(t *testing.T) {
	o := bytes.NewBufferString("")
	isTerminalFunc = func(fd uintptr) bool { return true }
	rootCmd.SetArgs([]string{"testdata/dummy.go"})
	rootCmd.SetOut(o)
	err := rootCmd.Execute()

	assert.Nil(t, err)
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), highlightedGoCode)
}

func TestExecuteWithAnalyseUnknownFile(t *testing.T) {
	o := bytes.NewBufferString("")
	isTerminalFunc = func(fd uintptr) bool { return true }
	rootCmd.SetArgs([]string{"testdata/dummy.go.unknown"})
	rootCmd.SetOut(o)
	err := rootCmd.Execute()

	assert.Nil(t, err)
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), unhighlightedGoCode)
}

func TestMultipleFiles(t *testing.T) {
	o := bytes.NewBufferString("")
	isTerminalFunc = func(fd uintptr) bool { return true }
	rootCmd.SetArgs([]string{"testdata/dummy.go", "testdata/dummyfile"})
	rootCmd.SetOut(o)
	err := rootCmd.Execute()

	assert.Nil(t, err)
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), highlightedGoCode)
	assert.Contains(t, o.String(), "[0m[38;5;231mThis is dummy.[0m")
}

func TestMultipleFilesWithInvalidFileError(t *testing.T) {
	o := bytes.NewBufferString("")
	isTerminalFunc = func(fd uintptr) bool { return true }
	rootCmd.SetArgs([]string{"testdata/dummy.go", "InvalidFilename", "testdata/dummyfile"})
	rootCmd.SetOut(o)
	err := rootCmd.Execute()

	assert.Nil(t, err)
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), highlightedGoCode)
	assert.Contains(t, o.String(), invalidFileErrorMsg())
	assert.Contains(t, o.String(), "[38;5;231mThis is dummy.[0m")
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
	assert.Contains(t, o.String(), "version ")
}

func TestListThemesFlag(t *testing.T) {
	o := bytes.NewBufferString("")
	rootCmd.SetArgs([]string{"--list-themes"})
	rootCmd.SetOut(o)
	err := rootCmd.Execute()
	resetFlags()

	assert.Nil(t, err)
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), "Theme: ")
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
	rootCmd.SetArgs([]string{"-t", "monokai"})
	rootCmd.SetIn(i)
	rootCmd.SetOut(o)
	err := rootCmd.Execute()

	assert.Nil(t, err)
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), highlightedGoCode)
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
		outfile := "testdata/output.out"
		cmd := exec.Command("bash", "-c", "echo 'package main' | ./nyan > "+outfile)
		err := cmd.Run()
		data, err := ioutil.ReadFile(outfile)
		assert.Nil(t, err)
		assert.Equal(t, "package main\n", string(data))
	})
}

func resetFlags() {
	showVersion = false
	listThemes = false
	rootCmd.Flags().Set("help", "false")
}

func resetTheme() {
	theme = "monokai"
}

func invalidFileErrorMsg() string {
	if runtime.GOOS == "windows" {
		return "open InvalidFilename: The system cannot find the file specified."
	}
	return "open InvalidFilename: no such file or directory"
}
