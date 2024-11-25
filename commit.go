package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"

	"github.com/urfave/cli/v2"
)

const (
	MAX_RECURSIVE_DEPTH = 100
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
		if !os.IsExist(err) {
			return CommitError{err}
		}
	}

	hashFilePath := hashDir + "/" + hash[2:]
	hashFile, err := os.Create(hashFilePath)
	if err != nil {
		if !os.IsExist(err) {
			return CommitError{err}
		}
	}

	if _, err = hashFile.Write(b); err != nil {
		return CommitError{err}
	}

	log.Printf("[DEBUG] Successfully created blob for file %s\n", entryPath)

	return nil
}

func getCurrentFunctionName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}

func traverseDirDepth(dirpath string, depth int) error {
	if depth == MAX_RECURSIVE_DEPTH {
		return fmt.Errorf("[ERROR] Max recursive depth (%d) reached for func %s", MAX_RECURSIVE_DEPTH, getCurrentFunctionName())
	}

	entries, err := os.ReadDir(dirpath)
	if err != nil {
		return CommitError{err}
	}

	for _, e := range entries {

		if e.Name() == ".git" || e.Name() == ".gogit" {
			continue
		}

		if e.IsDir() {
			log.Printf("[DEBUG] Traversing directory %s:\n", e.Name())

			newDirPath := dirpath + "/" + e.Name()
			traverseDirDepth(newDirPath, depth+1)
			continue
		}

		filePath := dirpath + "/" + e.Name()
		if err := createBlob(filePath); err != nil {
			return CommitError{err}
		}
	}

	return nil
}

// lists all dir files and creates blobs for every file in the dir
// if more dirs are found, traverse them
//
// NOTE This function is implemented like this to prevent accidental
// infinite recursive calls and throw an error when that happens
func traverseDir(dirpath string) error {
	return traverseDirDepth(dirpath, 0)
}

func Commit(c *cli.Context) error {

	cwd, err := os.Getwd()
	if err != nil {
		return CommitError{err}
	}

	return traverseDir(cwd)
}
