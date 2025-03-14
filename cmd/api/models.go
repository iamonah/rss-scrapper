package main

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/onahvictor/rss-scrapper/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

func NewUser(d database.User) *User {
	return &User{
		ID:        d.ID,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
		Name:      d.Name,
		ApiKey:    d.ApiKey,
	}
}

type Feed struct {
	ID            uuid.UUID    `json:"id"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
	Name          string       `json:"name"`
	Url           string       `json:"url"`
	UserID        uuid.UUID    `json:"user_id"`
	LastFetchedAt sql.NullTime `json:"lastfetched_at"`
}

func NewFeed(feed database.Feed) *Feed {
	return &Feed{
		ID:            feed.ID,
		CreatedAt:     feed.CreatedAt,
		UpdatedAt:     feed.UpdatedAt,
		Name:          feed.Name,
		Url:           feed.Url,
		UserID:        feed.UserID,
		LastFetchedAt: feed.LastFetchedAt,
	}
}

func NewFeeds(dbFeeds []database.Feed) []Feed {
	feeds := []Feed{}
	for _, dbfeed := range dbFeeds {
		feeds = append(feeds, Feed(dbfeed))
	}
	return feeds
}

type FeedsFollow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FeedID    uuid.UUID `json:"feed_id"`
	UserID    uuid.UUID `json:"user_id"`
}

func NewFeedFollows(feed database.FeedsFollow) *FeedsFollow {
	return &FeedsFollow{
		ID:        feed.ID,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		FeedID:    feed.FeedID,
		UserID:    feed.UserID,
	}
}

func NewUserFeeds(dbFeeds []database.FeedsFollow) []*FeedsFollow {
	feeds := []*FeedsFollow{}
	for _, dbfeed := range dbFeeds {
		feeds = append(feeds, NewFeedFollows(dbfeed))
	}
	return feeds
}

type Post struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	Url         string    `json:"url"`
	FeedID      uuid.UUID `json:"feed_id"`
}

func UserPost(dbpost database.Post) *Post {
	var description *string
	if dbpost.Description.Valid {
		description = &dbpost.Description.String
	}
	return &Post{
		ID:          dbpost.ID,
		CreatedAt:   dbpost.CreatedAt,
		UpdatedAt:   dbpost.UpdatedAt,
		Title:       dbpost.Title,
		Description: description,
		PublishedAt: dbpost.PublishedAt,
		Url:         dbpost.Url,
		FeedID:      dbpost.FeedID,
	}
}

func UserPosts(dbposts []database.Post)[]*Post{
	posts := []*Post{}
	for _, dbpost := range dbposts{
		posts = append(posts, UserPost(dbpost))
	}
	return posts
}