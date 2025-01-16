package main

import (
	"fmt"
	"os"

	"github.com/zkrgu/gator/internal/config"
	"github.com/zkrgu/gator/internal/state"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		fmt.Printf("Could not read config file\n")
		os.Exit(1)
	}
	fmt.Printf("%v\n", conf)

	s := state.State{
		Config: &conf,
	}

	cmds := state.Commands{
		Commands: make(map[string]func(*state.State, state.Command) error),
	}

	cmds.Register("login", state.HandlerLogin)

	cmd := state.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	err = cmds.Run(&s, cmd)
	if err != nil {
		os.Exit(1)
	}

}
