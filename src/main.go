package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func setupRouter(router *mux.Router) {
	router.Methods("GET").Path("/test").HandlerFunc(testGet)
}

func testGet(w http.ResponseWriter, r *http.Request) {
	log.Println("test get called")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Test"))
}

func main() {
	router := mux.NewRouter()
	setupRouter(router)

	log.Fatal(http.ListenAndServe(":8080", router))
}