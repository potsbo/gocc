package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

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

	t, err := token.Tokenize(os.Args[1])
	if err != nil {
		return err
	}

	{ // must start with Num
		if t == nil {
			return errors.New("No code given")
		}
		if t.Kind != token.Num {
			return errors.New("Must start with numbers")
		}
		i, err := strconv.Atoi(t.Str)
		if err != nil {
			return err
		}
		fmt.Printf("  mov rax, %d\n", i)
		t = t.Next
	}

	for {
		if t.Kind == token.Eof {
			break
		}

		if t.Kind == token.Reserved {
			if t.Str == "+" {
				t = t.Next
				i, err := strconv.Atoi(t.Str)
				if err != nil {
					return err
				}
				fmt.Printf("  add rax, %d\n", i)
				t = t.Next
				continue
			}

			if t.Str == "-" {
				t = t.Next
				i, err := strconv.Atoi(t.Str)
				if err != nil {
					return err
				}
				fmt.Printf("  sub rax, %d\n", i)
				t = t.Next
				continue
			}
		}

		return errors.New(fmt.Sprintf("Unexpected char: %q", t.Str))
	}
	fmt.Println("  ret")

	return nil
}
