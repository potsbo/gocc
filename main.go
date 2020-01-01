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

	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".global _main")
	fmt.Println("_main:")

	str, i, err := strtoint(os.Args[1])
	if err != nil {
		return err
	}

	fmt.Printf("  mov rax, %d\n", i)
	for {
		if len(str) == 0 {
			break
		}
		if str[0] == '+' {
			str, i, err = strtoint(str[1:])
			if err != nil {
				return err
			}
			fmt.Printf("  add rax, %d\n", i)
			continue
		}

		if str[0] == '-' {
			str, i, err = strtoint(str[1:])
			if err != nil {
				return err
			}
			fmt.Printf("  sub rax, %d\n", i)
			continue
		}

		return errors.New("Unexpected char")
	}
	fmt.Println("  ret")

	return nil
}

func strtoint(str string) (string, int, error) {
	cnt := 0
	for _, char := range str {
		if _, err := strconv.Atoi(string(char)); err == nil {
			cnt++
		} else {
			break
		}
	}
	i, err := strconv.Atoi(str[0:cnt])
	return str[cnt:], i, err
}
