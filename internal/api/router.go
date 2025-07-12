package api

import (
	"fmt"
	"net/http"
	"packhaus/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(router chi.Router) {
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{\"ok\": true}"))
	})

	router.Route("/api", func(r chi.Router) {
		r.Use(middleware.AuthMiddlware)

		r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
			val := r.Context().Value(middleware.ContextKeyUserID)
			userid, ok := val.(string)
			if !ok || userid == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			fmt.Printf("User ID: %s\n", userid)
			w.Write([]byte("hello user: " + userid))

		})
	})
}
