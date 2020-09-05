package utils

import (
	"os"
)

func Port() string {
	p := os.Getenv("PORT")
	if p == "" {
		p = "4000"
	}
	p = ":" + p
	return p
}

func Host() string {
	p := os.Getenv("HOST")
	if p == "" {
		p = "0.0.0.0"
	}
	p = ":" + p
	return p
}

func DbPort() string {
	p := os.Getenv("DB_PORT")
	if p == "" {
		p = "5432"
	}
	return p
}

func DbHost() string {
	p := os.Getenv("DB_HOST")
	if p == "" {
		p = "localhost"
	}
	return p
}

func DbUser() string {
	p := os.Getenv("DB_USER")
	if p == "" {
		p = "postgres"
	}
	return p
}

func DbName() string {
	p := os.Getenv("DB_NAME")
	if p == "" {
		p = "conduit"
	}
	return p
}

func DbPassword() string {
	p := os.Getenv("DB_PASSWORD")
	return p
}
