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
