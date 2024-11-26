package utils

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "mySecurePassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("HashPassword failed: %v", err)
	}
	if hash == password {
		t.Error("HashPassword should not return plain text password")
	}
	if len(hash) == 0 {
		t.Error("HashPassword should not return empty hash")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "mySecurePassword123"
	wrongPassword := "wrongPassword123"

	hash, _ := HashPassword(password)

	// Test correct password
	if !CheckPasswordHash(password, hash) {
		t.Error("CheckPasswordHash should return true for correct password")
	}

	// Test wrong password
	if CheckPasswordHash(wrongPassword, hash) {
		t.Error("CheckPasswordHash should return false for wrong password")
	}
}
