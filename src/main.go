package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tealeg/xlsx"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type FileResponse struct {
	Name    string `json:"name,omitempty"`
	Content string `json:"content,omitempty"`
}

type JsonNode struct {
	ID       int        `json:"id"`
	Name     string     `json:"name"`
	Children []JsonNode `json:"children"`
}

func setupRouter(router *mux.Router) {
	router.Methods("GET").Path("/test").HandlerFunc(GetTest)
	router.Methods("POST").Path("/testFile").HandlerFunc(PostFileTest)
	router.Methods("POST").Path("/parse").HandlerFunc(XlsxToJSON)
}

func GetTest(w http.ResponseWriter, r *http.Request) {
	log.Println("Test GET called")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Test"))
}

func GetFileFromRequest(r *http.Request) (string, multipart.File) {
	_ = r.ParseMultipartForm(32 << 20) // Max size limit for a file
	file, header, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	name := strings.Split(header.Filename, ".")
	return name[0], file
}

func PostFileTest(w http.ResponseWriter, r *http.Request) {
	log.Println("Test POST file request called")
	var buf bytes.Buffer
	name, file := GetFileFromRequest(r)
	defer file.Close()
	fmt.Printf("File name %s\n", name)
	_, _ = io.Copy(&buf, file)

	content := buf.String()
	fmt.Println(content)
	buf.Reset()
	res := FileResponse{Name: name, Content: content}
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		panic(err)
	}
}

func XlsxToJSON(w http.ResponseWriter, r *http.Request) {
	name, file := GetFileFromRequest(r)
	log.Printf("Received .xlsx file to parse - %s\n", name)
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, file)
	wb, err := xlsx.OpenBinary(buf.Bytes())
	if err != nil {
		panic(err)
	}

	sheets := wb.Sheets
	sh, _ := wb.Sheet[sheets[0].Name]

	for i := 0; i < sh.MaxRow; i++ {
		fmt.Println(sh.Row(i).Cells)
		// TODO
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	router := mux.NewRouter()
	setupRouter(router)

	log.Fatal(http.ListenAndServe(":8080", router))
}
