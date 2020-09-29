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
	dbErr := connectDb(5)
	if dbErr != nil {
		panic(fmt.Sprintf("could not connect to db: %s", dbErr))
	}
	defer DB.Close()
	router := mux.NewRouter()
	handlers.UseRoutes(router)
	models.AutoMigrate()
	SetSignature()
	port := utils.Port()
	host := utils.Host()
	srv := &http.Server{
		Handler:      router,
		Addr:         host + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func tryConnect() error {
	return DB.InitPostgres(utils.DbHost(), utils.DbPort(), utils.DbUser(), utils.DbName(), utils.DbPassword())
}

func connectDb(tries uint) error {
	var err error
	for i := 0; i < int(tries); i++ {
		err = tryConnect()
		if err == nil {
			return nil
		}
		time.Sleep(time.Second * 3)
	}
	return err
}

func SetSignature() {
	s := os.Getenv("SIGNATURE")
	if s != "" {
		auth.SetSignature(s)
	}
}
