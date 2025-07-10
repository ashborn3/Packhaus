package main

import (
	"net/http"
	"packhaus/internal/api"
	"packhaus/internal/config"
	"packhaus/internal/db"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		panic("error loading .env: " + err.Error())
	}

	pool, err := db.Connect()
	if err != nil {
		panic("error connecting to db: " + err.Error())
	}
	defer pool.Close()

	router := chi.NewRouter()

	router.Use(middleware.Logger)

	api.RegisterRoutes(router)

	http.ListenAndServe(":8888", router)
}
