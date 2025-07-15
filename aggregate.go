package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mshagirov/gator/internal/database"
)

func scrapeFeeds(s *state) error {
	f, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	markParams := database.MarkFeedFetchedParams{
		ID:        f.ID,
		UpdatedAt: time.Now(),
	}
	if err := s.db.MarkFeedFetched(context.Background(), markParams); err != nil {
		return err
	}
	rssData, err := fetchFeed(context.Background(), f.Url)
	if err != nil {
		return err
	}
	fmt.Printf("Channel: %v\n---\nLink: %v\nDescription: %v\n---\n",
		rssData.Channel.Title,
		rssData.Channel.Link,
		rssData.Channel.Description,
	)
	for _, fi := range rssData.Channel.Item {
		fmt.Printf("%v\n", fi.Title)
	}
	return nil
}
