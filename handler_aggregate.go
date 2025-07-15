package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"database/sql"

	"github.com/google/uuid"
	"github.com/mshagirov/gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("Error: missing required argument, \"time_between_reqs\"!\nE.g. 5s, 3m, 1h, ...\nUsage:\n  agg <time_between_reqs>")
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
	fmt.Printf(`Fetching: %v
 |- Link: %v
 |- Description: %v
 \__
`,
		rssData.Channel.Title,
		rssData.Channel.Link,
		rssData.Channel.Description,
	)
	var p database.Post
	var t, pubTime time.Time
	for _, feedItem := range rssData.Channel.Item {
		t = time.Now()
		pubTime, err = time.Parse(time.RFC1123Z, feedItem.PubDate)
		if err != nil {
			fmt.Println(err)
		}
		p, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   t,
			UpdatedAt:   t,
			Title:       feedItem.Title,
			Url:         feedItem.Link,
			Description: sql.NullString{Valid: true, String: feedItem.Description},
			PublishedAt: pubTime,
			FeedID:      f.ID,
		})
		if err != nil {
			if !strings.Contains(err.Error(), "duplicate key value") {
				fmt.Printf("%v\n", err)
			}
		} else {
			previewPost(p)
		}
		// fmt.Println(feedItem.PubDate)
	}
	return nil
}

func previewPost(p database.Post) {
	fmt.Printf(`    |- Title:   %v
    |  Updated: %v
    |
`, p.Title, p.UpdatedAt)
}
