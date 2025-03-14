package main

import (
	"fmt"
	"net/http"

	"github.com/onahvictor/rss-scrapper/internal/auth"
	"github.com/onahvictor/rss-scrapper/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (app *application) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			WriteJSON(w, err.Error(), http.StatusForbidden, nil)
			return
		}

		user, err := app.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			WriteJSON(w, fmt.Sprintf("couldn't get user: %v", err), http.StatusNotFound, nil)
			return
		}

		handler(w, r, user)
	}
}
