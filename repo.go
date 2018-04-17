package main

import (
	"bytes"
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

type Cache map[string][]byte

func (c *Cache) Writer(uuid string) (io.WriteCloser, error) {
	return c.NewWriter(uuid), nil
}

func (c *Cache) Reader(uuid string) (io.ReadCloser, error) {
	return c.NewReader(uuid), nil
}

func (c *Cache) Remove(uuid string) error {
	delete(*c, uuid)
	return nil
}

type CacheReader struct{
	reader io.Reader
}

func (c *Cache) NewReader(key string) *CacheReader {
	return &CacheReader{
		reader: bytes.NewReader((*c)[key]),
	}
}

func (cw *CacheReader) Read(p []byte) (int, error) {
	return cw.reader.Read(p)
}

func (cw *CacheReader) Close() error {
	return nil
}

type CacheWriter struct{
	writer *bytes.Buffer;
	cache *Cache;
	key string;
}

func (c *Cache) NewWriter(key string) *CacheWriter {
	return &CacheWriter{
		writer: bytes.NewBuffer(nil),
		cache: c,
		key: key,
	}
}

func (cw *CacheWriter) Write(p []byte) (int, error) {
	return cw.writer.Write(p)
}

func (cw *CacheWriter) Close() error {
	(*cw.cache)[cw.key] = cw.writer.Bytes()
	return nil
}
