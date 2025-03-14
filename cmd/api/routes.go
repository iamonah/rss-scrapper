package main

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *application) Mount() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	mux.Use(middleware.Logger)

	v1mux := chi.NewRouter()
	v1mux.Get("/healthz", handlerReadiness)
	v1mux.Post("/users", app.registerUser)
	v1mux.Get("/feeds", app.GetFeeds)
	v1mux.Get("/users", app.middlewareAuth(app.getUser))
	v1mux.Post("/feeds", app.middlewareAuth(app.createFeed))

	v1mux.Post("/feed-follows", app.middlewareAuth(app.createFeedFollow))
	v1mux.Get("/feeds-follows", app.middlewareAuth(app.userFeeds))
	v1mux.Delete("/feed-follows/{feedFollowID}", app.middlewareAuth(app.DeleteFeedFollow))
	v1mux.Get("/userfeed", app.middlewareAuth(app.getPostsForUser))

	mux.Mount("/v1/api", v1mux)

	return mux
}
