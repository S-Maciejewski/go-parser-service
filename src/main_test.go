package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

type Response struct {
	Name    string `json:"name"`
	Content string `json:"content"`
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

	var resBody Response
	err = json.Unmarshal(recorder.Body.Bytes(), &resBody)
	if err != nil {
		t.Error(err)
	}
	checksum := md5.Sum([]byte(resBody.Content))
	if !(resBody.Name == "test" && hex.EncodeToString(checksum[:]) == "343d5e9e9766c9c617fbf3a7c7c64779") {
		t.Error("Name of the file and / or content checksum incorrect")
	}
}
