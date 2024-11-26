package cli

import (
	"fmt"
	"regexp"
	"strings"
)

func validateFirstName(firstName string) error {
	if len(firstName) < 1 || len(firstName) > 255 {
		return fmt.Errorf("first name must be between 1 and 255 characters")
	}
	return nil
}

func validateLastName(lastName string) error {
	if len(lastName) < 1 || len(lastName) > 255 {
		return fmt.Errorf("last name must be between 1 and 255 characters")
	}
	return nil
}

func validateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if !emailRegex.MatchString(strings.ToLower(email)) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one number")
	}
	if !regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one special character")
	}
	return nil
}
