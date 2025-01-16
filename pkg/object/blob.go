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

type Blob struct {
	sourcePath string
	contents   io.Reader
	hash       string
}

func (b *Blob) Hash() string {
	return b.hash
}

func (b *Blob) Info() string {
	return fmt.Sprintf("Blob %s \t%s\n", b.Hash(), b.Name())
}

func (b *Blob) Name() string {
	parts := strings.Split(b.sourcePath, "/")
	return parts[len(parts)-1]
}

func NewBlob(filepath string) *Blob {
	file, err := os.Open(filepath)
	if err != nil {
		return nil
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("[FATAL] [NewBlob] Error creating blob: %s\n", err)
	}

	hasher := sha1.New()
	_, err = hasher.Write(b)
	if err != nil {
		log.Fatalf("[FATAL] [NewBlob] Error writting to sha-1 hasher: %s\n", err)
	}

	return &Blob{
		sourcePath: filepath,
		hash:       hex.EncodeToString(hasher.Sum(nil)),
		contents:   bytes.NewBuffer(b),
	}

}

func (b *Blob) Save() {

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("[FATAL] Error saving blob: %s\n", err)
	}

	objectsDir := cwd + "/.gogit/objects/"
	hashDir := objectsDir + b.hash[:2]
	if err := os.Mkdir(hashDir, 0755); err != nil {
		if !os.IsExist(err) {
			log.Fatalf("[FATAL] Error saving blob: %s\n", err)
		}
	}

	hashFilePath := hashDir + "/" + b.hash[2:]
	hashFile, err := os.Create(hashFilePath)
	if err != nil {
		if !os.IsExist(err) {
			log.Fatalf("[FATAL] Error saving blob: %s\n", err)
		}
	}

	if _, err = io.Copy(hashFile, b.contents); err != nil {
		log.Fatalf("[FATAL] Error saving blob: %s\n", err)
	}

}
