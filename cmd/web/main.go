package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	const addr = ":8080"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "DGSIS Website")
	})

	log.Printf("server listening on %s", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
