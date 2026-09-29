package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "PasteBox")
	})

	log.Printf("PasteBox listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
