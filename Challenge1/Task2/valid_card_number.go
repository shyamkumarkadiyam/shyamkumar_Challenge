package main

import (
	"fmt"
	"regexp"
)

func validateCreditCardNumber(cardNumber string) string {
	pattern := `^[456]((\d{3})(-?\d{4}){3}|(\d{4})(-?\d{4}){2}(\d{4}))$`
	regex := regexp.MustCompile(pattern)

	if !regex.MatchString(cardNumber) {
		return "Invalid"
	}

	noHyphenCardNumber := removeHyphens(cardNumber)
	repeatedDigitsPattern := `(\d)\1{3}`
	if regexp.MustCompile(repeatedDigitsPattern).MatchString(noHyphenCardNumber) {
		return "Invalid"
	}

	return "Valid"
}

func removeHyphens(s string) string {
	return regexp.MustCompile(`-`).ReplaceAllString(s, "")
}

func main() {
	var n int
	fmt.Scan(&n)

	for i := 0; i < n; i++ {
		var cardNumber string
		fmt.Scan(&cardNumber)
		fmt.Println(validateCreditCardNumber(cardNumber))
	}
}
