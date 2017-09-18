package main

import (
	"fmt"
	"log"
	"net/http"
)

const version = "1.0.0"

func main() {
	println("Starting up....")
	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s\n", version)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "rightstuff inv api")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
