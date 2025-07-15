package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/mshagirov/gator/internal/config"
	"github.com/mshagirov/gator/internal/database"
)

type state struct {
	db     *database.Queries
	Config *config.Config
}

type command struct {
	Name string
	Args []string
}

type commands struct {
	Command map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if _, exists := c.Command[cmd.Name]; !exists {
		return fmt.Errorf("Unknown command: %s", cmd.Name)
	}
	if err := c.Command[cmd.Name](s, cmd); err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) error {
	c.Command[name] = f
	return nil
}

func getCommands() commands {
	Commands := commands{
		Command: make(map[string]func(*state, command) error),
	}
	// Register CLI
	Commands.register("login", handlerLogin)
	Commands.register("register", handlerRegister)
	Commands.register("reset", handlerReset)
	Commands.register("users", handlerUsers)
	Commands.register("agg", handlerAgg)
	Commands.register("feeds", handlerFeeds)
	Commands.register("addfeed", middlewareLoggedIn(handlerAddfeed))
	Commands.register("follow", middlewareLoggedIn(handlerFollow))
	Commands.register("following", middlewareLoggedIn(handlerFollowing))
	Commands.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	return Commands
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	dryHandler := func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.Config.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
	return dryHandler
}

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

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	var u_suffix string
	for _, u := range users {
		u_suffix = ""
		if u == s.Config.CurrentUserName {
			u_suffix = " (current)"
		}
		fmt.Printf("* %s%s\n", u, u_suffix)
	}
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf(`Error: missing required argument!
Usage:
  register <name>`)
	}
	if _, err := s.db.GetUser(context.Background(), cmd.Args[0]); err == nil {
		return fmt.Errorf("User %s already exists.", cmd.Args[0])
	}
	time_now := time.Now()
	params := database.CreateUserParams{
		ID:        uuid.New(), //int32,
		CreatedAt: time_now,   //time.Time
		UpdatedAt: time_now,
		Name:      cmd.Args[0],
	}
	u, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return err
	}
	if err := s.Config.SetUser(u.Name); err != nil {
		return err
	}
	*s.Config = config.Read()
	fmt.Printf("User %s with uuid=%v created at %s",
		u.Name, u.ID, u.CreatedAt.Format("2006-01-02 15:04:05"),
	)
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf(`Error: missing required argument!
Usage:
  login <username>`)
	}
	// Try to query user from DB
	u, err := s.db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		fmt.Printf("User %s does NOT exists! Please register the user before login in.", cmd.Args[0])
		os.Exit(1)
	}
	if err := s.Config.SetUser(u.Name); err != nil {
		return err
	}
	*s.Config = config.Read()
	fmt.Printf("User is set to: %s\n", s.Config.CurrentUserName)
	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf(`Error:"reset" does NOT accept arguments!
Usage:
  reset`)
	}
	if err := s.db.Reset(context.Background()); err != nil {
		return err
	}
	fmt.Println("!!!Deleted all users from the database!!!")
	return nil
}
