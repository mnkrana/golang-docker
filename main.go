package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", index)
	http.ListenAndServe(":80", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	log.Println("request " + r.URL.Path)
	io.WriteString(w, "Hello, golang docker test!")
}
