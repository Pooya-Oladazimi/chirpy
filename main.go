package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	PORT        = 8080
	PATH_PREFIX = "/app/"
)

func main() {
	serverRouter := http.NewServeMux()
	fs := http.FileServer(http.Dir("."))
	serverRouter.Handle(PATH_PREFIX, http.StripPrefix(PATH_PREFIX, fs))
	serverRouter.HandleFunc("/healthz/", checkHealth)
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", PORT),
		Handler:        serverRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Printf("Server is running on %d ...\n", PORT)
	log.Fatal(server.ListenAndServe())
}

func checkHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		w.WriteHeader(500)
	}
}
