package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/c0nrad/ctfwriteups/datastore"
	"github.com/c0nrad/ctfwriteups/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

func CreateNewsletterSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	decoder := json.NewDecoder(r.Body)
	var data models.NewsletterSubscription
	err := decoder.Decode(&data)
	if err != nil {
		zap.L().Error("Unable to decode newsletter subscription", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data.ID = primitive.NewObjectID()
	data.TS = time.Now()

	err = datastore.SaveNewsletterSubscription(ctx, data)
	if err != nil {
		zap.L().Error("Unable to save newsletter subscription", zap.Error(err), zap.Any("newsletter subscription", data))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(data)
}

func RemoveNewsletterSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	decoder := json.NewDecoder(r.Body)
	var data models.NewsletterSubscription
	err := decoder.Decode(&data)
	if err != nil {
		zap.L().Error("Unable to decode newsletter subscription", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = datastore.RemoveNewsletterSubscription(ctx, data.Email)
	if err != nil {
		zap.L().Error("Unable to remove newsletter subscription", zap.Error(err), zap.Any("newsletter subscription", data))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(data)
}
