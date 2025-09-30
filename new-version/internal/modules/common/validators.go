package common

func ValidateFieldNotEmpty(field string) bool {
	if field == "" {
		return false
	}
	return true
}
