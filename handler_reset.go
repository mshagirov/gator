package main

import (
	"context"
	"fmt"
)

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
