package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/onahvictor/rss-scrapper/internal/database"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {

}

func (app *application) registerUser(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Name string `json:"name"`
	}

	err := ReadJSON(r, &data)
	if err != nil {
		WriteJSON(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	dataValue := database.CreateUserParams{
		ID:        uuid.New(),
		Name:      data.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	user, err := app.DB.CreateUser(r.Context(), dataValue)
	if err != nil {
		WriteJSON(w, "server error", http.StatusInternalServerError, nil)
		return
	}

	WriteJSON(w, NewUser(user), http.StatusCreated, nil)
}

func (app *application) getUser(w http.ResponseWriter, r *http.Request, user database.User) {
	WriteJSON(w, NewUser(user), http.StatusOK, nil)
}

func (app *application) getPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := app.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		WriteJSON(w, "server error", http.StatusInternalServerError, nil)
		return
	}
	WriteJSON(w, map[string]any{"feeds":UserPosts(posts)}, http.StatusOK, nil)

}
