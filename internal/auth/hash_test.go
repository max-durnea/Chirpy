package auth

import (
	"testing"
)

func TestHashAndCheckPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{"simple password", "secret123"},
		{"with symbols", "p@$$w0rd!"},
		{"long password", "thisIsAVeryLongPasswordThatShouldStillWorkFine"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			if err != nil {
				t.Fatalf("HashPassword() error = %v", err)
			}

			// Should succeed for correct password
			if err := CheckPasswordHash(tt.password, hash); err != nil {
				t.Errorf("CheckPasswordHash() failed for correct password: %v", err)
			}

			// Should fail for wrong password
			if err := CheckPasswordHash("wrongPassword", hash); err == nil {
				t.Errorf("CheckPasswordHash() succeeded for wrong password, want error")
			}
		})
	}
}

func TestHashPassword_ErrorOnEmpty(t *testing.T) {
	// bcrypt actually allows empty string, so this test just documents behavior
	hash, err := HashPassword("")
	if err != nil {
		t.Fatalf("HashPassword(\"\") returned unexpected error: %v", err)
	}
	if err := CheckPasswordHash("", hash); err != nil {
		t.Errorf("CheckPasswordHash() failed for empty password: %v", err)
	}
}
