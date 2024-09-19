package validate

import "regexp"

func ValidateEmail(errors *map[string]string, email string) {
	// Validate email format using regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		(*errors)["email"] = "invalid email format"
	}
}

// ValidateEmail validates the email format
func ValidatePassword(password string) map[string]string {
	errors := make(map[string]string)

	// Validate password length
	if len(password) < 8 {
		errors["password"] = "password must be at least 8 characters long"
		return errors
	}

	// Validate uppercase letter
	uppercaseRegex := regexp.MustCompile(`[A-Z]`)
	if !uppercaseRegex.MatchString(password) {
		errors["password"] = "password must contain at least one uppercase letter"
		return errors
	}

	// Validate lowercase letter
	lowercaseRegex := regexp.MustCompile(`[a-z]`)
	if !lowercaseRegex.MatchString(password) {
		errors["password"] = "password must contain at least one lowercase letter"
		return errors
	}

	// Validate digit
	digitRegex := regexp.MustCompile(`[0-9]`)
	if !digitRegex.MatchString(password) {
		errors["password"] = "password must contain at least one digit"
		return errors
	}

	// Validate special character
	specialCharRegex := regexp.MustCompile(`[!@#$%^&*]`)
	if !specialCharRegex.MatchString(password) {
		errors["password"] = "password must contain at least one special character"
		return errors
	}

	return errors
}
