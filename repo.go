package main

import (
	"io"
	"os"
)

// Repo is an interface that specifies methods to obtain io.ReadCloser
// and io.WriteCloser for elements and delete elements. Elemnts might be files
// or cache entries, by some identifiyed, e.g. UUID.
type Repo interface {
	Writer(string) (io.WriteCloser, error)
	Reader(string) (io.ReadCloser, error)
	Remove(string) error
}

// FileRepo is a Repo to store elements to disk as files.
type FileRepo struct{}

func NewFileRepo() *FileRepo {
	return &FileRepo{}
}

func (*FileRepo) Writer(uuid string) (io.WriteCloser, error) {
	return os.Create(uuid)
}

func (*FileRepo) Reader(uuid string) (io.ReadCloser, error) {
	return os.Open(uuid)
}

func (*FileRepo) Remove(uuid string) error {
	return os.Remove(uuid)
}
