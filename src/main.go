package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	ExUsage   = 64
	ExLox     = 65
	ExRuntime = 70
)

func main() {
	if len(os.Args) > 2 {
		fmt.Fprintln(os.Stderr, "Usage: jlox [script]")
		os.Exit(ExUsage)
	}

	if len(os.Args) == 2 {
		file := os.Args[1]
		runFile(file)
	} else {
		runPrompt(os.Stdin)
	}
}

func runFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		printError(fmt.Errorf("opening %q: %w", path, err))
		die(err)
	}

	err = run(f)
	if err != nil {
		printError(fmt.Errorf("running %q: %w", path, err))
		die(err)
	}
}

func runPrompt(r io.Reader) {
	s := bufio.NewScanner(r)

	prompt := "> "
	fmt.Print(prompt)

	for s.Scan() {
		line := s.Text()

		err := run(strings.NewReader(line))
		if err != nil {
			printError(err)

			var runErr loxRuntimeError
			if !errors.As(err, &runErr) {
				die(err)
			}
		}

		fmt.Printf("\n%s", prompt)
	}

	if err := s.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "while processing input: %s\n", err)
		die(err)
	}
}

func run(r io.Reader) error {
	s, err := NewScanner(r)
	if err != nil {
		return fmt.Errorf("intializing scanner: %w", err)
	}

	tokens, err := s.ScanTokens()
	if err != nil {
		return fmt.Errorf("scanning for tokens: %w", err)
	}

	p := NewParser(tokens)
	expr, err := p.Parse()
	if err != nil {
		return fmt.Errorf("while parsing: %w", err)
	}

	expr, err = Evaluate(expr)
	if err != nil {
		return err
	}

	fmt.Println(expr)
	return nil
}

func die(e error) {
	var runErr loxRuntimeError
	if errors.As(e, &runErr) {
		os.Exit(ExRuntime)
	}

	var loxErr loxError
	if errors.As(e, &loxErr) {
		os.Exit(ExLox)
	}

	os.Exit(1)
}

func printError(err error) {
	var e LoxLanguageError
	if errors.As(err, &e) {
		// only print the underlying error if it's a lox
		// error so that don't clutter the output with
		// needless context
		fmt.Fprintln(os.Stderr, e.Error())
		return
	}

	fmt.Fprint(os.Stderr, err.Error())
}

type LoxLanguageError interface {
	IsLoxLanguageError()
	Error() string
}

type loxError struct {
	line    int
	message string
}

func (e loxError) Error() string {
	return fmt.Sprintf("[line %d] Error: %s", e.line, e.message)
}

func (e loxError) IsLoxLanguageError() {}
