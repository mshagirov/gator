package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mshagirov/gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("Error: missing required argument time_between_reqs! E.g. \"10s\", \"1h\", ...\nUsage:\n  agg 15m")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) error {
	f, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting next feed to fetch: %w", err)
	}

	markParams := database.MarkFeedFetchedParams{
		ID:        f.ID,
		UpdatedAt: time.Now(),
	}
	if err := s.db.MarkFeedFetched(context.Background(), markParams); err != nil {
		return fmt.Errorf("Couldn't mark feed as fetched: %w", err)
	}
	rssData, err := fetchFeed(context.Background(), f.Url)
	if err != nil {
		return fmt.Errorf("Couldn't fetch feed: %w", err)
	}
	fmt.Printf("Fetching: %v\n---\nLink: %v\nDescription: %v\n---\n",
		rssData.Channel.Title,
		rssData.Channel.Link,
		rssData.Channel.Description,
	)
	for _, fi := range rssData.Channel.Item {
		fmt.Printf("%v\n", fi.Title)
	}
	fmt.Println("")
	return nil
}
