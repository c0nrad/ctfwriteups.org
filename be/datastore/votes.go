package datastore

import (
	"context"

	"github.com/c0nrad/ctfwriteups/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SaveVote(ctx context.Context, upvote models.Vote) error {
	_, err := Client.Database(DB).Collection(VoteCollection).InsertOne(ctx, upvote)
	return err
}

func GetVotesForUser(ctx context.Context, userIDStr string) ([]models.Vote, error) {
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return nil, err
	}

	var votes []models.Vote
	cursor, err := Client.Database(DB).Collection(VoteCollection).Find(ctx, bson.M{"userid": userID})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &votes)
	return votes, err
}

func GetVote(ctx context.Context, voteIDStr string) (*models.Vote, error) {
	voteID, err := primitive.ObjectIDFromHex(voteIDStr)
	if err != nil {
		return nil, err
	}

	var vote models.Vote
	err = Client.Database(DB).Collection(VoteCollection).FindOne(ctx, bson.M{"_id": voteID}).Decode(&vote)
	return &vote, err
}

func GetVotesForWriteup(ctx context.Context, writeupIDStr string) ([]models.Vote, error) {
	writeupID, err := primitive.ObjectIDFromHex(writeupIDStr)
	if err != nil {
		return nil, err
	}

	var votes []models.Vote
	cursor, err := Client.Database(DB).Collection(VoteCollection).Find(ctx, bson.M{"writeupid": writeupID})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &votes)
	return votes, err
}

func DeleteVote(ctx context.Context, voteIDStr string) error {
	voteID, err := primitive.ObjectIDFromHex(voteIDStr)
	if err != nil {
		return err
	}

	_, err = Client.Database(DB).Collection(VoteCollection).DeleteOne(ctx, bson.M{"_id": voteID})
	return err
}
