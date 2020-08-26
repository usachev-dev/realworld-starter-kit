package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"../api_errors"
	"time"
)

var signature = []byte("Change this with SIGNATURE env variable")

func expireDuration() time.Duration {
	return time.Hour * 100
}

func SetSignature() {
	s := os.Getenv("SIGNATURE")
	if s != "" {
		signature = []byte(s)
	}
}

type Claims struct {
	jwt.StandardClaims
	Email string `json:"email"`
}

func PasswordToHash(password string) string {
	phash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(phash)
}

func CheckPassword(password string, storedHash string) error {
	return bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
}

func GetToken(email string) *jwt.Token {
	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expireDuration()).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Email: email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	return token
}

func GetTokenString(email string) string {
	result, err := GetToken(email).SignedString(signature)
	if err != nil {
		panic(fmt.Sprintf("Could not get token string from token: %s", err))
	}
	return result
}

func ValidateTokenString(tokenString string, email string) *api_errors.E {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, KeyFunc)
	if err != nil {
		return api_errors.NewError(http.StatusUnauthorized).Add("Auth", fmt.Sprintf("could not authenticate token: %s", err.Error()))
	}
	claims, ok := token.Claims.(*Claims)
	if ok && token.Valid && claims.Email == email {
		return &api_errors.Ok
	} else {
		return api_errors.NewError(http.StatusUnauthorized).Add("Auth", fmt.Sprintf("could not authorize token: %s", claims.Valid().Error()))
	}
}

func GetEmailFromTokenString(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, KeyFunc)
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(*Claims)
	if ok && claims.Email != "" {
		return claims.Email, nil
	} else {
		return "", fmt.Errorf("could not get email from token")
	}
}

func KeyFunc(token *jwt.Token) (interface{}, error)  {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}
	return signature, nil
}