package user

import (
	"fmt"
	"net/mail"
	"strings"
)

func ValidateEmailFormat(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func ValidatePasswordLength(pass string, minLen int) bool {
	return len(pass) < minLen
}

func ValidatePasswordHasNumber(pass string) bool {
	return strings.ContainsAny(pass, "1234567890")
}

func ValidatePasswordHasCapitalizedLetter(pass string) bool {
	return strings.Contains(pass, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
}

func ValidatePasswordHasSpecialSymbol(pass string) bool {
	return strings.Contains(pass, "@#$%&/?.,-_+=)(^;:!")
}

func FullPasswordValidation(pass string) string {
	if !ValidatePasswordLength(pass, 10) {
		return fmt.Sprintf("%s", PasswordTooShort(pass))
	}

	if !ValidatePasswordHasNumber(pass) {
		return fmt.Sprintf("%s", PasswordHasNoNumber(pass))
	}

	if !ValidatePasswordHasCapitalizedLetter(pass) {
		return fmt.Sprintf("%s", PasswordHasNoCapitalizedLetter(pass))
	}

	if !ValidatePasswordHasSpecialSymbol(pass) {
		return fmt.Sprintf("%s", PasswordHasNoSpecialSymbol(pass))
	}

	return ""
}
