package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/c0nrad/ctfwriteups/datastore"
	"github.com/c0nrad/ctfwriteups/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

func GetWriteupsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	writeups, err := datastore.GetWriteups(ctx)
	if err != nil {
		zap.L().Error("Unable to get writeups", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(writeups)
}

func CreateWriteupHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := GetUser(ctx, r)
	if err != nil {
		http.Error(w, "not logged in", http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var data models.Writeup
	err = decoder.Decode(&data)
	if err != nil {
		zap.L().Error("Unable to decode writeup", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data.ID = primitive.NewObjectID()
	data.TS = time.Now()
	data.SubmitterID = user.ID

	if !IsValidURL(data.URL) {
		zap.L().Error("Invalid URL", zap.Error(err), zap.String("url", data.URL))
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	if data.SubmitterID.Hex() == "" || data.SubmitterID.IsZero() {
		zap.L().Error("No submitter ID", zap.Error(err), zap.String("submitterID", data.SubmitterID.Hex()))
		http.Error(w, "No submitter ID", http.StatusBadRequest)
		return
	}

	ctf, err := datastore.GetCTF(ctx, data.CTFID.Hex())
	if err != nil {
		zap.L().Error("No CTF ID", zap.Error(err), zap.String("ctfID", data.CTFID.Hex()))
		http.Error(w, "No CTF ID", http.StatusBadRequest)
		return
	}
	data.CTFName = ctf.Name
	data.CTFEndDate = ctf.EndDate

	challenge, err := datastore.GetChallenge(ctx, data.ChallengeID.Hex())
	if err != nil {
		zap.L().Error("No challenge ID", zap.Error(err), zap.String("challengeID", data.ChallengeID.Hex()))
		http.Error(w, "No challenge ID", http.StatusBadRequest)
		return
	}
	data.ChallengeName = challenge.Name
	data.ChallengeCategory = challenge.Category
	data.Tags = challenge.Tags

	err = datastore.SaveWriteup(ctx, data)
	if err != nil {
		zap.L().Error("Unable to save writeup", zap.Error(err), zap.Any("writeup", data))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = datastore.IncrementWriteupCount(ctx, data.CTFID.Hex())
	if err != nil {
		zap.L().Error("Unable to increment writeup count", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(data)
}

func IsValidURL(urlStr string) bool {
	u, err := url.Parse(urlStr)
	if err != nil {
		return false
	}
	if u.Scheme != "https" && u.Scheme != "http" {
		return false
	}

	return true
}

func GetWriteupHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	writeupID := mux.Vars(r)["writeupID"]

	writeup, err := datastore.GetWriteup(ctx, writeupID)
	if err != nil {
		zap.L().Error("Unable to get writeup", zap.Error(err), zap.String("writeupID", writeupID))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(writeup)
}

func UpdateWriteupHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	writeupID := mux.Vars(r)["writeupID"]

	// if err := CanEditWriteup(ctx, r, writeupID); err != nil {
	// 	zap.L().Error("Unable to edit writeup", zap.Error(err), zap.String("writeupID", writeupID))
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	decoder := json.NewDecoder(r.Body)
	var data models.Writeup
	err := decoder.Decode(&data)
	if err != nil {
		zap.L().Error("Unable to decode writeup", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeup, err := datastore.GetWriteup(ctx, writeupID)
	if err != nil {
		zap.L().Error("Unable to get writeup", zap.Error(err), zap.String("writeupID", writeupID))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data.ID = writeup.ID
	data.TS = writeup.TS
	data.VoteCount = writeup.VoteCount
	data.SubmitterID = writeup.SubmitterID

	err = datastore.UpdateWriteup(ctx, data)
	if err != nil {
		zap.L().Error("Unable to update writeup", zap.Error(err), zap.String("writeupID", writeupID))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(data)
}

func DeleteWriteupHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	writeupID := mux.Vars(r)["writeupID"]

	writeup, err := datastore.GetWriteup(ctx, writeupID)
	if err != nil {
		zap.L().Error("Unable to get writeup", zap.Error(err), zap.String("writeupID", writeupID))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = datastore.DeleteWriteup(ctx, writeupID)
	if err != nil {
		zap.L().Error("Unable to delete writeup", zap.Error(err), zap.String("writeupID", writeupID))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = datastore.DecrementWriteupCount(ctx, writeup.CTFID.Hex())
	if err != nil {
		zap.L().Error("Unable to decrement writeup count", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode("ok")
}

func GetWriteupsForCTFHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctfID := mux.Vars(r)["ctfID"]

	writeups, err := datastore.GetWriteupsForCTF(ctx, ctfID)
	if err != nil {
		zap.L().Error("Unable to get writeups", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(writeups)
}
