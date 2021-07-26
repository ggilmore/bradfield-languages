package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	ExUsage = 64
	ExLox   = 65
)

var hadError bool

func main() {
	if len(os.Args) > 2 {
		fmt.Fprintln(os.Stderr, "Usage: jlox [script]")
		os.Exit(ExUsage)
	}

	var err error

	if len(os.Args) == 2 {
		file := os.Args[1]
		err = runFile(file)
	} else {
		err = runPrompt(os.Stdin)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		if _, ok := err.(*loxError); ok {
			os.Exit(ExLox)
		}

		os.Exit(1)
	}
}

func runFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("while opening source file %q: %w", path, err)
	}

	err = run(f)
	if err != nil {
		return fmt.Errorf("while running source file %q: %w", path, err)
	}

	return nil
}

func runPrompt(r io.Reader) error {
	s := bufio.NewScanner(r)

	fmt.Print("> ")
	for s.Scan() {
		line := s.Text()

		err := run(strings.NewReader(line))
		if err != nil {
			return fmt.Errorf("while processing source line %s", err)
		}

		fmt.Printf("\n> ")
	}

	if err := s.Err(); err != nil {
		return fmt.Errorf("while processing input: %s", err)
	}

	return nil
}

func run(r io.Reader) error {
	s, err := NewScanner(r)
	if err != nil {
		return fmt.Errorf("while intializing scanner: %w", err)
	}

	tokens, err := s.ScanTokens()
	if err != nil {
		return fmt.Errorf("while scanning for tokens: %w", err)
	}

	p := NewParser(tokens)
	expr, err := p.Parse()
	if err != nil {
		return fmt.Errorf("while parsing: %w", err)
	}

	fmt.Printf("%s\n", expr)

	return nil
}

type loxError struct {
	line    int
	message string
}

func (e *loxError) Error() string {
	return fmt.Sprintf("[line %d] Error: %s", e.line, e.message)
}

func NewLoxError(line int, message string) error {
	hadError = true

	return &loxError{
		line,
		message,
	}
}
