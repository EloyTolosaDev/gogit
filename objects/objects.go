package objects

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	MAX_RECURSIVE_DEPTH = 100
)

type Object interface {
	Info() string // prints its info -> needed for a tree to save it's child object's info
	Save() error  // saves to disk
	Hash() string // generates hash with it's current buffer -> needeed for Save
	Type() string // get's the object type -> Blob, Tree or Commit (for the moment)
}

type Blob struct {
	file io.ReadWriter
	Path string // needed for Info() method -> maybe in the future we make it "private"
}

func NewBlob(filepath string) (*Blob, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("[DEBUG] [NewBlob] error creating blob: %s", err.Error())
	}
	return &Blob{f, filepath}, nil
}

func (b *Blob) Info() string {
	return fmt.Sprintf("%s %s \t%s\n", b.Type(), b.Hash(), b.Name())
}

func (b *Blob) Type() string {
	return "Blob"
}

func (b *Blob) Name() string {
	parts := strings.Split(b.Path, "/")
	return parts[len(parts)-1]
}

func (b *Blob) Hash() string {
	bytes, err := io.ReadAll(b.file)
	if err != nil {
		return ""
	}

	hasher := sha1.New()
	if _, err := hasher.Write(bytes); err != nil {
		return ""
	}

	return hex.EncodeToString(hasher.Sum(nil))
}

func (b *Blob) Save() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	hash := b.Hash()

	objectsDir := cwd + "/.gogit/objects/"
	hashDir := objectsDir + hash[:2]
	if err := os.Mkdir(hashDir, 0755); err != nil {
		if !os.IsExist(err) {
			return err
		}
	}

	hashFilePath := hashDir + "/" + hash[2:]
	hashFile, err := os.Create(hashFilePath)
	if err != nil {
		if !os.IsExist(err) {
			return err
		}
	}
	defer hashFile.Close()

	_, err = io.Copy(hashFile, b.file)
	if err != nil {
		return err
	}

	return nil
}

type Tree struct {
	dirpath string
	objs    []Object
	reader  io.ReadWriter
}

func NewTree(dirpath string) *Tree {
	return newTree(dirpath, 0)
}

func newTree(dirpath string, depth int) *Tree {
	if depth == MAX_RECURSIVE_DEPTH {
		return nil //, fmt.Errorf("[ERROR] Max recursive depth (%d) reached for func %s", MAX_RECURSIVE_DEPTH, getCurrentFunctionName())
	}

	entries, err := os.ReadDir(dirpath)
	if err != nil {
		return nil
	}

	t := &Tree{
		dirpath: dirpath,
	}

	// for every entry entry, check if its a file or a directory
	// and create, accordingly, a tree or a blob object
	var o Object
	for _, entry := range entries {
		if entry.Name() == ".git" || entry.Name() == ".gogit" {
			continue
		}

		entryPath := dirpath + "/" + entry.Name()
		if entry.IsDir() {
			o = newTree(entryPath, depth+1)
		} else {
			o, err = NewBlob(entryPath)
		}
		if err != nil {
			return nil
		}

		t.AddObject(o)
	}

	return t
}

func (t *Tree) AddObject(object Object) {
	t.objs = append(t.objs, object)
}

func (t *Tree) Info() string {
	return ""
}

func (t *Tree) Save() error {
	return nil
}

func (t *Tree) Hash() string {
	return ""
}

func (t *Tree) Type() string {
	return "Tree"
}
