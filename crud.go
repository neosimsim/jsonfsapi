package main

import (
	"net/http"
	"io"
	"encoding/json"
	"log"
)

// Reads the JSON from the http.ResponseWriter deserialize to an interface{}
// and sets the property UUID to a new UUID. A given UUID, if any, will we overwritten.
// Encode the JSON to a file afterwards.
func CreateElements(repo Repo, w http.ResponseWriter, req *http.Request) {
	file := genUUID()
	log.Print("Open ", file)
	f, err := repo.Writer(file)
	defer func() {
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
	w.Header().Set("Location", file)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func StoreElement(uuid string, w io.Reader) error {
	// find JSON with propertis.uuid == uuid
	// if none, create new
	return nil
}

func LoadElement(uuid string, w io.Writer) error {
	// find JSON with propertis.uuid == uuid
	return nil
}

func ReadElements(repo Repo, w http.ResponseWriter, req *http.Request) {
	file := req.URL.Query().Get("uuid")
	log.Print("Open ", file)
	f, err := repo.Reader(file)
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
	uuid := req.URL.Query().Get("uuid")

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
	file := req.URL.Query().Get("uuid")
	log.Print("Delete ", file)
	err := repo.Remove(file)
	if err != nil {
		log.Print(err)
	}
	w.WriteHeader(http.StatusNoContent)
}
