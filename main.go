package main

import (
	"fmt"
	"interpreter/lexer"
	"interpreter/parser"
	"os"

	"github.com/chzyer/readline"
)

func main() {
	repl()
}

func repl() {
	rl, err := readline.New("> ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open repl")
	}

	for {
		line, err := rl.Readline()
		if err != nil {
			return
		}

		lexer := lexer.NewLexer([]byte(line))
		tokens, errors := lexer.Tokenize()
		if errors != nil {
			for _, err := range errors {
				fmt.Printf("Error: %s\n", err)
			}
		}

		fmt.Printf("Tokens: \n")
		for _, tok := range tokens {
			fmt.Printf("%v\n", tok)	
		}

		parser := parser.NewParser(tokens)
		root, errors := parser.Parse()

		if len(errors) != 0 {
			for _, err := range errors {
				fmt.Printf("Error: %s\n", err)
			}
			continue
		}

		fmt.Printf("%v\n", root)
	}

}
