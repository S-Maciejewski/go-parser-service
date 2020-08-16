package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func setupRouter(router *mux.Router) {
	router.Methods("GET").Path("/test").HandlerFunc(testGet)
	router.Methods("POST").Path("/testFile").HandlerFunc(testPostFile)
}

func testGet(w http.ResponseWriter, r *http.Request) {
	log.Println("Test GET called")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Test"))
}

func testPostFile(w http.ResponseWriter, r *http.Request) {
	log.Println("Test POST file request called")
	// Limit max memory for a single file
	r.ParseMultipartForm(32 << 20)
	var buf bytes.Buffer
	file, header, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	name := strings.Split(header.Filename, ".")
	fmt.Printf("File name %s\n", name[0])
	io.Copy(&buf, file)

	content := buf.String()
	fmt.Println(content)
	buf.Reset()
	return
}

func main() {
	router := mux.NewRouter()
	setupRouter(router)

	log.Fatal(http.ListenAndServe(":8080", router))
}
