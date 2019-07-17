package main

import (
	"io/ioutil"
	"os"

	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/spf13/cobra"
)

var showVersion bool

var rootCmd = &cobra.Command{
	Use:     "nyan [FILE]",
	Short:   "Colorized cat",
	Long:    "Colorized cat",
	Example: `$ nyan FILE`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if showVersion {
			cmd.Println("Version 0.0.0 (not yet released)")
			return nil
		}
		if len(args) < 1 {
			cmd.Help()
			return nil
		}

		var data []byte
		var err error

		filename := args[0]
		if filename == "-" {
			data, err = ioutil.ReadAll(os.Stdin)
			if err != nil {
				return err
			}
		} else {
			data, err = ioutil.ReadFile(filename)
			if err != nil {
				return err
			}
		}

		lexer := lexers.Match(filename)
		if lexer == nil {
			lexer = lexers.Fallback
		}
		iterator, _ := lexer.Tokenise(nil, string(data))
		style := styles.Get("monokai")
		formatter := formatters.Get("terminal256")
		formatter.Format(os.Stdout, style, iterator)
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, `show version`)
}

func main() {
	rootCmd.Execute()
}
