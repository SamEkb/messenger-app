package main

import (
	"log"
	"net/http"
	"os"
)

const PORT = "8002"

func main() {
	http.HandleFunc("/health/live", liveHandler)
	http.HandleFunc("/health/ready", readyHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = PORT
	}
	addr := ":" + port
	log.Printf("Chat service listening on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Chat service failed: %v", err)
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
