package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/shivtriv12/chirpy/internal/auth"
	"github.com/shivtriv12/chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
		Refresh_Token string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid Email or Password", err)
		return
	}
	err = auth.CheckPasswordsWithHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid Email or Password", err)
		return
	}

	jwtTime := time.Hour
	token, err := auth.MakeJWT(user.ID, cfg.JWT_Secret, jwtTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Not able to Create JWT", err)
		return
	}
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "not able to create refresh token", err)
	}
	err = cfg.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "not able to store refresh token", err)
	}
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			Token:       token,
			IsChirpyRed: user.IsChirpyRed,
		},
		Refresh_Token: refreshToken,
	})
}
