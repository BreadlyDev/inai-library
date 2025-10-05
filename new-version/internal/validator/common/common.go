package common

import "fmt"

func FieldIsRequired(field string) string {
	return fmt.Sprintf("%s is required", field)
}

func IsFieldNotEmpty(field string) bool {
	return !(field == "")
}
