package object

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Tree struct {
	path     string
	objs     []Object
	contents io.ReadWriteCloser
	hash     string
}

func (t *Tree) Hash() string {
	return t.hash
}

func (t *Tree) Info() string {
	return fmt.Sprintf("Tree %s \t%s\n", t.Hash(), t.DirName())
}

func (t *Tree) DirName() string {
	parts := strings.Split(t.path, "/")
	return parts[len(parts)-1]
}

const (
	MAX_RECURSIVE_DEPTH = 100
)

func NewTree(dirpath string) *Tree {
	return newTreeRecursive(dirpath, 1)
}

func newTreeRecursive(dirpath string, depth int) *Tree {
	if depth == MAX_RECURSIVE_DEPTH {
		// NOTE log something
		return nil
	}

	entries, err := os.ReadDir(dirpath)
	if err != nil {
		log.Fatalf("[FATAL] Error when creating Tree object: %s\n", err)
		return nil
	}

	objs := make([]Object, 0)
	var o Object
	for _, e := range entries {

		// FIX Instead of this, the program should read the contents of .gogitignore (xd)
		// and avoid the names inside the file
		if e.Name() == ".git" || e.Name() == ".gogit" {
			continue
		}

		if e.IsDir() {
			log.Printf("[DEBUG] Traversing directory %s:\n", e.Name())

			newDirPath := dirpath + "/" + e.Name()
			o = newTreeRecursive(newDirPath, depth+1)
		} else {
			filePath := dirpath + "/" + e.Name()
			b := NewBlob(filePath)
			b.Save()
			o = b
		}

		objs = append(objs, o)

	}

	return &Tree{
		path: dirpath,
		objs: objs,
	}
}

func (t *Tree) Save() {

	buffer := &bytes.Buffer{}
	for _, o := range t.objs {
		if _, err := buffer.WriteString(o.Info()); err != nil {
			log.Fatalf("[FATAL] Error writting object into tree file: %s\n", err)
		}
	}

	hasher := sha1.New()
	if _, err := hasher.Write(buffer.Bytes()); err != nil {
		log.Fatalf("[FATAL] Error writting tree file contents: %s\n", err)
	}

	t.hash = hex.EncodeToString(hasher.Sum(nil))

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("[FATAL] Error saving blob: %s\n", err)
	}

	objectsDir := cwd + "/.gogit/objects/"
	hashDir := objectsDir + t.hash[:2]
	if err := os.Mkdir(hashDir, 0755); err != nil {
		if !os.IsExist(err) {
			log.Fatalf("[FATAL] Error saving blob: %s\n", err)
		}
	}

	hashFilePath := hashDir + "/" + t.hash[2:]
	hashFile, err := os.Create(hashFilePath)
	if err != nil {
		if !os.IsExist(err) {
			log.Fatalf("[FATAL] Error saving blob: %s\n", err)
		}
	}

	if _, err = io.Copy(hashFile, buffer); err != nil {
		log.Fatalf("[FATAL] Error saving blob: %s\n", err)
	}

}
