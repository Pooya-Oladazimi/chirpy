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
}

type User struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

type LoginResp struct {
	User
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
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
	userDataResp := User{ID: newUser.ID, CreatedAt: newUser.CreatedAt, UpdatedAt: newUser.UpdatedAt, Email: newUser.Email, IsChirpyRed: newUser.IsChirpyRed}
	resp, err := json.Marshal(&userDataResp)
	if err != nil {
		writeErrorResponse(w, 500, err.Error())
		return
	}
	writeOkResponse(w, 201, resp)
}

func (cfg *ApiConfig) UpdateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	userData := NewUserData{}
	if err := decoder.Decode(&userData); err != nil {
		writeErrorResponse(w, 500, "Something went wrong.")
		return
	}
	hashed_pass, err := auth.HashPassword(userData.Password)
	if err != nil {
		writeErrorResponse(w, 500, "Something went wrong.")
		return
	}
	userId, ok := r.Context().Value("userIdInToken").(uuid.UUID)
	if !ok {
		writeErrorResponse(w, 500, "Something went wrong.")
		return
	}
	updateUserDbData := database.UpdateUserParams{UserID: userId, NewEmail: userData.Email, Password: hashed_pass}
	updatedUser, err := cfg.DbQueries.UpdateUser(r.Context(), updateUserDbData)
	if err != nil {
		writeErrorResponse(w, 500, "Something went wrong.")
		return
	}
	userResp := User{ID: updatedUser.ID, UpdatedAt: updatedUser.UpdatedAt, CreatedAt: updatedUser.CreatedAt, Email: updatedUser.Email, IsChirpyRed: updatedUser.IsChirpyRed}
	userRespBody, err := json.Marshal(&userResp)
	if err != nil {
		writeErrorResponse(w, 500, "Something went wrong.")
		return
	}
	w.WriteHeader(200)
	w.Write(userRespBody)
}

func (cfg *ApiConfig) Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	userData := LoginData{}
	if err := decoder.Decode(&userData); err != nil {
		writeErrorResponse(w, 400, err.Error())
		return
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
	expiresIn := time.Duration(ONE_HOUR_IN_SECONDS) * time.Second
	singedToken, err := auth.MakeJWT(user.ID, cfg.JWT_SECRET, expiresIn)
	if err != nil {
		writeErrorResponse(w, 500, "Something went wrong.")
		return
	}
	refreshTokenDbParams := database.CreateRefreshTokenParams{
		Token:     auth.MakeRefreshToken(),
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add(time.Duration(ONE_DAY_IN_SECONDS) * time.Second),
	}
	refreshToken, err := cfg.DbQueries.CreateRefreshToken(r.Context(), refreshTokenDbParams)
	if err != nil {
		writeErrorResponse(w, 500, "Something went wrong.")
		return
	}
	userDataResp := LoginResp{
		User:         User{ID: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email, IsChirpyRed: user.IsChirpyRed},
		Token:        singedToken,
		RefreshToken: refreshToken.Token,
	}

	resp, err := json.Marshal(&userDataResp)
	if err != nil {
		writeErrorResponse(w, 500, err.Error())
		return
	}
	writeOkResponse(w, 200, resp)
}

func (cfg *ApiConfig) RefreshUserToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		writeErrorResponse(w, 401, "Unauthorized")
		return
	}
	tokenInDb, err := cfg.DbQueries.GetRefreshToken(r.Context(), token)
	if err != nil {
		writeErrorResponse(w, 401, "Unauthorized")
		return
	}
	if time.Now().UTC().After(tokenInDb.ExpiresAt) || tokenInDb.RevokedAt.Valid {
		writeErrorResponse(w, 401, "Unauthorized")
		return
	}
	expiresIn := time.Duration(ONE_HOUR_IN_SECONDS) * time.Second
	singedToken, err := auth.MakeJWT(tokenInDb.UserID, cfg.JWT_SECRET, expiresIn)
	type tokenResp struct {
		Token string `json:"token"`
	}
	tokenResponse := tokenResp{Token: singedToken}
	tokenRespBody, err := json.Marshal(&tokenResponse)
	if err != nil {
		writeErrorResponse(w, 500, "Something went wrong.")
		return
	}
	w.WriteHeader(200)
	w.Write(tokenRespBody)
}

func (cfg *ApiConfig) RevokeUserToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		writeErrorResponse(w, 401, "Unauthorized")
		return
	}
	if err := cfg.DbQueries.RevokeRefreshToken(r.Context(), token); err != nil {
		writeErrorResponse(w, 500, "Something went wrong.")
		return
	}
	w.WriteHeader(204)
}
