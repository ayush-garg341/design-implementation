package main

import (
	"fmt"
	"math"
)

func main() {
	num := aToI("1337c0d3")
	fmt.Println(num)

	num = aToI("   -042")
	fmt.Println(num)

	num = aToI("0-1")
	fmt.Println(num)

	num = aToI("words and 987")
	fmt.Println(num)

	num = aToI("9223372036854775808")
	fmt.Println(num)
}

func aToI(s string) int {
	symbol := '+'
	found := false

	var validChars []rune
	for _, c := range s {
		if isAlphabet(c) {
			break
		} else if !isDigit(c) {
			if found || len(validChars) != 0 {
				break
			}
			if c == ' ' {
				continue
			}
			if c == '+' || c == '-' {
				symbol = c
				found = true
			}

		} else {
			validChars = append(validChars, c)
		}
	}

	l := len(validChars)
	baseNum := 1
	if symbol == '-' {
		baseNum = -1
	}

	num := 0
	for i := 0; i < l; i++ {
		digit := int(validChars[i] - '0')
		num = num*10 + digit

		num, restricted := restrictTo32bit(num * baseNum)
		if restricted {
			return num
		}
	}

	return num * baseNum
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func isAlphabet(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '.'
}

func restrictTo32bit(num int) (int, bool) {
	if num >= int(math.Pow(2, 31)) {
		num = int(math.Pow(2, 31)) - 1
		return num, true
	}
	if num < int(math.Pow(-2, 31)) {
		num = int(math.Pow(-2, 31))
		return num, true
	}
	return num, false
}
