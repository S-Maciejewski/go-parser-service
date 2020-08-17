package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

type Form struct {
	Value map[string][]string
	File  map[string][]*multipart.FileHeader
}

func TestGetTest(t *testing.T) {
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(GetTest)
	handler.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Wrong status code - got %v, expected %v", status, http.StatusOK)
	}
}

func TestPostFileTest(t *testing.T) {
	testFilePath := "../testData/test.txt"
	file, err := os.Open(testFilePath)
	if err != nil {
		t.Error(err)
	}

	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(testFilePath))
	if err != nil {
		_ = writer.Close()
		t.Error(err)
	}
	_, _ = io.Copy(part, file)
	_ = writer.Close()

	req := httptest.NewRequest("POST", "/testFile", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	recorder := httptest.NewRecorder()

	PostFileTest(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Wrong status code - got %v, expected %v", status, http.StatusOK)
	}

	t.Log(recorder.Body.String())
}
