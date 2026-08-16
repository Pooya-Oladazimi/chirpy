package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Pooya-Oladazimi/chirpy.git/internal/auth"
	"github.com/Pooya-Oladazimi/chirpy.git/internal/database"
	"github.com/google/uuid"
)

type newChirpData struct {
	Body string `json:"body"`
}
type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *ApiConfig) CreateChirp(w http.ResponseWriter, r *http.Request) {
	userChirpData, err := loadNewChirpData(r)
	if err != nil && strings.Contains(strings.ToLower(err.Error()), "invalid uuid") {
		writeErrorResponse(w, 401, "You need to register to use this endpoint.")
		return
	}
	if err != nil {
		writeErrorResponse(w, 500, "Something went wrong.")
		return
	}
	cleanChirpBody, err := validateChirp(userChirpData.Body)
	if err != nil {
		writeErrorResponse(w, 500, err.Error())
		return
	}
	userToken, _ := auth.GetBearerToken(r.Header)
	userId, _ := auth.ValidateJWT(userToken, cfg.JWT_SECRET)
	user, err := cfg.DbQueries.GetUser(r.Context(), userId)
	if err != nil {
		writeErrorResponse(w, 401, "You need to register to use this endpoint.")
		return
	}
	chirpInsertData := database.CreateChirpParams{Body: cleanChirpBody, UserID: user.ID}
	newChirp, err := cfg.DbQueries.CreateChirp(r.Context(), chirpInsertData)
	if err != nil {
		writeErrorResponse(w, 500, "Something went wrong.")
		return
	}
	newChirpResp := Chirp{ID: newChirp.ID, CreatedAt: newChirp.CreatedAt, UpdatedAt: newChirp.UpdatedAt, Body: newChirp.Body, UserID: newChirp.UserID}
	resp, err := json.Marshal(&newChirpResp)
	if err != nil {
		writeErrorResponse(w, 500, "Something went wrong.")
		return
	}
	writeOkResponse(w, 201, resp)
}

func (cfg *ApiConfig) GetChirp(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("chirpID")
	chirpUuid, err := uuid.Parse(id)
	if err != nil {
		writeErrorResponse(w, 404, "Chirp does not exist.")
		return
	}
	dbChirp, err := cfg.DbQueries.GetChirp(r.Context(), chirpUuid)
	if err != nil {
		writeErrorResponse(w, 404, "Chirp does not exist.")
		return
	}
	chirp := Chirp{ID: dbChirp.ID, CreatedAt: dbChirp.CreatedAt, UpdatedAt: dbChirp.UpdatedAt, Body: dbChirp.Body, UserID: dbChirp.UserID}
	chirpResp, err := json.Marshal(&chirp)
	if err != nil {
		writeErrorResponse(w, 500, "Something went wrong.")
		return
	}
	writeOkResponse(w, 200, chirpResp)
}

func (cfg *ApiConfig) GetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.DbQueries.GetAllChirps(r.Context())
	if err != nil {
		writeErrorResponse(w, 500, "Something went wrong.")
		return
	}
	jsonChirps := make([]Chirp, 0)
	for _, chirp := range chirps {
		jsonChirps = append(jsonChirps, Chirp{ID: chirp.ID, Body: chirp.Body, CreatedAt: chirp.CreatedAt, UpdatedAt: chirp.UpdatedAt, UserID: chirp.UserID})
	}
	allChirpsResp, err := json.Marshal(jsonChirps)
	if err != nil {
		writeErrorResponse(w, 500, "Something went wrong.")
		return
	}
	writeOkResponse(w, 200, allChirpsResp)
}

func validateChirp(chirpBody string) (string, error) {
	if len(chirpBody) > 140 {
		return "", fmt.Errorf("Chirp is too long.")
	}
	return cleanChirp(chirpBody), nil

}

func loadNewChirpData(r *http.Request) (newChirpData, error) {
	decoder := json.NewDecoder(r.Body)
	postData := newChirpData{}
	if err := decoder.Decode(&postData); err != nil {
		return newChirpData{}, err
	}
	return postData, nil
}

func cleanChirp(chirp string) string {
	FORBIDEN_WORDS := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(chirp, " ")
	for i := range words {
		for _, fw := range FORBIDEN_WORDS {
			if strings.EqualFold(words[i], fw) {
				words[i] = "****"
				break
			}
		}
	}
	return strings.Join(words, " ")
}
