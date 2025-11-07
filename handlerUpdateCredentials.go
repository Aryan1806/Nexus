package main

import (
	"encoding/json"
	"net/http"

	"github.com/shivtriv12/chirpy/internal/auth"
	"github.com/shivtriv12/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateCredentials(w http.ResponseWriter, r *http.Request) {
	type paramaters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "auth error", err)
		return
	}

	userId, err := auth.ValidateJWT(accessToken, cfg.JWT_Secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect accessToken", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := paramaters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "not able to access body", err)
		return
	}

	hashedPassword, err := auth.HashPasswords(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "not able to hash password", err)
		return
	}

	user, err := cfg.dbQueries.UpdateUserByUserId(r.Context(), database.UpdateUserByUserIdParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
		ID:             userId,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to get user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}
