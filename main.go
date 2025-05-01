package main

import (
	"fmt"
	"interpreter/interpret"
	"interpreter/lexer"
	"interpreter/parser"
	"interpreter/types"
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
		
		// fmt.Printf("Tokens: \n")
		// for _, tok := range tokens {
		// 	fmt.Printf("%v\n", tok)
		// }

		parser := parser.NewParser(tokens)
		root, errors := parser.Parse()

		if len(errors) != 0 {
			for _, err := range errors {
				fmt.Printf("Error: %s\n", err)
			}
			continue
		}

		// fmt.Printf("%v\n", root)

		typechecker := types.Checker{
			Errors: make([]error, 0),
		}
		ok := typechecker.Visit(root)
		if !ok {
			fmt.Println("Got typerrors: ")
			for _, err := range typechecker.Errors {
				fmt.Printf("%v\n", err)
			}
			continue
		}

		interpreter := interpret.Interpreter{}
		interpreter.Visit(root)
	}

}
