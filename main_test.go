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
	highlightedGoCode   = "[38;5;197mpackage[0m[38;5;231m"
	unhighlightedGoCode = "[38;5;231mpackage main[0m[38;5;231m"
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

func TestExecute(t *testing.T) {
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

func TestExecuteWithAnalyseUnknownFile(t *testing.T) {
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
	rootCmd.SetArgs([]string{"testdata/dummy.go", "testdata/dummyfile"})
	rootCmd.SetOut(&o)
	rootCmd.SetErr(&e)
	err := rootCmd.Execute()

	assert.Nil(t, err)
	assert.Empty(t, e.String())
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), highlightedGoCode)
	assert.Contains(t, o.String(), "[0m[38;5;231mThis is dummy.[0m")
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
func TestInvalidTheme(t *testing.T) {
	var o, e bytes.Buffer
	isTerminalFunc = func(fd uintptr) bool { return true }
	rootCmd.SetArgs([]string{"testdata/dummy.go", "--theme", "invalid"})
	rootCmd.SetOut(&o)
	rootCmd.SetErr(&e)
	err := rootCmd.Execute()
	resetStrings()

	assert.Nil(t, err)
	assert.Empty(t, e.String())
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), "[1m[38;5;231mpackage")
}

func TestValidTheme(t *testing.T) {
	var o, e bytes.Buffer
	isTerminalFunc = func(fd uintptr) bool { return true }
	rootCmd.SetArgs([]string{"testdata/dummy.go", "--theme", "vim"})
	rootCmd.SetOut(&o)
	rootCmd.SetErr(&e)
	err := rootCmd.Execute()
	resetStrings()

	assert.Nil(t, err)
	assert.Empty(t, e.String())
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), "[38;5;164mpackage[0m")
}

func TestVersionFlag(t *testing.T) {
	var o, e bytes.Buffer
	rootCmd.SetArgs([]string{"-v"})
	rootCmd.SetOut(&o)
	rootCmd.SetErr(&e)
	err := rootCmd.Execute()
	resetFlags()

	assert.Nil(t, err)
	assert.Empty(t, e.String())
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), "version ")
}

func TestListThemesFlag(t *testing.T) {
	var o, e bytes.Buffer
	rootCmd.SetArgs([]string{"--list-themes"})
	rootCmd.SetOut(&o)
	rootCmd.SetErr(&e)
	err := rootCmd.Execute()
	resetFlags()

	assert.Nil(t, err)
	assert.Empty(t, e.String())
	assert.NotNil(t, o.String())
	assert.Contains(t, o.String(), "Theme: ")
	assert.Contains(t, o.String(), "Sample Code in Go")
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

func resetStrings() {
	language = ""
	theme = "monokai"
}

func invalidFileErrorMsg() string {
	if runtime.GOOS == "windows" {
		return "open InvalidFilename: The system cannot find the file specified."
	}
	return "open InvalidFilename: no such file or directory"
}

func _unhighlightedGoCode() string {
	if runtime.GOOS == "windows" {
		return "package main"
	}
	return unhighlightedGoCode
}
