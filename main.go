package main

import (
	"io/ioutil"
	"os"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"github.com/toshimaru/nyan/styles"
)

var (
	showVersion bool
	theme       string
)

var rootCmd = &cobra.Command{
	Use:     "nyan [FILE]",
	Short:   "Colorized cat",
	Long:    "Colorized cat",
	Example: `$ nyan FILE`,
	RunE:    cmdMain,
}

var isTerminalFunc = isatty.IsTerminal

func init() {
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, `show version`)
	rootCmd.PersistentFlags().StringVarP(&theme, "theme", "t", "monokai", "color theme")
}

func main() {
	rootCmd.SetOutput(os.Stdout)
	if err := rootCmd.Execute(); err != nil {
		rootCmd.SetOutput(os.Stderr)
		rootCmd.Println(err)
		os.Exit(1)
	}
}

func cmdMain(cmd *cobra.Command, args []string) (err error) {
	if showVersion {
		cmd.Println("Version 0.0.0 (not yet released)")
		return
	}

	var data []byte
	var lexer chroma.Lexer

	if len(args) < 1 || args[0] == "-" {
		if data, err = ioutil.ReadAll(cmd.InOrStdin()); err != nil {
			return
		}
		lexer = lexers.Analyse(string(data))
	} else {
		filename := args[0]
		if data, err = ioutil.ReadFile(filename); err != nil {
			return
		}
		lexer = lexers.Match(filename)
	}

	if isTerminalFunc(os.Stdout.Fd()) {
		if lexer == nil {
			lexer = lexers.Fallback
		}
		iterator, _ := lexer.Tokenise(nil, string(data))
		formatter := formatters.Get("terminal256")
		formatter.Format(cmd.OutOrStdout(), styles.Get(theme), iterator)
	} else {
		cmd.Print(string(data))
	}
	return
}
