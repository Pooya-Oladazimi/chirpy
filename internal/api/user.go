package api

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type NewUserData struct {
	Email string `json:"email"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *ApiConfig) RegisterUser(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	userData := NewUserData{}
	if err := decoder.Decode(&userData); err != nil {
		writeErrorResponse(w, 400, err.Error())
		return
	}
	newUser, err := cfg.DbQueries.CreateUser(r.Context(), userData.Email)
	if err != nil {
		writeErrorResponse(w, 500, err.Error())
		return
	}
	userDataResp := User{ID: newUser.ID, CreatedAt: newUser.CreatedAt, UpdatedAt: newUser.UpdatedAt, Email: newUser.Email}
	resp, err := json.Marshal(&userDataResp)
	if err != nil {
		writeErrorResponse(w, 500, err.Error())
		return
	}
	writeOkResponse(w, 201, resp)
}
