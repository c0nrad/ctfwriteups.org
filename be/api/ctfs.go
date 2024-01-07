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

func GetCTFsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctfs, err := datastore.GetCTFs(ctx)
	if err != nil {
		zap.L().Error("Unable to get ctfs", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(ctfs)
}

func CreateCTFHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	decoder := json.NewDecoder(r.Body)
	var data models.CTF
	err := decoder.Decode(&data)
	if err != nil {
		zap.L().Error("Unable to decode ctf", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data.ID = primitive.NewObjectID()
	data.TS = time.Now()

	err = datastore.SaveCTF(ctx, data)
	if err != nil {
		zap.L().Error("Unable to save ctf", zap.Error(err), zap.Any("ctf", data))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(data)
}

func GetCTFHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctfID := mux.Vars(r)["ctfID"]

	ctf, err := datastore.GetCTF(ctx, ctfID)
	if err != nil {
		zap.L().Error("Unable to get ctf", zap.Error(err), zap.String("ctfID", ctfID))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(ctf)
}

func UpdateCTFHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctfID := mux.Vars(r)["ctfID"]

	// if err := CanEditCTF(ctx, r, ctfID); err != nil {
	// 	zap.L().Error("Unable to edit ctf", zap.Error(err), zap.String("ctfID", ctfID))
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	decoder := json.NewDecoder(r.Body)
	var data models.CTF
	err := decoder.Decode(&data)
	if err != nil {
		zap.L().Error("Unable to decode ctf", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctf, err := datastore.GetCTF(ctx, ctfID)
	if err != nil {
		zap.L().Error("Unable to get ctf", zap.Error(err), zap.String("ctfID", ctfID))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data.ID = ctf.ID
	data.TS = ctf.TS
	data.ChallengeCount = ctf.ChallengeCount
	data.WriteupCount = ctf.WriteupCount
	data.SubmitterID = ctf.SubmitterID

	err = datastore.UpdateCTF(ctx, data)
	if err != nil {
		zap.L().Error("Unable to update ctf", zap.Error(err), zap.String("ctfID", ctfID))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(data)
}

func DeleteCTFHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctfID := mux.Vars(r)["ctfID"]

	err := datastore.DeleteCTF(ctx, ctfID)
	if err != nil {
		zap.L().Error("Unable to delete ctf", zap.Error(err), zap.String("ctfID", ctfID))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode("ok")
}
