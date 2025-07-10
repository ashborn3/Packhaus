package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(router chi.Router) {
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{\"ok\": true}"))
	})
}
