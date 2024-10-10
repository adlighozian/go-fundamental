package server

import (
	"go-axiata/https"
	"net/http"
)

func MapRouters(mux *http.ServeMux, handlers https.Handlers) {
	mux.HandleFunc("/api/posts", https.AuthMiddleware(handlers.GetPost(), ""))
	mux.HandleFunc("/api/posts/", https.AuthMiddleware(handlers.DetailPost(), ""))
	mux.HandleFunc("/register", handlers.Register())
	mux.HandleFunc("/login", handlers.Login())
	mux.HandleFunc("/logout", handlers.Logout())
}
