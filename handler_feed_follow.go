package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mshagirov/gator/internal/database"
)

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("Error: missing required argument!\nUsage:\n  unfollow <url>")
	}
	url := cmd.Args[0]
	f, err := s.db.GetFeedWithUrl(context.Background(), url)
	if err != nil {
		return err
	}
	params := database.DeleteFeedForUserIdParams{
		UserID: user.ID,
		FeedID: f.ID,
	}
	if err := s.db.DeleteFeedForUserId(context.Background(), params); err != nil {
		return err
	}
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("Error: too many arguments!\nUsage:\n  following")
	}
	u_follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return err
	}
	for _, v := range u_follows {
		fmt.Printf("- \"%v\" : \"%v\"\n", v.FeedName, v.FeedUrl)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("Error: missing required argument!\nUsage:\n  follow <url>")
	}
	url := cmd.Args[0]
	f, err := s.db.GetFeedWithUrl(context.Background(), url)
	if err != nil {
		return err
	}
	time_now := time.Now()
	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time_now,
		UpdatedAt: time_now,
		UserID:    user.ID,
		FeedID:    f.ID,
	}
	feed_follow_row, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}
	fmt.Printf("%v subscribed to %v\n", feed_follow_row.UserName, feed_follow_row.FeedName)
	return nil
}
