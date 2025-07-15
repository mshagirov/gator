package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mshagirov/gator/internal/database"
)

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("Error: too many arguments!\nUsage:\n  feeds")
	}
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Feed, URL, User, Created, Modified:")
	var u database.User
	for _, f := range feeds {
		u, err = s.db.GetUserWithID(context.Background(), f.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("%v, %v, %v, %v, %v\n",
			f.Name,
			f.Url,
			u.Name,
			f.CreatedAt.Format("2006-01-02 15:04:05"),
			f.UpdatedAt.Format("2006-01-02 15:04:05"),
		)
	}
	return nil
}

func handlerAddfeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("Error: missing required argument(s)!\nUsage:\n  addfeed <name> <url>")
	}
	time_now := time.Now()
	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time_now,
		UpdatedAt: time_now,
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	}
	f, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return err
	}
	// fmt.Printf("%+v\n", f)
	follow_params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time_now,
		UpdatedAt: time_now,
		UserID:    user.ID,
		FeedID:    f.ID,
	}
	feed_follow_row, err := s.db.CreateFeedFollow(context.Background(), follow_params)
	if err != nil {
		return err
	}
	fmt.Printf("%v subscribed to %v\n", feed_follow_row.UserName, feed_follow_row.FeedName)

	return nil
}
