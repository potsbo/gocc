package util

import "strconv"

func Strtoint(str string) (string, int, error) {
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

func IsDigit(c rune) bool {
	digits := []rune("0123456789")

	for _, d := range digits {
		if c == d {
			return true
		}
	}

	return false
}
