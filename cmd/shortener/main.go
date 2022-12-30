package main

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/vladimirimekov/url-shortener"
	"github.com/vladimirimekov/url-shortener/internal/handlers"
	"github.com/vladimirimekov/url-shortener/internal/middlewares"
	"github.com/vladimirimekov/url-shortener/internal/storage"
	"log"
	"net/http"
)

const userKey string = "userid"

func main() {

	cfg := urlshortener.GetConfig()

	s := storage.Storage{Filename: cfg.Filename}
	h := handlers.Handler{Storage: s, LengthOfShortname: cfg.ShortnameLength, Host: cfg.BaseURL, UserKey: userKey}
	m := middlewares.UserCookies{Storage: s, Secret: cfg.Secret, UserKey: userKey}

	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	r.Use(middlewares.GZIPRead)
	r.Use(middlewares.GZIPWrite)
	r.Use(m.CheckUserCookies)

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.MainHandler)
	})

	r.Route("/api/user/urls", func(r chi.Router) {
		r.Get("/", h.AllShorterURLsHandler)
	})
	r.Post("/", h.MainHandler)

	r.Route("/api/shorten", func(r chi.Router) {
		r.Post("/", h.ShortenHandler)
	})

	log.Fatal(http.ListenAndServe(cfg.ServerAddress, r))
}
