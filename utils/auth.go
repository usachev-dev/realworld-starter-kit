package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

var signature = "Change this with SIGNATURE env variable"

func expireDuration() time.Duration {
	return time.Hour * 100
}

func SetSignature() {
	s := os.Getenv("SIGNATURE")
	if s != "" {
		signature = s
	}
}

type Claims struct {
	jwt.StandardClaims
	Username string `json:"username"`
}

func PasswordToHash(password string) string {
	phash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(phash)
}

func GetToken(username string) *jwt.Token {
	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expireDuration()).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Username: username,
	}
	fmt.Printf("Claims are %+v\n", claims)
	fmt.Println()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	fmt.Printf("Token is %+v\n", token)
	return token
}

func GetTokenString(username string) string {
	result, err := GetToken(username).SignedString(signature)
	if err != nil {
		panic(fmt.Sprintf("Could not get token string from token: %s", err))
	}
	fmt.Printf("Token string is %s\n", result)
	return result
}
