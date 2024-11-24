package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
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
	}

	return nil
}

// this function reads a file into memory, creates a sha-1 hash,
// creates the folders to hold it, and compresses it
func createBlob(entryPath string) error {
	log.Printf("[DEBUG] Creating blob for file %s\n", entryPath)

	file, err := os.Open(entryPath)
	if err != nil {
		return CommitError{err}
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return CommitError{err}
	}

	hasher := sha1.New()
	if _, err := hasher.Write(b); err != nil {
		return CommitError{err}
	}

	hash := hex.EncodeToString(hasher.Sum(nil))

	cwd, err := os.Getwd()
	if err != nil {
		return CommitError{err}
	}

	objectsDir := cwd + "/.gogit/objects/"
	hashDir := objectsDir + hash[:2]
	if err := os.Mkdir(hashDir, 0755); err != nil {
		return CommitError{err}
	}

	hashFilePath := hashDir + "/" + hash[2:]
	hashFile, err := os.Create(hashFilePath)
	if err != nil {
		return CommitError{err}
	}

	if _, err = hashFile.Write(b); err != nil {
		return CommitError{err}
	}

	log.Fatalf("[DEBUG] Successfully created blob for file %s\n", entryPath)

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

	entries, err := os.ReadDir(cwd)
	if err != nil {
		return CommitError{err}
	}

	for _, e := range entries {

		// NOTE: actually, when there's a directory, we want to recursively
		// create blobs for their inner files
		if e.IsDir() {
			continue
		}

		filePath := cwd + "/" + e.Name()
		if err := createBlob(filePath); err != nil {
			return CommitError{err}
		}
	}

	return nil
}
