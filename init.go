package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var InitCommand = &cli.Command{
	Name:  "init",
	Usage: "Initialize gogit repository (create necessary files and folders)",
	// Flags: []cli.Flag{
	// 	&cli.StringFlag{
	// 		Name:    "message",
	// 		Aliases: []string{"m"},
	// 	},
	// },
	Action: Init,
}

func Init(c *cli.Context) error {
	var cwd string

	cwd, err := os.Getwd()
	if err != nil {
		return CommitError{err}
	}

	mainDir := ".gogit"
	files := []string{"HEAD", "config", "hooks/post-commit", "hooks/pre-commit", "info/exclude",
		"objects/info", "objects/pack"}
	dirs := []string{"index", "info", "objects", "refs", "hooks"}

	cwd = cwd + "/" + mainDir
	if err := os.Mkdir(cwd, 0755); err != nil {
		if !os.IsExist(err) {
			return CommitError{err}
		} else {
			log.Printf("[DEBUG] Dir %s already exists\n", cwd)
		}
	}

	for _, dir := range dirs {
		dirname := cwd + "/" + dir
		if err := os.Mkdir(dirname, 0755); err != nil {
			if !os.IsExist(err) {
				return CommitError{err}
			} else {
				log.Printf("[DEBUG] Dir %s already exists\n", dirname)
			}
		}
	}

	for _, file := range files {
		filename := cwd + "/" + file
		if _, err := os.Create(filename); err != nil {
			if !os.IsExist(err) {
				return CommitError{err}
			} else {
				log.Printf("[DEBUG] file %s already exists\n", filename)
			}
		}
	}

	return nil
}
