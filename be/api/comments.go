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

// r.HandleFunc("/api/v1/comments", CreateCommentHandler).Methods("POST")
// r.HandleFunc("/api/v1/comments/{commentID}", DeleteCommentHandler).Methods("DELETE")
// r.HandleFunc("/api/v1/writeups/{writeupID}/comments", GetCommentsForWriteupHandler).Methods("GET")

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := GetUser(ctx, r)
	if err != nil {
		http.Error(w, "not logged in", http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var data models.Comment
	err = decoder.Decode(&data)
	if err != nil {
		zap.L().Error("Unable to decode comment", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data.ID = primitive.NewObjectID()
	data.UserID = user.ID
	data.TS = time.Now()
	data.Username = user.Username

	err = datastore.SaveComment(ctx, data)
	if err != nil {
		zap.L().Error("Unable to save comment", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = datastore.IncrementCommentCount(ctx, data.WriteupID.Hex())
	if err != nil {
		zap.L().Error("Unable to increment comment count", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(data)
}

func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := GetUser(ctx, r)
	if err != nil {
		http.Error(w, "not logged in", http.StatusUnauthorized)
		return
	}

	commentID := mux.Vars(r)["commentID"]

	comment, err := datastore.GetComment(ctx, commentID)
	if err != nil {
		zap.L().Error("Unable to get comment", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if comment.UserID.Hex() != user.ID.Hex() {
		zap.L().Error("User does not own comment", zap.Error(err))
		http.Error(w, "User does not own comment", http.StatusBadRequest)
		return
	}

	err = datastore.DeleteComment(ctx, commentID)
	if err != nil {
		zap.L().Error("Unable to delete comment", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = datastore.DecrementCommentCount(ctx, comment.WriteupID.Hex())
	if err != nil {
		zap.L().Error("Unable to decrement comment count", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode("ok")
}

func GetCommentsForWriteupHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	writeupID := mux.Vars(r)["writeupID"]

	comments, err := datastore.GetCommentsForWriteup(ctx, writeupID)
	if err != nil {
		zap.L().Error("Unable to get comments", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(comments)
}
