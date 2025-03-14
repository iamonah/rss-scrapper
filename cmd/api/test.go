package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/onahvictor/rss-scrapper/internal/database"
)

func Scrapping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequest time.Duration,
	ctx context.Context,
) {
	ticker := time.NewTicker(timeBetweenRequest)
	defer ticker.Stop()

	jobs := make(chan Feed, concurrency)
	var wg sync.WaitGroup

	for range concurrency {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for feed := range jobs {
				scrapeFeed(db, nil, database.Feed(feed))
			}
		}()
	}

	for {
		select {
		case <-ticker.C:
			feeds, err := db.GetNextFeedToFetch(context.Background(), int32(concurrency))
			if err != nil {
				log.Printf("error fetching feeds: %v", err)
				continue
			}
			for _, feed := range feeds {
				jobs <- Feed(feed)
			}
		case <-ctx.Done():
			close(jobs)
			wg.Wait()
			return

		}
	}
}
