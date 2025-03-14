package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/onahvictor/rss-scrapper/internal/database"
)

func (app *application) createFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	var data struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	err := ReadJSON(r, &data)
	if err != nil {
		WriteJSON(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	feed := database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      data.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Url:       data.URL,
		UserID:    user.ID,
	}

	userFeed, err := app.DB.CreateFeed(r.Context(), feed)
	if err != nil {
		WriteJSON(w, "server error", http.StatusInternalServerError, nil)
		return
	}

	WriteJSON(w, NewFeed(userFeed), http.StatusCreated, nil)
}

func (app *application) GetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := app.DB.GetFeeds(r.Context())
	if err != nil {
		WriteJSON(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError, nil)
		return
	}

	responseFeed := map[string]any{
		"feeds": NewFeeds(feeds),
	}
	WriteJSON(w, responseFeed, http.StatusCreated, nil)
}
