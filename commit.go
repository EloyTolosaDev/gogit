package main

import (
	"fmt"
	"os"
	"runtime"

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
	Action: Commit,
}

type CommitError struct {
	err error
}

func (e CommitError) Error() string {
	return fmt.Sprintf("[ERROR] Error in commit: %s\n", e.err.Error())
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

	t := NewTree(cwd)
	t.Save()

	return nil
}
