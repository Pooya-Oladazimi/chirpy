package api

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
)

type UserWebHooks struct {
	Event string `json:"event"`
	Data  struct {
		UserId uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func (cfg *ApiConfig) PolkaWebHooks(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	userWebhookData := UserWebHooks{}
	if err := decoder.Decode(&userWebhookData); err != nil {
		writeErrorResponse(w, 204, "")
		return
	}
	if userWebhookData.Event != "user.upgraded" {
		writeErrorResponse(w, 204, "")
		return
	}
	_, err := cfg.DbQueries.UpgradeUserToChirpyRed(r.Context(), userWebhookData.Data.UserId)
	if err != nil {
		writeErrorResponse(w, 404, "User not found.")
		return
	}
	writeOkResponse(w, 204, []byte(""))
}
