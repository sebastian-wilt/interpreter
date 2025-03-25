package main

import (
	"fmt"
	"interpreter/lexer"
)

func main() {
	src := "if else false true for in while fun return val var continue fall match"
	// src := "( += );"

	lexer := lexer.NewLexer([]byte(src))

	for _, tok := range lexer.Tokenize() {
		fmt.Printf("%v\n", tok)	
	}

}
