package main

import (
	"fmt"
	"interpreter/lexer"
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
				fmt.Printf("Error: %s", err)
			}
		}

		for _, tok := range tokens {
			fmt.Printf("%v\n", tok)
		}

	}

}
