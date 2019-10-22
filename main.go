package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"github.com/toshimaru/nyan/styles"
)

var (
	isTerminalFunc = isatty.IsTerminal
	version        = "dev"

	listThemes  bool
	showVersion bool
	theme       string
)

var rootCmd = &cobra.Command{
	Use:   "nyan [OPTION]... [FILE]...",
	Short: "Colored cat command.",
	Long:  "Colored cat command which supports syntax highlighting.",
	Example: `$ nyan FILE
$ nyan FILE1 FILE2
$ nyan -t solarized-dark FILE1`,
	RunE: cmdMain,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, `show version`)
	rootCmd.PersistentFlags().StringVarP(&theme, "theme", "t", "monokai", fmt.Sprintf("color theme\navailable themes: %s", styles.Names()))
	rootCmd.PersistentFlags().BoolVarP(&listThemes, "list-themes", "T", false, `list available themes`)
}

func main() {
	rootCmd.SetOutput(colorable.NewColorableStdout())
	if err := rootCmd.Execute(); err != nil {
		rootCmd.SetOutput(os.Stderr)
		rootCmd.Println(err)
		os.Exit(1)
	}
}

func cmdMain(cmd *cobra.Command, args []string) (err error) {
	if showVersion {
		cmd.Println("version", version)
		return
	} else if listThemes {
		printThemes(cmd)
		return
	}

	var data []byte

	if len(args) < 1 || args[0] == "-" {
		if data, err = ioutil.ReadAll(cmd.InOrStdin()); err != nil {
			return
		}
		lexer := lexers.Analyse(string(data))
		printData(&data, cmd, lexer)
	} else {
		for _, filename := range args {
			if data, err = ioutil.ReadFile(filename); err != nil {
				cmd.Println(err)
			}
			lexer := lexers.Match(filename)
			printData(&data, cmd, lexer)
		}
	}

	return
}

func printData(data *[]byte, cmd *cobra.Command, lexer chroma.Lexer) {
	if isTerminalFunc(os.Stdout.Fd()) {
		if lexer == nil {
			lexer = lexers.Fallback
		}
		iterator, _ := lexer.Tokenise(nil, string(*data))
		formatter := formatters.Get("terminal256")
		formatter.Format(cmd.OutOrStdout(), styles.Get(theme), iterator)
	} else {
		cmd.Print(string(*data))
	}
}

const sampleCode = `package main

import "fmt"

func main() {
	fmt.Println("Hello nyan cat command 😺")
}
`

func printThemes(cmd *cobra.Command) {
	for _, theme = range styles.Names() {
		cmd.Println("Theme: ", theme)
		code := []byte(sampleCode)
		lexer := lexers.Get("go")
		printData(&code, cmd, lexer)
		cmd.Println()
	}
}
