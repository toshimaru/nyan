package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
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
	number      bool
)

var rootCmd = &cobra.Command{
	Use:   "nyan [flags] [FILE]...",
	Short: "Colored cat command.",
	Long:  "Colored cat command which supports syntax highlighting.",
	Example: `$ nyan FILE
$ nyan FILE1 FILE2 FILE3
$ nyan -t solarized-dark FILE
$ nyan -l go FILE`,
	RunE:          cmdMain,
	SilenceErrors: true,
	SilenceUsage:  false,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&listThemes, "list-themes", "T", false, `List available color themes`)
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, `Show version`)
	rootCmd.PersistentFlags().StringVarP(&theme, "theme", "t", "monokai", fmt.Sprintf("Set color theme for syntax highlighting\nAvailable themes: %s", styles.Names()))
	rootCmd.PersistentFlags().StringVarP(&language, "language", "l", "", "Specify language for syntax highlighting")
	rootCmd.PersistentFlags().BoolVarP(&number, "number", "n", false, "Output with line numbers")

	rootCmd.SetOut(colorable.NewColorableStdout())
	rootCmd.SetErr(colorable.NewColorableStderr())
}

// Execute root commands. This is called by `main.main()`.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func cmdMain(cmd *cobra.Command, args []string) (err error) {
	if checkSpecialFlags(cmd) {
		return
	}

	var data []byte
	var lexer chroma.Lexer

	if language != "" {
		lexer = lexers.Get(language)
	}

	cmd.SilenceUsage = true

	if len(args) < 1 || args[0] == "-" {
		if data, err = io.ReadAll(cmd.InOrStdin()); err != nil {
			cmd.PrintErrln("Error:", err)
			return err
		}
		if lexer == nil {
			lexer = lexers.Analyse(string(data))
		}
		printData(&data, cmd, lexer)
	} else {
		var lastErr error
		for _, filename := range args {
			if data, err = os.ReadFile(filename); err != nil {
				cmd.PrintErrln("Error:", err)
				lastErr = err
				continue
			}
			if language == "" {
				lexer = lexers.Match(filename)
			}
			printData(&data, cmd, lexer)
		}
		if lastErr != nil {
			return lastErr
		}
	}
	return
}

func printData(data *[]byte, cmd *cobra.Command, lexer chroma.Lexer) {
	out := cmd.OutOrStdout()
	if number {
		w := &numberWriter{
			w:           out,
			currentLine: 1,
		}
		out = w
		defer w.Flush()
	}

	if isTerminalFunc(os.Stdout.Fd()) {
		if lexer == nil {
			lexer = lexers.Fallback
		}
		iterator, _ := lexer.Tokenise(nil, string(*data))
		formatter := formatters.Get("terminal256")
		formatter.Format(out, styles.Get(theme), iterator)
	} else {
		fmt.Fprint(out, string(*data))
	}
}

func checkSpecialFlags(cmd *cobra.Command) bool {
	if showVersion {
		cmd.Println("version", version)
		return true
	} else if listThemes {
		printThemes(cmd)
		return true
	}
	return false
}

const sampleCode = `
  // Sample Code in Go
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

type numberWriter struct {
	w           io.Writer
	currentLine uint64
	buf         []byte
}

func (w *numberWriter) Write(p []byte) (n int, err error) {
	// Early return.
	// Can't calculate the line numbers until the line breaks are made, so store them all in a buffer.
	if !bytes.Contains(p, []byte{'\n'}) {
		w.buf = append(w.buf, p...)
		return len(p), nil
	}

	var (
		original = p
		tokenLen uint
	)
	for i, c := range original {
		tokenLen++
		if c != '\n' {
			continue
		}

		token := p[:tokenLen]
		p = original[i+1:]
		tokenLen = 0

		format := "%6d\t%s%s"
		if w.currentLine > 999999 {
			format = "%d\t%s%s"
		}

		_, er := fmt.Fprintf(w.w, format, w.currentLine, string(w.buf), string(token))
		if er != nil {
			return i + 1, er
		}
		w.buf = w.buf[:0]
		w.currentLine++
	}

	if len(p) > 0 {
		w.buf = append(w.buf, p...)
	}
	return len(original), nil
}

func (w *numberWriter) Flush() error {
	terminalReset := []byte("\u001B[0m")
	if bytes.Compare(w.buf, terminalReset) == 0 {
		// In almost all cases, a control code is passed last to reset the terminal's color code.
		// This is not a printable character and should not be counted as a line, so it is output as is without a line number.
		_, err := fmt.Fprintf(w.w, "%s", string(w.buf))
		return err
	}

	format := "%6d\t%s"
	if w.currentLine > 999999 {
		format = "%d\t%s"
	}
	_, err := fmt.Fprintf(w.w, format, w.currentLine, string(w.buf))
	w.buf = w.buf[:0]
	return err
}
