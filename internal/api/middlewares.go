package api

import (
	"net/http"

	"github.com/Pooya-Oladazimi/chirpy.git/internal/auth"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *ApiConfig) MiddlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userToken, err := auth.GetBearerToken(r.Header)
		if err != nil {
			writeErrorResponse(w, 401, "Unauthorized")
			return
		}
		userId, err := auth.ValidateJWT(userToken, cfg.JWT_SECRET)
		if err != nil || userId == uuid.Nil {
			writeErrorResponse(w, 401, "Unauthorized")
			return
		}
		next.ServeHTTP(w, r)
	})
}
