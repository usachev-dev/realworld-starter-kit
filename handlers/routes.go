package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
)

func UseRoutes(r *mux.Router) {
	r.HandleFunc("/users", createUserHandle).Methods(http.MethodPost)
	r.HandleFunc("/users/login", signInHandle).Methods(http.MethodPost)
	r.HandleFunc("/profiles/{username}", getProfileHandle).Methods(http.MethodGet)

	authRoutes := r.NewRoute().Subrouter()
	authRoutes.Use(AuthRequest)
	authRoutes.HandleFunc("/user", getUserHandle).Methods(http.MethodGet)
	authRoutes.HandleFunc("/user", updateUserHandle).Methods(http.MethodPut)
}
