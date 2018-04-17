package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreate(t *testing.T) {
	cache := Cache{}

	message := `{ "title": "I'm new" }`
	req := httptest.NewRequest("POST", "http://localhost", strings.NewReader(message))
	w := httptest.NewRecorder()
	ServeElementsFactory(&cache)(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusCreated {
		t.Error("Expected status code ", http.StatusCreated, " got ", resp.StatusCode)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Error("Content-Type unset or unequal to 'application/json'")
	}

	location := w.Header().Get("Location")
	if location == "" {
		t.Fatal("Location header is missing")
	}

	uuid := location
	entry, exists := cache[uuid]
	if !exists {
		t.Error("Element has not been stored to repo.")
	}

	var result map[string]interface{}
	if err := json.Unmarshal(entry, &result); err != nil {
		t.Fatal("Stored element is not json serializable, ", err, entry)
	}
	if result["uuid"] != uuid {
		t.Error("Stored element ", result, " should have the returned UUID ", uuid)
	}
}

func TestRead(t *testing.T) {
	cache := Cache{"deadbeef": []byte(`{"uuid": "deadbeef", "title": "Get me"}`)}

	req := httptest.NewRequest("GET", "http://localhost?uuid=deadbeef", nil)
	w := httptest.NewRecorder()
	ServeElementsFactory(&cache)(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Error("Expected status code ", http.StatusOK, " got ", resp.StatusCode)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Error("Content-Type unset or unequal to 'application/json'")
	}

	decoder := json.NewDecoder(w.Body)
	var result map[string]interface{}

	if err := decoder.Decode(&result); err != nil {
		t.Fatal("Response is not JSON serializable: ", err)
	}
	if result["uuid"] != "deadbeef" {
		t.Error("The UUID should not be updated.")
	}
	if result["title"] != "Get me" {
		t.Error("The content should be updated.")
	}
}

func TestUpdate(t *testing.T) {
	uuid := "deadbeef"
	cache := Cache{uuid: []byte(`{"uuid": "deadbeef", "title": "Update me"}`)}

	message := `{"uuid": "abad1dea", "title": "Updated"}`
	req := httptest.NewRequest("PUT", "http://localhost?uuid=deadbeef", strings.NewReader(message))
	w := httptest.NewRecorder()
	ServeElementsFactory(&cache)(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusNoContent {
		t.Error("Expected status code ", http.StatusNoContent, " got ", resp.StatusCode)
	}

	entry, exists := cache[uuid]
	if !exists {
		t.Error("Element has vanished.")
	}

	var result map[string]interface{}
	if err := json.Unmarshal(entry, &result); err != nil {
		t.Fatal("Stored element is not json serializable. Error:", err, ". Entry: ", string(entry))
	}
	if result["uuid"] != uuid {
		t.Error("The UUID should not be updated.")
	}
	if result["title"] != "Updated" {
		t.Error("The content should be updated.")
	}
}

func TestDelete(t *testing.T) {
	cache := Cache{"deadbeef": []byte("Hello World!")}

	req := httptest.NewRequest("DELETE", "http://localhost?uuid=deadbeef", strings.NewReader("{}"))
	w := httptest.NewRecorder()

	ServeElementsFactory(&cache)(w, req)
	resp := w.Result()

	if _, exists := cache["deadbeef"]; exists {
		t.Error("Element has not been deleted from repo.")
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Error("Expected status code ", http.StatusNoContent, " got ", resp.StatusCode)
	}
}
