package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/mshagirov/gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	var postLimit int32

	if len(cmd.Args) > 1 {
		return fmt.Errorf("Error: too many arguments!\nUsage:\n  browse\n  browse <NUMBER_OF_POSTS>")
	} else if len(cmd.Args) == 1 {
		postLimit64, err := strconv.ParseInt(cmd.Args[0], 10, 32)
		if err != nil {
			return err
		}
		postLimit = int32(postLimit64)
	} else {
		postLimit = 2
	}
	posts, err := s.db.GetPostsForUserId(context.Background(),
		database.GetPostsForUserIdParams{
			UserID: user.ID,
			Limit:  postLimit,
		})
	if err != nil {
		return err
	}
	for _, p := range posts {
		printPost(p)
	}
	return nil
}

func printPost(p database.GetPostsForUserIdRow) {
	var desc string
	if p.Description.Valid {
		desc = p.Description.String
	} else {
		desc = ""
	}
	fmt.Printf(`
>  %v
   Feed: %v
   PubDate: %v
   URL: %v
   ---
   %v
`,
		p.Title,
		p.FeedUrl,
		p.PublishedAt,
		p.Url,
		desc,
	)
}
