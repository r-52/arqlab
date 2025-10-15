package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

const version = "0.1.0-pre"

func main() {
	modeRepl := flag.Bool("repl", false, "start an interactive REPL session")
	filePath := flag.String("file", "", "path to a JavaScript file to execute")
	showVersion := flag.Bool("version", false, "print the interpreter version")

	flag.Parse()

	if *showVersion {
		fmt.Println("es6-interpreter", version)
		return
	}

	switch {
	case *modeRepl:
		if err := startREPL(); err != nil {
			exitWithError(err)
		}
	case *filePath != "":
		if err := runFile(*filePath); err != nil {
			exitWithError(err)
		}
	default:
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintln(os.Stderr, "  es6-interpreter -repl")
		fmt.Fprintln(os.Stderr, "  es6-interpreter -file program.js")
		fmt.Fprintln(os.Stderr, "  es6-interpreter -version")
		os.Exit(2)
	}
}

func startREPL() error {
	// TODO: Implement full REPL once lexer, parser, and VM are ready.
	return errors.New("REPL is not implemented yet")
}

func runFile(path string) error {
	source, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read script: %w", err)
	}

	// TODO: Thread source through lexer -> parser -> VM pipeline.
	_ = source

	return errors.New("file execution is not implemented yet")
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}
