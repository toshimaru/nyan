package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	data, err := ioutil.ReadFile(`README.md`)
	if err != nil {
		// エラー処理
		panic("Error!")
	}
	fmt.Print(string(data))
}
