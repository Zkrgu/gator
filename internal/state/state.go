package state

import (
	"fmt"

	"github.com/zkrgu/gator/internal/config"
)

type State struct {
	Config *config.Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Commands map[string]func(*State, Command) error
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("No arguments supplied\n")
	}
	err := s.Config.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Set user to %s\n", s.Config.CurrentUserName)

	return nil
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Commands[name] = f
}

func (c *Commands) Run(s *State, cmd Command) error {
	fun, ok := c.Commands[cmd.Name]

	if !ok {
		return fmt.Errorf("No Command with that name found\n")
	}

	return fun(s, cmd)
}
