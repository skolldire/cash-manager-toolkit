package validation

func ValidateStatusCode(status int, codes []int) bool {
	isValidStatusCode := false
	for _, code := range codes {
		if status == code {
			isValidStatusCode = true
			break
		}
	}
	return isValidStatusCode
}
