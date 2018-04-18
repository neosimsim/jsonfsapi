package main

import (
	"encoding/json"
	"fmt"
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
	http.HandleFunc("/elements", ServeElementsFactory(NewFileRepo("elements")))
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("assets"))))

	log.Fatal(http.ListenAndServe(":12345", nil))
}

func ServeElementsFactory(repo Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case "GET":
			ReadElements(repo, w, req)
		case "POST":
			CreateElements(repo, w, req)
		case "PUT":
			UpdateElements(repo, w, req)
		case "DELETE":
			DeleteElements(repo, w, req)
		}
	}
}

func genUUID() string {
	f, _ := os.Open("/dev/urandom")
	b := make([]byte, 16)
	f.Read(b)
	f.Close()
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
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

func findByUUID(uuid string) (*os.File, error) {
	return findByProperty("uuid", fmt.Sprintf("^%s$", uuid))
}

func findByProperty(propertyPath, regexp string) (*os.File, error) {
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
