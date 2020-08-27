package auth_test

import (
	"fmt"
	"net/http"
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
	result := auth.ValidateTokenStringWithEmail("alkjndasoljsewoldsglgndfsg", "sdasd")
	if result == nil {
		t.Fatalf("random string is not validated as token")
	}
}

func TestGetTokenFromRequest(t *testing.T) {
	token := "sadsdsf"
	r, _ := http.NewRequest(http.MethodGet, "", nil)
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	parsedToken, err := auth.GetTokenFromRequest(r)
	if err != nil {
		t.Fatalf("could not get token from request header: %s", err)
	}
	if parsedToken != token {
		t.Fatalf("token from header is not equal to provided")
	}
}
