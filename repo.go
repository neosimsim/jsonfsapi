package main

import (
	"io"
	"os"
)

type Repo interface {
	Writer(string) (io.WriteCloser, error)
	Reader(string) (io.ReadCloser, error)
	Remove(string) error
}

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
