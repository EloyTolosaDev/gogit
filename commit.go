package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/EloyTolosaDev/gogit/objects"
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
	}

	return nil
}

// this function creates a tree object from a directory
// a tree object is a file with a list of trees and blobs
// inside the folder with the following information:
// -- mode, type, name and sha
func createTree(dirpath string) (*objects.Tree, error) {
	t := objects.NewTree(dirpath)
	if err := t.Save(); err != nil {
		return nil, err
	}
	return t, nil
}

func getCurrentFunctionName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}

func Commit(c *cli.Context) error {

	cwd, err := os.Getwd()
	if err != nil {
		return CommitError{err}
	}

	_, err = createTree(cwd)
	if err != nil {
		return CommitError{err}
	}

	return nil
}
