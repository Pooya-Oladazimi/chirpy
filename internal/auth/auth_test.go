package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.MustParse("1c022932-06c2-44dc-b453-367649307867")
	secret := "test-secret"

	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	got, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT returned error: %v", err)
	}

	if got != userID {
		t.Fatalf("ValidateJWT returned %s, want %s", got, userID)
	}
}

func TestValidateJWTRejectsWrongSecret(t *testing.T) {
	userID := uuid.MustParse("1c022932-06c2-44dc-b453-367649307867")

	token, err := MakeJWT(userID, "test-secret", time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	if _, err := ValidateJWT(token, "wrong-secret"); err == nil {
		t.Fatal("ValidateJWT returned nil error for wrong secret")
	}
}

func TestValidateJWTRejectsExpiredToken(t *testing.T) {
	userID := uuid.MustParse("1c022932-06c2-44dc-b453-367649307867")

	token, err := MakeJWT(userID, "test-secret", -time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	if _, err := ValidateJWT(token, "test-secret"); err == nil {
		t.Fatal("ValidateJWT returned nil error for expired token")
	}
}

func TestGetBearerToken(t *testing.T) {
	token, err := GetBearerToken(http.Header{
		"Authorization": []string{"Bearer test-token"},
	})
	if err != nil {
		t.Fatalf("GetBearerToken returned error: %v", err)
	}
	if token != "test-token" {
		t.Fatalf("GetBearerToken returned %q, want %q", token, "test-token")
	}
}

func TestGetBearerTokenRejectsMissingBearer(t *testing.T) {
	if _, err := GetBearerToken(http.Header{}); err == nil {
		t.Fatal("GetBearerToken returned nil error for missing header")
	}

	if _, err := GetBearerToken(http.Header{
		"Authorization": []string{"Basic test-token"},
	}); err == nil {
		t.Fatal("GetBearerToken returned nil error for non-bearer token")
	}
}
