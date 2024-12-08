package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

var Commands = []*cli.Command{
	CommitCommand,
	InitCommand,
	ConfigCommand,
}

var DefaultAction = func(c *cli.Context) error {
	fmt.Println("Hello world, I'm using urfave/cli for the first time!")
	return nil
}

func main() {

	app := &cli.App{
		Name:        "gogit",
		Description: "my own version of git",
		// the function called when no parameters or flags added
		Action: DefaultAction,
		// available commands after "gogit" (i.e commit, add, pull,)
		Commands: Commands,
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}

}
