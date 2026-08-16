package main

import (
	"database/sql"
	"fmt"
	"github.com/Pooya-Oladazimi/chirpy.git/internal/api"
	"github.com/Pooya-Oladazimi/chirpy.git/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	PORT         = 8080
	PATH_PREFIX  = "/app/"
	API_PREFIX   = "/api/"
	ADMIN_PREFIX = "/admin/"
	GET_METHOD   = "GET "
	POST_METHOD  = "POST "
)

func main() {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dbQueries := database.New(db)
	serverRouter := http.NewServeMux()
	cfg := api.ApiConfig{}
	cfg.FileserverHits.Store(0)
	cfg.Platform = os.Getenv("PLATFORM")
	cfg.DbQueries = dbQueries
	cfg.JWT_SECRET = os.Getenv("JWT_SECRET")
	fs := http.FileServer(http.Dir("."))
	serverRouter.Handle(PATH_PREFIX, http.StripPrefix(PATH_PREFIX, cfg.MiddlewareMetricsInc(fs)))
	serverRouter.HandleFunc(GET_METHOD+API_PREFIX+"healthz/", cfg.CheckHealth)
	serverRouter.HandleFunc(GET_METHOD+ADMIN_PREFIX+"metrics/", cfg.GetHits)
	serverRouter.HandleFunc(POST_METHOD+ADMIN_PREFIX+"reset/", cfg.Reset)
	serverRouter.HandleFunc(POST_METHOD+API_PREFIX+"users/", cfg.RegisterUser)
	serverRouter.Handle(POST_METHOD+API_PREFIX+"chirps/", cfg.MiddlewareAuth(http.HandlerFunc(cfg.CreateChirp)))
	serverRouter.HandleFunc(GET_METHOD+API_PREFIX+"chirps/", cfg.GetAllChirps)
	serverRouter.HandleFunc(GET_METHOD+API_PREFIX+"chirps/{chirpID}", cfg.GetChirp)
	serverRouter.HandleFunc(POST_METHOD+API_PREFIX+"login/", cfg.Login)
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
