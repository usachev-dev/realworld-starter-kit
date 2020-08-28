package auth_test

import (
	"../auth"
	"testing"
)

func TestGetTokenString(t *testing.T) {
	email := "saergdgfg"
	result := auth.GetTokenString(email)
	if result == "" {
		t.Fatalf("could not get token string")
	}
	if len(result) < len(email) {
		t.Fatalf("token is too short")
	}
}

func TestRandomTokenIsInvalid(t *testing.T) {
	result := auth.ValidateTokenStringWithEmail("alkjndasoljsewoldsglgndfsg", "sdasd")
	if result == nil {
		t.Fatalf("random string is not validated as token")
	}
}
