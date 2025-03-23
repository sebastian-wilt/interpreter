package main

import (
	"fmt"
	"interpreter/lexer"
)

func main() {
	src := "( += );"

	lexer := lexer.NewLexer([]byte(src))

	for _, tok := range lexer.Tokenize() {
		fmt.Printf("%v\n", tok)	
	}

}
