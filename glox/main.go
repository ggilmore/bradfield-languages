package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ggilmore/bradfield-languages/glox/errutil"
	"github.com/ggilmore/bradfield-languages/glox/interpreter"
	"github.com/ggilmore/bradfield-languages/glox/parser"
	"github.com/ggilmore/bradfield-languages/glox/scanner"
)

const (
	ExUsage   = 64
	ExLox     = 65
	ExRuntime = 70
)

func main() {
	if len(os.Args) > 2 {
		fmt.Fprintln(os.Stderr, "Usage: glox [script]")
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

	runner := newRunner()
	err = runner.Run(f)
	if err != nil {
		printError(fmt.Errorf("running %q: %w", path, err))
		die(err)
	}
}

func runPrompt(r io.Reader) {
	runner := newRunner()
	s := bufio.NewScanner(r)

	prompt := "> "
	fmt.Print(prompt)

	for s.Scan() {
		line := s.Text()

		err := runner.Run(strings.NewReader(line))
		if err != nil {
			printError(err)

			var runErr interpreter.Error
			if !errors.Is(err, &runErr) {
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

type runner struct {
	interpreter *interpreter.Interpreter
}

func newRunner() *runner {
	return &runner{
		interpreter: interpreter.New(),
	}
}

func (r *runner) Run(input io.Reader) error {
	s, err := scanner.New(input)
	if err != nil {
		return fmt.Errorf("intializing scanner: %w", err)
	}

	tokens, err := s.Scan()
	if err != nil {
		return fmt.Errorf("scanning for tokens: %w", err)
	}

	statements, err := parser.NewParser(tokens).Parse()
	if err != nil {
		return fmt.Errorf("while parsing: %w", err)
	}

	resolver := interpreter.NewResolver(r.interpreter)
	err = resolver.Resolve(statements)
	if err != nil {
		return fmt.Errorf("while resolving: %w", err)
	}

	err = r.interpreter.Interpret(statements)
	if err != nil {
		return fmt.Errorf("while interpreting: %w", err)
	}

	return nil
}

func die(e error) {
	var runErr interpreter.Error
	if errors.Is(e, &runErr) {
		os.Exit(ExRuntime)
	}

	var loxErr errutil.LoxLanguageError
	if errors.As(e, &loxErr) {
		os.Exit(ExLox)
	}

	os.Exit(1)
}

func printError(err error) {
	var e errutil.LoxLanguageError
	if errors.As(err, &e) {
		// only print the underlying error if it's a lox
		// error so that don't clutter the output with
		// needless context
		fmt.Fprintln(os.Stderr, e.Error())
		return
	}

	fmt.Fprint(os.Stderr, err.Error())

}
