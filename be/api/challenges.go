package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/c0nrad/ctfwriteups/datastore"
	"github.com/c0nrad/ctfwriteups/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

func GetChallengesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctfIDStr := mux.Vars(r)["ctfID"]

	challenges, err := datastore.GetChallengesForCTF(ctx, ctfIDStr)
	if err != nil {
		zap.L().Error("Unable to get challenges", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(challenges)
}

func CreateChallengeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	decoder := json.NewDecoder(r.Body)
	var data models.Challenge
	err := decoder.Decode(&data)
	if err != nil {
		zap.L().Error("Unable to decode challenge", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data.ID = primitive.NewObjectID()
	data.TS = time.Now()

	if data.CTFID.Hex() == "" || data.CTFID.IsZero() {
		zap.L().Error("No CTF ID", zap.Error(err), zap.String("ctfID", data.CTFID.Hex()))
		http.Error(w, "No CTF ID", http.StatusBadRequest)
		return
	}

	err = datastore.SaveChallenge(ctx, data)
	if err != nil {
		zap.L().Error("Unable to save challenge", zap.Error(err), zap.Any("challenge", data))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = datastore.IncrementChallengeCount(ctx, data.CTFID.Hex())
	if err != nil {
		zap.L().Error("Unable to increment challenge count", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(data)
}

func GetChallengeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	challengeID := mux.Vars(r)["challengeID"]

	challenge, err := datastore.GetChallenge(ctx, challengeID)
	if err != nil {
		zap.L().Error("Unable to get challenge", zap.Error(err), zap.String("challengeID", challengeID))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(challenge)
}

func UpdateChallengeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	challengeID := mux.Vars(r)["challengeID"]

	// if err := CanEditChallenge(ctx, r, challengeID); err != nil {
	// 	zap.L().Error("Unable to edit challenge", zap.Error(err), zap.String("challengeID", challengeID))
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	decoder := json.NewDecoder(r.Body)
	var newChallenge models.Challenge
	err := decoder.Decode(&newChallenge)
	if err != nil {
		zap.L().Error("Unable to decode challenge", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	challenge, err := datastore.GetChallenge(ctx, challengeID)
	if err != nil {
		zap.L().Error("Unable to get challenge", zap.Error(err), zap.String("challengeID", challengeID))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newChallenge.ID = challenge.ID
	newChallenge.TS = challenge.TS

	err = datastore.UpdateChallenge(ctx, newChallenge)
	if err != nil {
		zap.L().Error("Unable to update challenge", zap.Error(err), zap.String("challengeID", challengeID))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = datastore.UpdateChallengeTags(ctx, newChallenge.ID.Hex(), newChallenge.Tags)
	if err != nil {
		zap.L().Error("Unable to update challenge tags", zap.Error(err), zap.String("challengeID", challengeID))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(newChallenge)
}

func DeleteChallengeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	challengeID := mux.Vars(r)["challengeID"]

	challenge, err := datastore.GetChallenge(ctx, challengeID)
	if err != nil {
		zap.L().Error("Unable to get challenge", zap.Error(err), zap.String("challengeID", challengeID))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = datastore.DeleteChallenge(ctx, challengeID)
	if err != nil {
		zap.L().Error("Unable to delete challenge", zap.Error(err), zap.String("challengeID", challengeID))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = datastore.DecrementChallengeCount(ctx, challenge.CTFID.Hex())
	if err != nil {
		zap.L().Error("Unable to decrement challenge count", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode("ok")
}
