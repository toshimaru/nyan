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
	Run: func(cmd *cobra.Command, args []string) {
		if showVersion {
			cmd.Println("Version 0.0.0 (not yet released)")
			return
		}

		var filename string
		var data []byte
		var err error

		if len(args) < 1 {
			cmd.Help()
			return
		} else if false {
			// TODO
			filename = "-"
			data, err = ioutil.ReadAll(os.Stdin)
			if err != nil {
				panic("Read Error!")
			}
		} else {
			filename = args[0]
			data, err = ioutil.ReadFile(filename)
			if err != nil {
				panic("Read Error!")
			}
		}

		style := styles.Get("swapoff")
		lexer := lexers.Match(filename)
		formatter := formatters.Get("terminal256")
		iterator, _ := lexer.Tokenise(nil, string(data))
		formatter.Format(os.Stdout, style, iterator)
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, `show version`)
}

func main() {
	rootCmd.Execute()
}
