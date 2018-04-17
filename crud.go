package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"path"
)

// Reads the JSON from the http.ResponseWriter deserialize to an interface{}
// and sets the property UUID to a new UUID. A given UUID, if any, will we overwritten.
// Encode the JSON to a file afterwards.
func CreateElements(repo Repo, w http.ResponseWriter, req *http.Request) {
	uuid := genUUID()
	log.Print("Open ", uuid)
	f, err := repo.Writer(uuid)
	defer func() {
		f.Close()
	}()
	if err != nil {
		log.Panic(err)
	}
	decoder := json.NewDecoder(req.Body)
	var jsonObj interface{}
	decoder.Decode(&jsonObj)
	log.Print("Store new JSON: ", jsonObj, " whith UUID ", uuid)
	jsonObj.(map[string]interface{})["uuid"] = uuid
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "\t")
	encoder.Encode(jsonObj)
	w.Header().Set("Location", uuid)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func ReadElement(repo Repo, w http.ResponseWriter, req *http.Request) {
	uuid := path.Base(req.URL.Path)
	log.Print("Open ", uuid)
	f, err := repo.Reader(uuid)
	defer func() {
		f.Close()
	}()
	if err != nil {
		log.Panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, f)
}

func ReadElements(repo Repo, w http.ResponseWriter, req *http.Request) {
	uuid := req.URL.Query().Get("uuid")
	query := Query{"uuid", uuid}
	log.Print("Open ", uuid)
	log.Print("query: ", query)
	f, err := repo.QueryReader(query)
	defer func() {
		f.Close()
	}()
	if err != nil {
		log.Panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, f)
}

func UpdateElements(repo Repo, w http.ResponseWriter, req *http.Request) {
	uuid := path.Base(req.URL.Path)

	log.Print("Open ", uuid)
	f, err := repo.Writer(uuid)
	defer func() {
		f.Close()
	}()
	if err != nil {
		log.Panic(err)
	}

	// make sure the UUID is not overwritten by the request
	decoder := json.NewDecoder(req.Body)
	var jsonObj map[string]interface{}
	decoder.Decode(&jsonObj)
	jsonObj["uuid"] = uuid

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "\t")
	encoder.Encode(jsonObj)
	io.Copy(f, req.Body)

	w.WriteHeader(http.StatusNoContent)
}

func DeleteElements(repo Repo, w http.ResponseWriter, req *http.Request) {
	uuid := path.Base(req.URL.Path)
	log.Print("Delete ", uuid)
	err := repo.Remove(uuid)
	if err != nil {
		log.Print(err)
	}
	w.WriteHeader(http.StatusNoContent)
}
