package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/micropub", handleMicropub)
	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleMicropub(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement Micropub endpoint
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Micropub endpoint not yet implemented"))
}
