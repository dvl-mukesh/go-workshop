package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world")
	})
	fmt.Println("Starting server")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
