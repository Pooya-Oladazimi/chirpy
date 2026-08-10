package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

const (
	PORT                           = 8080
	PATH_PREFIX                    = "/app/"
	API_PREFIX                     = "/api/"
	ADMIN_PREFIX                   = "/admin/"
	TEXT_CONTENT_TYPE_HEADER_VALUE = "text/plain; charset=utf-8"
	HTML_CONTENT_TYPE_HEADER_VALUE = "text/html; charset=utf-8"
	JSON_CONTENT_TYPE_HEADER_VALUE = "application/json"
	GET_METHOD                     = "GET "
	POST_METHOD                    = "POST "
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	serverRouter := http.NewServeMux()
	cfg := apiConfig{}
	cfg.fileserverHits.Store(0)
	fs := http.FileServer(http.Dir("."))
	serverRouter.Handle(PATH_PREFIX, http.StripPrefix(PATH_PREFIX, cfg.middlewareMetricsInc(fs)))
	serverRouter.HandleFunc(GET_METHOD+API_PREFIX+"healthz/", checkHealth)
	serverRouter.HandleFunc(GET_METHOD+ADMIN_PREFIX+"metrics/", cfg.getHits)
	serverRouter.HandleFunc(POST_METHOD+ADMIN_PREFIX+"reset/", cfg.reset)
	serverRouter.HandleFunc(POST_METHOD+API_PREFIX+"validate_chirp/", validateChirp)
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
	w.Header().Add("Content-Type", TEXT_CONTENT_TYPE_HEADER_VALUE)
	w.WriteHeader(200)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		w.WriteHeader(500)
	}
}

func (cfg *apiConfig) getHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", HTML_CONTENT_TYPE_HEADER_VALUE)
	w.WriteHeader(200)
	buf := make([]byte, 0)
	buf = fmt.Appendf(buf, `
	<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
	</html>
	`, cfg.fileserverHits.Load())
	_, err := w.Write(buf)
	if err != nil {
		w.WriteHeader(500)
	}
}

func (cfg *apiConfig) reset(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", TEXT_CONTENT_TYPE_HEADER_VALUE)
	w.WriteHeader(200)
	cfg.fileserverHits.Store(0)
	buf := make([]byte, 0)
	buf = fmt.Appendf(buf, "Hits reseted: %d\n", cfg.fileserverHits.Load())
	_, err := w.Write(buf)
	if err != nil {
		w.WriteHeader(500)
	}
}

func validateChirp(w http.ResponseWriter, r *http.Request) {
	type data struct {
		Body string `json:"body"`
	}
	type validResp struct {
		Value bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	postData := data{}
	w.Header().Add("Content-Type", JSON_CONTENT_TYPE_HEADER_VALUE)
	if err := decoder.Decode(&postData); err != nil {
		writeErrorResponse(w, 500, "something went wrong")
		return
	}
	if len(postData.Body) > 140 {
		writeErrorResponse(w, 400, "Chirp is too long.")
		return
	}
	w.WriteHeader(200)
	validRespJson := validResp{Value: true}
	validRespJsonData, err := json.Marshal(validRespJson)
	if err == nil {
		w.Write(validRespJsonData)
		return
	}
	w.Write([]byte("something went wrong"))
}

func writeErrorResponse(w http.ResponseWriter, errorCode int, errorMsg string) {
	type errorResp struct {
		ErrorMsg string `json:"error"`
	}
	w.WriteHeader(errorCode)
	errorRespJson := errorResp{ErrorMsg: errorMsg}
	errorRespJsonData, err := json.Marshal(errorRespJson)
	if err == nil {
		w.Write(errorRespJsonData)
		return
	}
	w.Write([]byte("something went wrong"))
}
