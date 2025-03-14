package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/onahvictor/rss-scrapper/internal/database"
)

func (app *application) createFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	var data struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	err := ReadJSON(r, &data)
	if err != nil {
		WriteJSON(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	feedFollow := database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    data.FeedID,
		UserID:    user.ID,
	}

	userFeedFollow, err := app.DB.CreateFeedFollows(r.Context(), feedFollow)
	if err != nil {
		WriteJSON(w, fmt.Sprintf("error: %v", err), http.StatusInternalServerError, nil)
		return
	}

	WriteJSON(w, NewFeedFollows(userFeedFollow), http.StatusCreated, nil)
}

func (app *application) userFeeds(w http.ResponseWriter, r *http.Request, user database.User) {
	userfeeds, err := app.DB.GetAllUserFeeds(r.Context(), user.ID)
	if err != nil {
		WriteJSON(w, "server error", http.StatusInternalServerError, nil)
		return
	}

	WriteJSON(w, map[string]any{"user_feeds": NewUserFeeds(userfeeds)}, http.StatusOK, nil)
}

func (app *application) DeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedfollowId, err := uuid.Parse(chi.URLParam(r, "feedFollowID"))
	if err != nil {
		WriteJSON(w, fmt.Sprintf("error parsing feed id%v", err), http.StatusOK, nil)
		return
	}

	err = app.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		FeedID:     feedfollowId,
		UserID: user.ID,
	})

	if err != nil {
		WriteJSON(w, fmt.Sprintf("error %v", err), http.StatusOK, nil)
		return
	}
	WriteJSON(w, "succesfully unfollowed a feed", http.StatusOK, nil)
}
