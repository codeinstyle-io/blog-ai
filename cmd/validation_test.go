package cmd

import (
	"strings"
	"testing"
)

func TestValidateFirstName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{"Valid name", "John", false, ""},
		{"Valid name with hyphen", "Jean-Pierre", false, ""},
		{"Valid name with quote", "O'Connor", false, ""},
		{"Empty name", "", true, "first name must be between 1 and 255 characters"},
		{"Too long name", string(make([]byte, 256)), true, "first name must be between 1 and 255 characters"},
		{"Name with numbers", "John123", true, "first name can only contain letters, hyphens (-), and simple quotes (')"},
		{"Name with special chars", "John!", true, "first name can only contain letters, hyphens (-), and simple quotes (')"},
		{"Name with spaces", "John Doe", true, "first name can only contain letters, hyphens (-), and simple quotes (')"},
		{"Unicode letters", "José", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFirstName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFirstName() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && err.Error() != tt.errMsg {
				t.Errorf("ValidateFirstName() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestValidateLastName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{"Valid name", "Smith", false, ""},
		{"Valid name with hyphen", "Smith-Jones", false, ""},
		{"Valid name with quote", "O'Brien", false, ""},
		{"Empty name", "", true, "last name must be between 1 and 255 characters"},
		{"Too long name", string(make([]byte, 256)), true, "last name must be between 1 and 255 characters"},
		{"Name with numbers", "Smith123", true, "last name can only contain letters, hyphens (-), and simple quotes (')"},
		{"Name with special chars", "Smith!", true, "last name can only contain letters, hyphens (-), and simple quotes (')"},
		{"Name with spaces", "van Smith", true, "last name can only contain letters, hyphens (-), and simple quotes (')"},
		{"Unicode letters", "González", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLastName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateLastName() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && err.Error() != tt.errMsg {
				t.Errorf("ValidateLastName() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{"Valid email", "test@example.com", false, ""},
		{"Invalid email", "test", true, "invalid email format"},
		{"Invalid domain", "test@", true, "invalid email format"},
		{"Missing @", "test.com", true, "invalid email format"},
		{"Multiple @", "test@@example.com", true, "invalid email format"},
		{"Special chars in local", "test+123@example.com", false, ""},
		{"Subdomain", "test@sub.example.com", false, ""},
		{"IP address domain", "test@[127.0.0.1]", false, ""},
		{"Unicode in domain", "test@ünicode.com", true, "invalid email format"},
		{"Empty string", "", true, "invalid email format"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && err.Error() != tt.errMsg {
				t.Errorf("ValidateEmail() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{"Valid password", "Password123!", false, ""},
		{"Too short", "Pass1!", true, "password must be at least 8 characters long"},
		{"No uppercase", "password123!", true, "password must contain at least one uppercase letter"},
		{"No lowercase", "PASSWORD123!", true, "password must contain at least one lowercase letter"},
		{"No number", "Password!!!", true, "password must contain at least one number"},
		{"No special char", "Password123", true, "password must contain at least one special character"},
		{"Exactly 8 chars", "Pass123!", false, ""},
		{"With spaces", "Pass 123!", true, "password cannot contain spaces"},
		{"Unicode chars", "Pässword123!", false, ""},
		{"Empty string", "", true, "password must be at least 8 characters long"},
		{"Only special chars", "!@#$%^&*()", true, "password must contain at least one uppercase letter"},
		{"Maximum length", "Password123!" + strings.Repeat("a", 60), false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && err.Error() != tt.errMsg {
				t.Errorf("ValidatePassword() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}
