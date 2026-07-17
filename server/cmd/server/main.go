package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	port := "8080"
	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(":" + port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
