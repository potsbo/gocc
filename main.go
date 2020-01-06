package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/potsbo/gocc/node"
	"github.com/potsbo/gocc/token"
	"github.com/srvc/fail"
)

func main() {
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

func compile() error {
	if len(os.Args) != 2 {
		errors.New("Wrong size of arguments")
	}

	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".global _main")
	fmt.Println("_main:")

	proc, err := token.Tokenize(os.Args[1])
	if err != nil {
		return fail.Wrap(err)
	}
	p := node.NewParser(proc)
	prog, err := p.Generate()
	if err != nil {
		return fail.Wrap(err)
	}
	fmt.Println(prog)

	fmt.Println("  pop rax")
	fmt.Println("  ret")

	return nil
}
