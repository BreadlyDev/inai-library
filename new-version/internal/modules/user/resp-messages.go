package user

import "fmt"

func WrongEmailFormat(email string) string {
	return fmt.Sprintf("wrong email format: %s", email)
}

func PasswordTooShort(pass string) string {
	return fmt.Sprintf("password is too short: %s", pass)
}

func PasswordHasNoNumber(pass string) string {
	return fmt.Sprintf("password has no number: %s", pass)
}

func PasswordHasNoCapitalizedLetter(pass string) string {
	return fmt.Sprintf("password has no capitalized letter: %s", pass)
}

func PasswordHasNoSpecialSymbol(pass string) string {
	return fmt.Sprintf("password has no special symbol: %s", pass)
}
