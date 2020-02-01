package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/potsbo/gocc/node"
	"github.com/potsbo/gocc/token"
	"github.com/srvc/fail"
)

var debug bool

func main() {
	if os.Getenv("GOCC_DEBUG") == "true" {
		debug = true
	}

	err := compile()
	if err != nil {
		aerr := fail.Unwrap(err)
		fmt.Fprintf(os.Stderr, "%s\n", aerr.Error())
		for _, f := range aerr.StackTrace {
			fmt.Fprintf(os.Stderr, "%s in %s:L%d\n", f.Func, f.File, f.Line)
		}
		os.Exit(1)
	}
}

func inspectTokens(tokens []token.Token) {
	for _, token := range tokens {
		fmt.Fprintf(os.Stderr, "%v\n", token)
	}
}

func compile() error {
	if len(os.Args) != 2 {
		errors.New("Wrong size of arguments")
	}

	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".global _main")
	fmt.Println("_main:")

	fmt.Println("# prologue")
	fmt.Println("  push rbp")
	fmt.Println("  mov rbp, rsp")
	fmt.Println("  sub rsp, 208") // 26 * 8
	fmt.Println("# prologue end")

	{
		proc, err := token.Tokenize(os.Args[1])
		if err != nil {
			return fail.Wrap(err)
		}
		if debug {
			inspectTokens(proc.Inspect())
		}
		p := node.NewParser(proc)
		prog, err := p.Generate()
		if err != nil {
			return fail.Wrap(err)
		}
		fmt.Println(prog)
		fmt.Println("  pop rax")
	}

	fmt.Println("# epilogue")
	fmt.Println("  mov rsp, rbp")
	fmt.Println("  pop rbp")
	fmt.Println("  ret")
	fmt.Println("# epilogue end")

	return nil
}
