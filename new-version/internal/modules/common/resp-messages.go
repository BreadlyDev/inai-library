package common

import "fmt"

func ReqFieldMessage(field string) string {
	return fmt.Sprintf("%s: must not be empty", field)
}
