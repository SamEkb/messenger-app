package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const PORT = "8001"

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/health/live", liveHandler)
	r.Get("/health/ready", readyHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = PORT
	}
	addr := ":" + port
	log.Printf("Auth service listening on %s\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Auth service failed: %v", err)
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
