package main

import (
	"errors"
	"fmt"
	"os"

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

	{ // must start with Num
		if proc.Finished() {
			return errors.New("No code given")
		}
		i, err := proc.ExtractNum()
		if err != nil {
			return err
		}

		fmt.Printf("  mov rax, %d\n", i)
	}

	for {
		if proc.Finished() {
			break
		}

		if proc.Consume("+") {
			i, err := proc.ExtractNum()
			if err != nil {
				return err
			}
			fmt.Printf("  add rax, %d\n", i)
			continue
		}
		if proc.Consume("-") {
			i, err := proc.ExtractNum()
			if err != nil {
				return err
			}
			fmt.Printf("  sub rax, %d\n", i)
			continue
		}
	}
	fmt.Println("  ret")

	return nil
}
