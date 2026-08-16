package api

import (
	"github.com/Pooya-Oladazimi/chirpy.git/internal/database"
	"sync/atomic"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	Platform       string
	DbQueries      *database.Queries
	JWT_SECRET     string
}

const (
	TEXT_CONTENT_TYPE_HEADER_VALUE = "text/plain; charset=utf-8"
	HTML_CONTENT_TYPE_HEADER_VALUE = "text/html; charset=utf-8"
	JSON_CONTENT_TYPE_HEADER_VALUE = "application/json"
	ONE_HOUR_IN_SECONDS            = 3600
)
