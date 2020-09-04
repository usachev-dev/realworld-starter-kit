package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
)

func UseRoutes(r *mux.Router) {
	authRoutes := r.NewRoute().Subrouter()
	authRoutes.Use(AuthRequest)
	authRoutes.HandleFunc("/user", getUserHandle).Methods(http.MethodGet)
	authRoutes.HandleFunc("/user", updateUserHandle).Methods(http.MethodPut)
	authRoutes.HandleFunc("/profiles/{username}/follow", followHandle).Methods(http.MethodPost)
	authRoutes.HandleFunc("/profiles/{username}/follow", unfollowHandle).Methods(http.MethodDelete)
	authRoutes.HandleFunc("/articles", createArticleHandle).Methods(http.MethodPost)
	authRoutes.HandleFunc("/articles/feed", feedArticlesHandle).Methods(http.MethodGet)
	authRoutes.HandleFunc("/articles/{slug}", deleteArticleHandle).Methods(http.MethodDelete)
	authRoutes.HandleFunc("/articles/{slug}", updateArticleHandle).Methods(http.MethodPut)
	authRoutes.HandleFunc("/articles/{slug}/favorite", favoriteArticleHandle).Methods(http.MethodPost)
	authRoutes.HandleFunc("/articles/{slug}/favorite", unfavoriteArticleHandle).Methods(http.MethodDelete)

	r.HandleFunc("/users", createUserHandle).Methods(http.MethodPost)
	r.HandleFunc("/users/login", signInHandle).Methods(http.MethodPost)
	r.HandleFunc("/profiles/{username}", getProfileHandle).Methods(http.MethodGet)
	r.HandleFunc("/articles/{slug}", getArticleHandle).Methods(http.MethodGet)
	r.HandleFunc("/articles", listArticlesHandle).Methods(http.MethodGet)
}
