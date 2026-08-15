package api

import (
	"fmt"
	"net/http"
)

func (cfg *ApiConfig) CheckHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", TEXT_CONTENT_TYPE_HEADER_VALUE)
	w.WriteHeader(200)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		w.WriteHeader(500)
	}
}

func (cfg *ApiConfig) GetHits(w http.ResponseWriter, r *http.Request) {
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
	`, cfg.FileserverHits.Load())
	_, err := w.Write(buf)
	if err != nil {
		w.WriteHeader(500)
	}
}

func (cfg *ApiConfig) Reset(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" {
		w.WriteHeader(401)
		return
	}
	w.Header().Add("Content-Type", TEXT_CONTENT_TYPE_HEADER_VALUE)
	w.WriteHeader(200)
	cfg.FileserverHits.Store(0)
	buf := make([]byte, 0)
	buf = fmt.Appendf(buf, "Hits reseted: %d\n", cfg.FileserverHits.Load())
	_, err := w.Write(buf)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	if err := cfg.DbQueries.DeleteAllUsers(r.Context()); err != nil {
		w.WriteHeader(500)
		return
	}

}
