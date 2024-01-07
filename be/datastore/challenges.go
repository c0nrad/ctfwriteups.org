package datastore

import (
	"context"

	"github.com/c0nrad/ctfwriteups/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SaveChallenge(ctx context.Context, challenge models.Challenge) error {
	_, err := Client.Database(DB).Collection(ChallengeCollection).InsertOne(ctx, challenge)
	return err
}

func GetChallenges(ctx context.Context) ([]models.Challenge, error) {
	var challenges []models.Challenge
	cursor, err := Client.Database(DB).Collection(ChallengeCollection).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &challenges)
	return challenges, err
}

func GetChallengesForCTF(ctx context.Context, ctfIDStr string) ([]models.Challenge, error) {
	ctfID, err := primitive.ObjectIDFromHex(ctfIDStr)
	if err != nil {
		return nil, err
	}

	var challenges []models.Challenge
	cursor, err := Client.Database(DB).Collection(ChallengeCollection).Find(ctx, bson.M{"ctfid": ctfID})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &challenges)
	return challenges, err
}

func GetChallenge(ctx context.Context, challengeIDStr string) (*models.Challenge, error) {
	challengeID, err := primitive.ObjectIDFromHex(challengeIDStr)
	if err != nil {
		return nil, err
	}

	var challenge models.Challenge
	err = Client.Database(DB).Collection(ChallengeCollection).FindOne(ctx, bson.M{"_id": challengeID}).Decode(&challenge)
	return &challenge, err
}

func UpdateChallenge(ctx context.Context, challenge models.Challenge) error {
	_, err := Client.Database(DB).Collection(ChallengeCollection).UpdateOne(ctx, bson.M{"_id": challenge.ID}, bson.M{"$set": challenge})
	return err
}

func DeleteChallenge(ctx context.Context, challengeIDStr string) error {
	challengeID, err := primitive.ObjectIDFromHex(challengeIDStr)
	if err != nil {
		return err
	}

	_, err = Client.Database(DB).Collection(ChallengeCollection).DeleteOne(ctx, bson.M{"_id": challengeID})
	return err
}

func GetChallengeByName(ctx context.Context, name string, category string, ctfIDStr string) (*models.Challenge, error) {
	ctfID, err := primitive.ObjectIDFromHex(ctfIDStr)
	if err != nil {
		return nil, err
	}

	var challenge models.Challenge
	err = Client.Database(DB).Collection(ChallengeCollection).FindOne(ctx, bson.M{"name": name, "category": category, "ctfid": ctfID}).Decode(&challenge)
	return &challenge, err
}
