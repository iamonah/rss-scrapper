package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/onahvictor/rss-scrapper/internal/database"
)

// the scraper function is a long running job and this scrapper function runs as long as our server runs
// having a good logging is essential
func startScrapping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequest time.Duration,
) {
	log.Printf("Scrapping on %v goroutines every %s duration", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	defer ticker.Stop()

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedToFetch(context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Println("error fetching feeds:", err)
			continue
		}
		//where you spawned a new goroutine for every feed in every batch, potentially creating tons of them over time.
		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	feed, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("error marking feed as fetched", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("error fetching feeds:", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		pubAT, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("couldn't parse date %v with err %v", item.PubDate, err)
			continue
		}
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: description,
			PublishedAt: pubAT,
			Url:         item.Link,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Printf("failed to create post: %v", err)
		}
	}
	log.Printf("feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}



// Define a new handler type that returns an error
// type HTTPHandlerWithErr func(http.ResponseWriter, *http.Request) error

// // Handle wraps your error-returning handlers
// func (r *Router) Handle(pattern string, handler HTTPHandlerWithErr) {
//     r.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
//         if err := handler(w, r); err != nil {
//             // Check if it's an HTTPError
//             var httpErr *httperror.HTTPError
//             if errors.As(err, &httpErr) {
//                 http.Error(w, err.Error(), httpErr.Code)
//                 slog.Debug("http error", "code", httpErr.Code, "err", err.Error())
//             } else {
//                 // Default to 500
//                 http.Error(w, err.Error(), http.StatusInternalServerError)
//                 slog.Error("internal server error", "err", err.Error())
//             }
//         }
//     })
// }