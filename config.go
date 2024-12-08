package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var ConfigCommand = &cli.Command{
	Name:  "config",
	Usage: "Set and get values from configuration files",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name: "local",
		},
		&cli.BoolFlag{
			Name: "global",
		},
		&cli.BoolFlag{
			Name: "system",
		},
	},
	Action: Config,
}

type ConfigError struct {
	err error
}

func (ce ConfigError) Error() string {
	return fmt.Sprintf("[ERROR] Error in Config command: %s", ce.err)
}

func Config(c *cli.Context) error {

	// _ := c.Bool("local")
	global := c.Bool("global")
	system := c.Bool("system")
	var filepath string

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	filepath = cwd + "/.git/config"

	if global {
		filepath = os.Getenv("HOME") + "/.gitconfig"
	} else if system {
		filepath = "/etc/gitconfig"
	}

	log.Printf("The config file is %s\n", filepath)

	return nil
}
