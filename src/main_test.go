package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTestGet(t *testing.T) {
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(TestGet)
	handler.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Wrong status code - got %v, expected %v", status, http.StatusOK)
	}
}
