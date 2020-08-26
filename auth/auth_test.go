package auth_test

import (
	"testing"
	"../auth"
)

func TestGetTokenString(t *testing.T) {
	auth.SetSignature()
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
	result := auth.ValidateTokenString("alkjndasoljsewoldsglgndfsg", "sdasd")
	if result == nil {
		t.Fatalf("random string is not validated as token")
	}
}
