package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Pooya-Oladazimi/chirpy.git/internal/auth"
	"github.com/Pooya-Oladazimi/chirpy.git/internal/database"
	"github.com/google/uuid"
)

type NewUserData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginData struct {
	NewUserData
	Expires_in_seconds int `json:"expires_in_seconds"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type LoginResp struct {
	User
	Token string `json:"token"`
}

func (cfg *ApiConfig) RegisterUser(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	userData := NewUserData{}
	if err := decoder.Decode(&userData); err != nil {
		writeErrorResponse(w, 400, err.Error())
		return
	}
	hashedPass, err := auth.HashPassword(userData.Password)
	if err != nil {
		writeErrorResponse(w, 500, "Something went wrong.")
		return
	}
	newUserDbParams := database.CreateUserParams{Email: userData.Email, HashedPassword: hashedPass}
	newUser, err := cfg.DbQueries.CreateUser(r.Context(), newUserDbParams)
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

func (cfg *ApiConfig) Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	userData := LoginData{}
	if err := decoder.Decode(&userData); err != nil {
		writeErrorResponse(w, 400, err.Error())
		return
	}
	if userData.Expires_in_seconds == 0 || userData.Expires_in_seconds > ONE_HOUR_IN_SECONDS {
		userData.Expires_in_seconds = ONE_HOUR_IN_SECONDS
	}
	user, err := cfg.DbQueries.GetUserByEmail(r.Context(), userData.Email)
	if err != nil {
		writeErrorResponse(w, 401, "Unauthorized")
		return
	}
	ok, err := auth.CheckPasswordHash(userData.Password, user.HashedPassword)
	if err != nil || !ok {
		writeErrorResponse(w, 401, "Unauthorized")
		return
	}
	expiresIn := time.Duration(userData.Expires_in_seconds) * time.Second
	singedToken, err := auth.MakeJWT(user.ID, cfg.JWT_SECRET, expiresIn)
	if err != nil {
		writeErrorResponse(w, 500, "Something went wrong.")
		return
	}
	userDataResp := LoginResp{
		User:  User{ID: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email},
		Token: singedToken,
	}

	resp, err := json.Marshal(&userDataResp)
	if err != nil {
		writeErrorResponse(w, 500, err.Error())
		return
	}
	writeOkResponse(w, 200, resp)

}
