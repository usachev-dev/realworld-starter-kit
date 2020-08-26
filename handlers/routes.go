package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
)

func UseRoutes(r *mux.Router) {
	r.HandleFunc("/users", createUserHandle).Methods(http.MethodPost)
	r.HandleFunc("/users/login", signInHandle).Methods(http.MethodPost)
}
