package cmd

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"unicode"
)

func ValidateFirstName(name string) error {
	if len(name) < 1 || len(name) > 255 {
		return fmt.Errorf("first name must be between 1 and 255 characters")
	}

	for _, r := range name {
		if !unicode.IsLetter(r) {
			return fmt.Errorf("first name can only contain letters")
		}
	}
	return nil
}

func ValidateLastName(name string) error {
	if len(name) < 1 || len(name) > 255 {
		return fmt.Errorf("last name must be between 1 and 255 characters")
	}

	for _, r := range name {
		if !unicode.IsLetter(r) {
			return fmt.Errorf("last name can only contain letters")
		}
	}
	return nil
}

func ValidateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

	if email == "" {
		return fmt.Errorf("invalid email format")
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return fmt.Errorf("invalid email format")
	}

	// Allow IP addresses in domain part
	domain := parts[1]
	if strings.HasPrefix(domain, "[") && strings.HasSuffix(domain, "]") {
		ipStr := domain[1 : len(domain)-1]
		if net.ParseIP(ipStr) == nil {
			return fmt.Errorf("invalid email format")
		}
		return nil
	}

	// Regular email validation
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	if len(password) > 72 {
		return fmt.Errorf("password must be less than 72 characters")
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, r := range password {
		if unicode.IsSpace(r) {
			return fmt.Errorf("password cannot contain spaces")
		}
		if r > unicode.MaxASCII {
			return fmt.Errorf("password can only contain ASCII characters")
		}
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsNumber(r):
			hasNumber = true
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?", r):
			hasSpecial = true
		}
	}

	switch {
	case !hasUpper:
		return fmt.Errorf("password must contain at least one uppercase letter")
	case !hasLower:
		return fmt.Errorf("password must contain at least one lowercase letter")
	case !hasNumber:
		return fmt.Errorf("password must contain at least one number")
	case !hasSpecial:
		return fmt.Errorf("password must contain at least one special character")
	}
	return nil
}
