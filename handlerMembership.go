package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/shivtriv12/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerMembership(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Event string `json:"event"`
		Data  struct {
			User_ID string `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find api key", err)
		return
	}
	if apiKey != cfg.POLKA_KEY {
		respondWithError(w, http.StatusUnauthorized, "API key is invalid", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameter{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding body", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userId, err := uuid.Parse(params.Data.User_ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to convert userid string to uuid", err)
		return
	}

	err = cfg.dbQueries.UpdateMembershipById(r.Context(), userId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "unable to update membership", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
