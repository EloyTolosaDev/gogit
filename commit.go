package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"runtime"
	"strings"

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

// this function creates a tree object from a directory
// a tree object is a file with a list of trees and blobs
// inside the folder with the following information:
// -- mode, type, name and sha
func createTreeWithDepth(dirpath string, depth int) error {
	if depth == MAX_RECURSIVE_DEPTH {
		return fmt.Errorf("[ERROR] Max recursive depth (%d) reached for func %s", MAX_RECURSIVE_DEPTH, getCurrentFunctionName())
	}

	objects := []fs.DirEntry{}
	treeBuilder := strings.Builder{}

	entries, err := os.ReadDir(dirpath)
	if err != nil {
		return CommitError{err}
	}

	// for every entry entry, check if its a file or a directory
	// and create, accordingly, a tree or a blob object
	for _, entry := range entries {
		if entry.Name() == ".git" || entry.Name() == ".gogit" {
			continue
		}

		entryPath := dirpath + "/" + entry.Name()
		if entry.IsDir() {
			createTreeWithDepth(entryPath, depth+1)
		} else {
			createBlob(entryPath)
		}
		objects = append(objects, entry)
	}

	for _, entry := range entries {
		t := "blob"
		if entry.IsDir() {
			t = "tree"
		}

		treeBuilder.WriteString(fmt.Sprintf("%s %s \t%s\n", t, hash, entry.Name()))
	}

	return nil
}

// calculates the sha-1 hash from the contents of the file
// and returns it
func hash(filepath string) (string, error) {

	file, err := os.Open(filepath)
	if err != nil {
		return "", CommitError{err}
	}
	defer file.Close()
	b, err := io.ReadAll(file)
	if err != nil {
		return "", CommitError{err}
	}

	hasher := sha1.New()
	if _, err := hasher.Write(b); err != nil {
		return "", CommitError{err}
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// this function reads a file into memory, creates a sha-1 hash,
// creates the folders to hold it, and compresses it
func createBlob(entryPath string) error {
	log.Printf("[DEBUG] Creating blob for file %s\n", entryPath)

	hash, err := hash(entryPath)
	if err != nil {
		return CommitError{err}
	}

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

// lists all dir files and creates blobs for every file in the dir
// if more dirs are found, traverse them
//
// NOTE This function is implemented like this to prevent accidental
// infinite recursive calls and throw an error when that happens
func createTree(dirpath string) error {
	return createTreeWithDepth(dirpath, 0)
}

func Commit(c *cli.Context) error {

	cwd, err := os.Getwd()
	if err != nil {
		return CommitError{err}
	}

	return createTree(cwd)
}
