package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/shivtriv12/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerDeleteChirpById(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "auth error", err)
		return
	}

	userId, err := auth.ValidateJWT(accessToken, cfg.JWT_Secret)
	if err != nil {
		respondWithError(w, 403, "incorrect accessToken", err)
		return
	}

	chirpId := r.PathValue("chirpid")
	if chirpId == "" {
		respondWithError(w, http.StatusBadRequest, "no chirp id provided", nil)
		return
	}

	chirpIdUUID, err := uuid.Parse(chirpId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to convert chirp string to uuid", err)
		return
	}

	dbChirp, err := cfg.dbQueries.GetChirp(r.Context(), chirpIdUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}
	if dbChirp.UserID != userId {
		respondWithError(w, http.StatusForbidden, "You can't delete this chirp", err)
		return
	}

	err = cfg.dbQueries.DeleteChirpById(r.Context(), dbChirp.ID)
	if err != nil {
		respondWithError(w, 404, "chirp not found", err)
		return
	}
	w.WriteHeader(204)
}
