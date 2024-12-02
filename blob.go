package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"
)

type Blob struct {
	contents io.Reader
	hash     string
}

func NewBlob(reader io.Reader) *Blob {

	b, err := io.ReadAll(reader)
	if err != nil {
		log.Fatalf("[FATAL] [NewBlob] Error creating blob: %s\n", err)
	}

	hasher := sha1.New()
	_, err = hasher.Write(b)
	if err != nil {
		log.Fatalf("[FATAL] [NewBlob] Error writting to sha-1 hasher: %s\n", err)
	}

	return &Blob{
		hash:     hex.EncodeToString(hasher.Sum(nil)),
		contents: bytes.NewBuffer(b),
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
