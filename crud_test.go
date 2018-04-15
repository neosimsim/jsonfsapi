package main

import (
	"net/http/httptest"
	"strings"
	"testing"
	"log"
	"os"
)

func TestCreate(t *testing.T) {
	devNull, _ := os.Create(os.DevNull)
	log.SetOutput(devNull)
	req := httptest.NewRequest("POST", "http://localhost?uuid=deadbeef", strings.NewReader("{}"))
	w := httptest.NewRecorder()
	ServeElementsFactory(NewFileRepo())(w, req)

	resp := w.Result()

	if resp.StatusCode != 201 {
		t.Error("Expected status code ", 201, " got ", resp.StatusCode)
	}
}
