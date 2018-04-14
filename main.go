package main

import (
	"fmt"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

type Uuided struct {
	Uuid string `json:"uuid"`
}

func (u Uuided) OpenFile() (*os.File, error) {
	return os.Create(u.Uuid)
}

var cache = []Uuided{Uuided{"123"}, Uuided{"123a"}, Uuided{"123b"}}

func main() {
	http.HandleFunc("/", ServeElements)
	log.Fatal(http.ListenAndServe(":12345", nil))
}

func ServeElements(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		file := req.URL.Query().Get("uuid")
		log.Print("Open ", file)
		f, err := os.Open(file)
		defer func() {
			f.Sync()
			f.Close()
		}()
		if err != nil {
			log.Panic(err)
		}
		io.Copy(w, f)
	} else if req.Method == "PUT" {
		file := req.URL.Query().Get("uuid")
		log.Print("Open ", file)
		f, err := os.Create(file)
		defer func() {
			f.Sync()
			f.Close()
		}()
		if err != nil {
			log.Panic(err)
		}
		io.Copy(f, req.Body)
	}
}

func write() {
	for _, uuided := range cache {
		log.Print(uuided)
		f, err := uuided.OpenFile()
		defer func() {
			f.Sync()
			f.Close()
		}()
		if err != nil {
			log.Print(err)
		}
		encoder := json.NewEncoder(f)
		encoder.SetIndent("", "\t")
		encoder.Encode(uuided)
	}
}

func findByUUID(uuid string) (* os.File, error) {
	return findByProperty("uuid", fmt.Sprintf("^%s$", uuid))
}

func findByProperty(propertyPath, regexp string) (* os.File, error) {
	return nil, nil
}

func read() {
	for _, file := range os.Args[1:] {
		log.Print("Open ", file)
		f, err := os.Open(file)
		defer func() {
			f.Sync()
			f.Close()
		}()
		if err != nil {
			log.Panic(err)
		}
		encoder := json.NewDecoder(f)
		uuided := Uuided{}
		encoder.Decode(&uuided)
		log.Print(uuided.Uuid)
	}
}
