package main

import (
	"encoding/json"
	"fmt"
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
	switch req.Method {
	case "GET":
		ReadElements(w, req)
	case "POST":
		CreateElements(w, req)
	case "PUT":
		UpdateElements(w, req)
	case "DELETE":
		DeleteElements(w, req)
	}
}

func genUUID() string {
	return "12123j1230123"
}

// Reads the JSON from the http.ResponseWriter deserialize to an interface{}
// and sets the property UUID to a new UUID. A given UUID, if any, will we overwritten.
// Encode the JSON to a file afterwards.
func CreateElements(w http.ResponseWriter, req *http.Request) {
	file := genUUID()
	log.Print("Open ", file)
	f, err := os.Create(file)
	defer func() {
		f.Sync()
		f.Close()
	}()
	if err != nil {
		log.Panic(err)
	}
	decoder := json.NewDecoder(req.Body)
	var jsonObj interface{}
	decoder.Decode(&jsonObj)
	log.Print(jsonObj)
	jsonObj.(map[string]interface{})["uuid"] = file
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "\t")
	encoder.Encode(jsonObj)
}

func ReadElements(w http.ResponseWriter, req *http.Request) {
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
}

func UpdateElements(w http.ResponseWriter, req *http.Request) {
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

func DeleteElements(w http.ResponseWriter, req *http.Request) {
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
