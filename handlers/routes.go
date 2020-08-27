package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
	"../auth"
)

func UseRoutes(r *mux.Router) {
	r.HandleFunc("/users", createUserHandle).Methods(http.MethodPost)
	r.HandleFunc("/users/login", signInHandle).Methods(http.MethodPost)

	authRoutes := r.NewRoute().Subrouter()
	authRoutes.Use(auth.AuthRequest)
	authRoutes.HandleFunc("/user", getUserHandle).Methods(http.MethodGet)
}
