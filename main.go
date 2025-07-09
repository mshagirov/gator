package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/mshagirov/gator/internal/config"
	"github.com/mshagirov/gator/internal/database"
)

func main() {
	Config := config.Read()
	db, err := sql.Open("postgres", Config.DbURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dbQueries := database.New(db)

	State := state{
		db:     dbQueries,
		Config: &Config,
	}

	Commands := getCommands()

	if len(os.Args) < 2 {
		fmt.Println("Error: missing a command!")
		os.Exit(1)
	}

	cmd := command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	if err := Commands.run(&State, cmd); err != nil {
		fmt.Println(err)
		fmt.Println("")
		os.Exit(1)
	}
}
