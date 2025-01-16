package command

import (
	"fmt"
	"os"

	"github.com/EloyTolosaDev/gogit/pkg/object"

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

func Commit(c *cli.Context) error {

	cwd, err := os.Getwd()
	if err != nil {
		return CommitError{err}
	}

	t := object.NewTree(cwd)
	t.Save()

	return nil
}
