package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	secret := "supersecretkey"
	userID := uuid.New()

	token, err := MakeJWT(userID, secret, time.Minute)
	if err != nil {
		t.Fatalf("failed to make JWT: %v", err)
	}
	gotID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("failed to validate JWT: %v", err)
	}

	if gotID != userID {
		t.Errorf("expected userID %v, got %v", userID, gotID)
	}
}
