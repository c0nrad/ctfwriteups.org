package api

import (
	"encoding/json"
	"net/http"
)

func GetMeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := GetUser(ctx, r)
	if err != nil {
		http.Error(w, "not logged in", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(user)
}
