package handlers_test

import (
	"../handlers"
	"fmt"
	"net/http"
	"testing"
)

func TestGetTokenFromRequest(t *testing.T) {
	token := "sadsdsf"
	r, _ := http.NewRequest(http.MethodGet, "", nil)
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	parsedToken1, err1 := handlers.GetTokenFromRequest(r)
	if err1 != nil {
		t.Fatalf("could not get token from request header: %s", err1)
	}
	if parsedToken1 != token {
		t.Fatalf("token from header is not equal to provided")
	}

	r.Header.Add("Authorization", fmt.Sprintf("Token %s", token))
	parsedToken2, err2 := handlers.GetTokenFromRequest(r)
	if err2 != nil {
		t.Fatalf("could not get token from request header: %s", err2)
	}
	if parsedToken2 != token {
		t.Fatalf("token from header is not equal to provided")
	}

}
