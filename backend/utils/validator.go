package utils

func IsValidCurrency(currency string) bool {
	switch currency {
	case "USD", "EUR", "INR":
		return true
	}
	return false
}
