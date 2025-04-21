package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/health/live", liveHandler)
	http.HandleFunc("/health/ready", readyHandler)

	addr := ":8001"
	log.Printf("Users service listening on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Users service failed: %v", err)
	}
}

func liveHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ready"))
}
