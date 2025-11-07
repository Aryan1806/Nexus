package main

import (
	"net/http"
	"time"

	"github.com/shivtriv12/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerPostRefreshToken(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "not able to extract token", err)
		return
	}

	user_id, err := cfg.dbQueries.GetUserByRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "token invalid or expired", err)
		return
	}

	jwt, err := auth.MakeJWT(user_id, cfg.JWT_Secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to generate jwt", err)
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		Token: jwt,
	})
}
