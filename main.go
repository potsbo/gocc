package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

func main() {
	err := compile()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func compile() error {
	if len(os.Args) != 2 {
		errors.New("Wrong size of arguments")
	}

	i, err := strconv.Atoi(os.Args[1])
	if err != nil {
		return err
	}

	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".global _main")
	fmt.Println("_main:")
	fmt.Printf("  mov rax, %d\n", i)
	fmt.Println("  ret")

	return nil
}
