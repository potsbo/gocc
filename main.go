package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/potsbo/gocc/node"
	"github.com/potsbo/gocc/token"
)

func main() {
	err := compile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
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
		return err
	}
	p := node.NewParser(proc)
	prog, err := p.Generate()
	if err != nil {
		return err
	}
	fmt.Println(prog)

	fmt.Println("  pop rax")
	fmt.Println("  ret")

	return nil
}
