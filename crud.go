package main

import (
	"net/http"
	"io"
	"encoding/json"
	"log"
	"os"
)

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
	log.Print("Store new JSON: ", jsonObj, " whith UUID ", file)
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
	file := req.URL.Query().Get("uuid")
	log.Print("Delete ", file)
	err := os.Remove(file)
	if err != nil {
		log.Print(err)
	}
}
