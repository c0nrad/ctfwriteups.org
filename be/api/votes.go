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

func GetVotesForUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := GetUser(ctx, r)
	if err != nil {
		http.Error(w, "not logged in", http.StatusUnauthorized)
		return
	}

	votes, err := datastore.GetVotesForUser(ctx, user.ID.Hex())
	if err != nil {
		zap.L().Error("Unable to get votes", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(votes)
}

func CreateVoteForWriteupHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := GetUser(ctx, r)
	if err != nil {
		http.Error(w, "not logged in", http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var vote models.Vote
	err = decoder.Decode(&vote)

	if err != nil {
		zap.L().Error("Unable to decode vote", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vote.ID = primitive.NewObjectID()
	vote.UserID = user.ID
	vote.IsUpvote = true

	if vote.WriteupID.Hex() == "" || vote.WriteupID.IsZero() {
		zap.L().Error("No writeup ID", zap.Error(err), zap.String("writeupID", vote.WriteupID.Hex()))
		http.Error(w, "No writeup ID", http.StatusBadRequest)
		return
	}

	err = datastore.SaveVote(ctx, vote)
	if err != nil {
		zap.L().Error("Unable to create vote", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// incrememt vote count on writeup
	err = datastore.IncrementVoteCount(ctx, vote.WriteupID.Hex())
	if err != nil {
		zap.L().Error("Unable to increment vote count", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(vote)
}

func DeleteVoteForWriteupHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := GetUser(ctx, r)
	if err != nil {
		http.Error(w, "not logged in", http.StatusUnauthorized)
		return
	}

	voteID := mux.Vars(r)["voteID"]

	vote, err := datastore.GetVote(ctx, voteID)
	if err != nil {
		zap.L().Error("Unable to get vote", zap.Error(err), zap.String("voteID", voteID))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if vote.UserID != user.ID {
		zap.L().Error("Unable to delete vote", zap.Error(err), zap.String("voteID", voteID))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = datastore.DeleteVote(ctx, voteID)
	if err != nil {
		zap.L().Error("Unable to delete vote", zap.Error(err), zap.String("voteID", voteID))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// decrememt vote count on writeup
	err = datastore.DecrementVoteCount(ctx, vote.WriteupID.Hex())
	if err != nil {
		zap.L().Error("Unable to decrement vote count", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(vote)
}
