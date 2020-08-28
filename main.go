package main

import (
	"./DB"
	"./handlers"
	"fmt"
	"os"

	"./auth"
	"./models"
	"./utils"
	"log"
	"time"

	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	dbErr := DB.InitPostgres(utils.DbHost(), utils.DbPort(), utils.DbUser(), utils.DbName(), utils.DbPassword())
	if dbErr != nil {
		panic(fmt.Sprintf("could not connect to db: %s", dbErr))
	}
	defer DB.Close()
	router := mux.NewRouter()
	handlers.UseRoutes(router)
	models.AutoMigrate()
	SetSignature()
	p := utils.Port()
	srv := &http.Server{
		Handler:      router,
		Addr:         "localhost" + p,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func SetSignature() {
	s := os.Getenv("SIGNATURE")
	if s != "" {
		auth.SetSignature(s)
	}
}
