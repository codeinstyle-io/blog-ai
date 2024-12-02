package cmd

import "testing"

func TestValidateFirstName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{"Valid name", "John", false, ""},
		{"Empty name", "", true, "first name must be between 1 and 255 characters"},
		{"Too long name", string(make([]byte, 256)), true, "first name must be between 1 and 255 characters"},
		{"Name with numbers", "John123", true, "first name can only contain letters"},
		{"Name with special chars", "John!", true, "first name can only contain letters"},
		{"Name with spaces", "John Doe", true, "first name can only contain letters"},
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
		{"Empty name", "", true, "last name must be between 1 and 255 characters"},
		{"Too long name", string(make([]byte, 256)), true, "last name must be between 1 and 255 characters"},
		{"Name with numbers", "Smith123", true, "last name can only contain letters"},
		{"Name with special chars", "Smith!", true, "last name can only contain letters"},
		{"Name with spaces", "van Smith", true, "last name can only contain letters"},
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
		{"Valid password", "Test1234!", false, ""},
		{"Too short", "Test1!", true, "password must be at least 8 characters"},
		{"No uppercase", "test1234!", true, "password must contain at least one uppercase letter"},
		{"No lowercase", "TEST1234!", true, "password must contain at least one lowercase letter"},
		{"No number", "TestTest!", true, "password must contain at least one number"},
		{"No special char", "Test1234", true, "password must contain at least one special character"},
		{"Exactly 8 chars", "Test1!Px", false, ""},
		{"With spaces", "Test 1!Px", true, "password cannot contain spaces"},
		{"Unicode chars", "Test1!のパ", true, "password can only contain ASCII characters"},
		{"Empty string", "", true, "password must be at least 8 characters"},
		{"Only special chars", "!@#$%^&*()", true, "password must contain at least one uppercase letter"},
		{"Maximum length", string(make([]byte, 72)) + "T1!", true, "password must be less than 72 characters"},
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
