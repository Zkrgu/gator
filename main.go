package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/zkrgu/gator/internal/config"
	"github.com/zkrgu/gator/internal/database"
	"github.com/zkrgu/gator/internal/state"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		fmt.Printf("Could not read config file\n")
		os.Exit(1)
	}

	db, err := sql.Open("postgres", conf.DBURL)
	if err != nil {
		fmt.Printf("Could not connect to databse\n")
		os.Exit(1)
	}

	dbQueries := database.New(db)

	s := state.State{
		Config:       &conf,
		DBConnection: dbQueries,
	}

	cmds := state.Commands{
		Commands: make(map[string]func(*state.State, state.Command) error),
	}

	cmds.Register("login", state.HandlerLogin)
	cmds.Register("register", state.HandlerRegister)
	cmds.Register("reset", state.HandlerReset)
	cmds.Register("users", state.HandlerUsers)
	cmds.Register("agg", state.HandlerAgg)
	cmds.Register("addfeed", state.HandlerAddFeed)
	cmds.Register("feeds", state.HandlerFeeds)
	cmds.Register("following", state.HandlerFollowing)
	cmds.Register("follow", state.HandlerFollow)
	cmds.Register("unfollow", state.HandlerUnfollow)
	cmds.Register("browse", state.HandlerBrowse)

	cmd := state.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	err = cmds.Run(&s, cmd)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

}
