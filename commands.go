package main

import (
	"context"
	"fmt"

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
