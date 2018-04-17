package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type Query struct {
	Key, Val string
} // Should be a map[string]string in the future to and/or queries

// Repo is an interface that specifies methods to obtain io.ReadCloser
// and io.WriteCloser for elements and delete elements. Elemnts might be files
// or cache entries, by some identifiyed, e.g. UUID.
type Repo interface {
	Writer(string) (io.WriteCloser, error)
	Reader(string) (io.ReadCloser, error)
	QueryReader(Query) (io.ReadCloser, error)
	Remove(string) error
}

// FileRepo is a Repo to store elements to disk as files.
type FileRepo struct {
	dir string
}

func NewFileRepo(dir string) *FileRepo {
	return &FileRepo{dir}
}

func (*FileRepo) Writer(uuid string) (io.WriteCloser, error) {
	return os.Create(uuid)
}

func (*FileRepo) Reader(uuid string) (io.ReadCloser, error) {
	return os.Open(uuid)
}

func (fr *FileRepo) QueryReader(q Query) (io.ReadCloser, error) {
	files, err := ioutil.ReadDir(fr.dir)
	if err != nil {
		log.Fatal(err)
	}

	fileNames := make([]string, len(files)+2)
	fileNames[0] = q.Key
	fileNames[1] = q.Val
	for i, file := range files {
		fileNames[i+2] = file.Name()
	}
	cmd := exec.Command("./jrep", fileNames...)
	jqCmd := exec.Command("xargs", "jq", "-s", ".")

	in, out := io.Pipe()
	cmd.Stdout = out
	jqCmd.Stdin = in

	var result bytes.Buffer
	jqCmd.Stdout = &result

	cmd.Start()
	jqCmd.Start()
	cmd.Wait()
	out.Close()
	jqCmd.Wait()

// 	out, _ := cmd.StdoutPipe()
// 	cmd.Start()
//
// 	in, _ := jqCmd.StdinPipe()
// 	result, _ := jqCmd.StdoutPipe()
// 	jqCmd.Start()

// 	io.Copy(os.Stdout, &result)

	return &ReadCloserWrapper{&result}, nil
}

type ReadCloserWrapper struct {
	Reader io.Reader
}

func (r *ReadCloserWrapper) Read(p []byte) (int, error) {
	return r.Reader.Read(p)
}

func (r *ReadCloserWrapper) Close() (error) {
	return nil
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

func (c *Cache) QueryReader(q Query) (io.ReadCloser, error) {
	return c.NewQueryReader(q), nil
}

func (c *Cache) Remove(uuid string) error {
	delete(*c, uuid)
	return nil
}

type CacheReader struct {
	reader io.Reader
}

func (c *Cache) NewReader(key string) *CacheReader {
	return &CacheReader{
		reader: bytes.NewReader((*c)[key]),
	}
}

func (c *Cache) NewQueryReader(q Query) *CacheReader {
	entry := (*c)[q.Val] // We simply assume the Val is an existing UUID
	entryArray := fmt.Sprintf("[ %s ]", string(entry))
	return &CacheReader{
		reader: bytes.NewReader([]byte(entryArray)),
	}
}

func (cw *CacheReader) Read(p []byte) (int, error) {
	return cw.reader.Read(p)
}

func (cw *CacheReader) Close() error {
	return nil
}

type CacheWriter struct {
	writer *bytes.Buffer
	cache  *Cache
	key    string
}

func (c *Cache) NewWriter(key string) *CacheWriter {
	return &CacheWriter{
		writer: bytes.NewBuffer(nil),
		cache:  c,
		key:    key,
	}
}

func (cw *CacheWriter) Write(p []byte) (int, error) {
	return cw.writer.Write(p)
}

func (cw *CacheWriter) Close() error {
	(*cw.cache)[cw.key] = cw.writer.Bytes()
	return nil
}
