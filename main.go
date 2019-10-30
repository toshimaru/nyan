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
	language    string
)

var rootCmd = &cobra.Command{
	Use:   "nyan [OPTION]... [FILE]...",
	Short: "Colored cat command.",
	Long:  "Colored cat command which supports syntax highlighting.",
	Example: `$ nyan FILE
$ nyan FILE1 FILE2 FILE3
$ nyan -t solarized-dark FILE`,
	RunE:          cmdMain,
	SilenceErrors: true,
	SilenceUsage:  false,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&listThemes, "list-themes", "T", false, `List available color themes`)
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, `Show version`)
	rootCmd.PersistentFlags().StringVarP(&theme, "theme", "t", "monokai", fmt.Sprintf("Set color theme for syntax highlighting\nAvailable themes: %s", styles.Names()))
	rootCmd.PersistentFlags().StringVarP(&language, "language", "l", "", "Specify language for syntax highlighting")

	rootCmd.SetOut(colorable.NewColorableStdout())
	rootCmd.SetErr(colorable.NewColorableStderr())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
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
	var lexer chroma.Lexer

	if len(args) < 1 || args[0] == "-" {
		if data, err = ioutil.ReadAll(cmd.InOrStdin()); err != nil {
			cmd.PrintErr(err, "\n")
			return
		}
		if language != "" {
			lexer = lexers.Get(language)
		} else {
			lexer = lexers.Analyse(string(data))
		}
		printData(&data, cmd, lexer)
	} else {
		var lastErr error
		for _, filename := range args {
			if data, err = ioutil.ReadFile(filename); err != nil {
				// FIXME: use PrintErrln after upstream is fixed
				cmd.PrintErr(err, "\n")
				lastErr = err
				continue
			}
			if language != "" {
				lexer = lexers.Get(language)
			} else {
				lexer = lexers.Match(filename)
			}
			printData(&data, cmd, lexer)
		}
		if lastErr != nil {
			cmd.SilenceUsage = true
			return lastErr
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

const sampleCode = `// Sample Code in Go
package main

import "fmt"

func main() {
	fmt.Println("Hello nyan cat command 😺")
}
`

func printThemes(cmd *cobra.Command) {
	for _, theme = range styles.Names() {
		cmd.Println("Theme:", theme)
		code := []byte(sampleCode)
		lexer := lexers.Get("go")
		printData(&code, cmd, lexer)
		cmd.Println()
	}
}
