package main

import (
	"io/ioutil"
	"os"

	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

func main() {
	filename := "README.md"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic("Read Error!")
	}

	style := styles.Get("swapoff")
	lexer := lexers.Match(filename)
	formatter := formatters.Get("terminal256")
	iterator, _ := lexer.Tokenise(nil, string(data))
	formatter.Format(os.Stdout, style, iterator)
}
