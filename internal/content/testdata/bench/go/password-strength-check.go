import "unicode"

func is_strong_password(s string) bool {
	if len(s) < 8 {
		return false
	}
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false
	for _, c := range s {
		if unicode.IsUpper(c) {
			hasUpper = true
		} else if unicode.IsLower(c) {
			hasLower = true
		} else if unicode.IsDigit(c) {
			hasDigit = true
		} else {
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasDigit && hasSpecial
}
