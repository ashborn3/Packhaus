package api

import (
	"net/http"
	"packhaus/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type controller struct {
	DB *pgxpool.Pool
}

func newController(db *pgxpool.Pool) *controller {
	return &controller{
		DB: db,
	}
}

func RegisterRoutes(router chi.Router, pool *pgxpool.Pool) {
	cntlr := newController(pool)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{\"ok\": true}"))
	})

	router.Post("/auth/signup", cntlr.SignupHandler)
	router.Post("/auth/login", cntlr.SigninHandler)

	router.Route("/api", func(r chi.Router) {
		r.Use(middleware.AuthMiddlware)
		r.Get("/me", cntlr.MeHandler)
		r.Post("/packages", cntlr.UploadPackageHandler)
	})

}
