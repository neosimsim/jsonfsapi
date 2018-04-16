package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreate(t *testing.T) {
	cache := Cache{}

	req := httptest.NewRequest("POST", "http://localhost", strings.NewReader(`{ "title": "I'm new"}`))
	w := httptest.NewRecorder()
	ServeElementsFactory(&cache)(w, req)

	resp := w.Result()

	location := w.Header().Get("Location")
	if location == "" {
		t.Fatal("Location header is missing")
	}

	uuid := location
	if _, exists := cache[uuid]; !exists {
		t.Error("Element has not been stored to repo.")
	}

	if resp.StatusCode != http.StatusCreated {
		t.Error("Expected status code ", http.StatusCreated, " got ", resp.StatusCode)
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
