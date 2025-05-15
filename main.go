package main

import (
	"fmt"
	"interpreter/interpret"
	"interpreter/lexer"
	"interpreter/parser"
	"interpreter/types"
	"io"
	"os"

	"github.com/chzyer/readline"
)

func main() {
	if len(os.Args) > 1 {
		interpretProgram(os.Args[1])
	}

	repl()
}

func interpretProgram(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}

	runProgram(content, path)
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

		runProgram([]byte(line), "repl")
	}

}

func runProgram(program []byte, file string) {
	lexer := lexer.NewLexer(program, file)
	tokens, errors := lexer.Tokenize()
	if errors != nil {
		for _, err := range errors {
			fmt.Printf("%s\n", err)
		}
	}

	// fmt.Printf("Tokens: \n")
	// for _, tok := range tokens {
	// 	fmt.Printf("%v\n", tok)
	// }

	parser := parser.NewParser(tokens, file)
	root, errors := parser.Parse()

	if len(errors) != 0 {
		for _, err := range errors {
			fmt.Printf("%s\n", err)
		}
		return
	}

	// fmt.Printf("%v\n", root)

	typechecker := types.NewChecker(file)
	ok := typechecker.Visit(root)
	if !ok {
		fmt.Println("Got typerrors: ")
		for _, err := range typechecker.Errors {
			fmt.Printf("%v\n", err)
		}
		return
	}

	interpreter := interpret.NewInterpreter()
	interpreter.Visit(root)
}
