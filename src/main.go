package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type FileResponse struct {
	Name    string `json:"name,omitempty"`
	Content string `json:"content,omitempty"`
}

func setupRouter(router *mux.Router) {
	router.Methods("GET").Path("/test").HandlerFunc(GetTest)
	router.Methods("POST").Path("/testFile").HandlerFunc(PostFileTest)
}

func GetTest(w http.ResponseWriter, r *http.Request) {
	log.Println("Test GET called")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Test"))
}

func PostFileTest(w http.ResponseWriter, r *http.Request) {
	log.Println("Test POST file request called")
	// Limit max memory for a single file
	_ = r.ParseMultipartForm(32 << 20)
	var buf bytes.Buffer
	file, header, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	name := strings.Split(header.Filename, ".")
	fmt.Printf("File name %s\n", name[0])
	_, _ = io.Copy(&buf, file)

	content := buf.String()
	fmt.Println(content)
	buf.Reset()
	res := FileResponse{Name: name[0], Content: content}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		panic(err)
	}
}

func main() {
	router := mux.NewRouter()
	setupRouter(router)

	log.Fatal(http.ListenAndServe(":8080", router))
}
