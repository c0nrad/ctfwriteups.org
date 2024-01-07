package api

import (
	"encoding/json"
	"net/http"

	"github.com/c0nrad/ctfwriteups/datastore"
	"github.com/c0nrad/ctfwriteups/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

func GetSeensForUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := GetUser(ctx, r)
	if err != nil {
		http.Error(w, "not logged in", http.StatusUnauthorized)
		return
	}

	seens, err := datastore.GetSeensForUser(ctx, user.ID.Hex())
	if err != nil {
		zap.L().Error("Unable to get seens", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(seens)
}

func IncrementSeenHandler(w http.ResponseWriter, r *http.Request) {
	// allow anonymous seen

	ctx := r.Context()

	user, err := GetUser(ctx, r)
	var userID primitive.ObjectID = primitive.NilObjectID
	if err == nil {
		userID = user.ID
	}

	decoder := json.NewDecoder(r.Body)
	var seen models.Seen
	err = decoder.Decode(&seen)

	if err != nil {
		zap.L().Error("Unable to decode seen", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	seen.ID = primitive.NewObjectID()
	seen.UserID = userID

	if seen.WriteupID.Hex() == "" || seen.WriteupID.IsZero() {
		zap.L().Error("No writeup ID", zap.Error(err), zap.String("writeupID", seen.WriteupID.Hex()))
		http.Error(w, "No writeup ID", http.StatusBadRequest)
		return
	}

	err = datastore.IncrementSeen(ctx, seen)
	if err != nil {
		zap.L().Error("Unable to create seen", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(seen)
}

func DeleteSeenHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := GetUser(ctx, r)
	if err != nil {
		http.Error(w, "not logged in", http.StatusUnauthorized)
		return
	}

	seenID := mux.Vars(r)["seenID"]

	seen, err := datastore.GetSeen(ctx, seenID)
	if err != nil {
		zap.L().Error("Unable to get seen", zap.Error(err), zap.String("seenID", seenID))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if seen.UserID != user.ID {
		zap.L().Error("Unable to delete seen", zap.Error(err), zap.String("seenID", seenID))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = datastore.DeleteSeen(ctx, seenID)
	if err != nil {
		zap.L().Error("Unable to delete seen", zap.Error(err), zap.String("seenID", seenID))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(seen)
}
