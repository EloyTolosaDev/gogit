package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var CommitCommand = &cli.Command{
	Name:  "commit",
	Usage: "Create a commit",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "message",
			Aliases: []string{"m"},
		},
	},
	Before: BeforeCommit,
	Action: Commit,
}

type CommitError struct {
	err error
}

func (e CommitError) Error() string {
	return fmt.Sprintf("[ERROR] Error in commit: %s\n", e.err.Error())
}

// Create .gogit dir and other dirs necessary to the command functionality
func BeforeCommit(c *cli.Context) error {
	var cwd string

	cwd, err := os.Getwd()
	if err != nil {
		return CommitError{err}
	}

	cwd = cwd + "/.gogit"
	if err := os.Mkdir(cwd, 0755); err != nil {
		if !os.IsExist(err) {
			return CommitError{err}
		} else {
			log.Printf("[DEBUG] Dir %s already exists\n", cwd)
		}
	}

	// create .gogit subfolders
	obj_folder := cwd + "/objects"
	if err := os.Mkdir(obj_folder, 0755); err != nil {
		if !os.IsExist(err) {
			return CommitError{err}
		} else {
			log.Printf("[DEBUG] Dir %s already exists\n", obj_folder)
		}
		return CommitError{err}
	}

	return nil
}

func Commit(c *cli.Context) error {

	// logic:
	// traverse all files in the directory sequentially, read their contents,
	// calculate their hsa-1 hash, and create the desired folder and object
	var cwd string

	cwd, err := os.Getwd()
	if err != nil {
		return CommitError{err}
	}

	log.Printf("[DEBUG] Listing files in the %s dir...\n", cwd)

	return nil
}
