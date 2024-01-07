package datastore

import (
	"context"

	"github.com/c0nrad/ctfwriteups/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

func NewObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func SaveWriteup(ctx context.Context, writeup models.Writeup) error {
	_, err := Client.Database(DB).Collection(WriteupCollection).InsertOne(ctx, writeup)
	return err
}

func GetWriteups(ctx context.Context) ([]models.Writeup, error) {
	var writeups []models.Writeup
	cursor, err := Client.Database(DB).Collection(WriteupCollection).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &writeups)
	return writeups, err
}

func GetWriteup(ctx context.Context, writeupIDStr string) (*models.Writeup, error) {
	writeupID, err := primitive.ObjectIDFromHex(writeupIDStr)
	if err != nil {
		return nil, err
	}

	var writeup models.Writeup
	err = Client.Database(DB).Collection(WriteupCollection).FindOne(ctx, bson.M{"_id": writeupID}).Decode(&writeup)
	return &writeup, err
}

func UpdateWriteup(ctx context.Context, writeup models.Writeup) error {
	_, err := Client.Database(DB).Collection(WriteupCollection).UpdateOne(ctx, bson.M{"_id": writeup.ID}, bson.M{"$set": writeup})
	return err
}

func DeleteWriteup(ctx context.Context, writeupIDStr string) error {
	writeupID, err := primitive.ObjectIDFromHex(writeupIDStr)
	if err != nil {
		return err
	}

	_, err = Client.Database(DB).Collection(WriteupCollection).DeleteOne(ctx, bson.M{"_id": writeupID})
	return err
}

func GetWriteupsForCTF(ctx context.Context, ctfIDStr string) ([]models.Writeup, error) {
	ctfID, err := primitive.ObjectIDFromHex(ctfIDStr)
	if err != nil {
		return nil, err
	}

	var writeups []models.Writeup
	cursor, err := Client.Database(DB).Collection(WriteupCollection).Find(ctx, bson.M{"ctfid": ctfID})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &writeups)
	return writeups, err
}

func IncrementVoteCount(ctx context.Context, writeupIDStr string) error {
	writeupID, err := primitive.ObjectIDFromHex(writeupIDStr)
	if err != nil {
		return err
	}

	_, err = Client.Database(DB).Collection(WriteupCollection).UpdateOne(ctx, bson.M{"_id": writeupID}, bson.M{"$inc": bson.M{"votecount": 1}})
	return err
}

func DecrementVoteCount(ctx context.Context, writeupIDStr string) error {
	writeupID, err := primitive.ObjectIDFromHex(writeupIDStr)
	if err != nil {
		return err
	}

	_, err = Client.Database(DB).Collection(WriteupCollection).UpdateOne(ctx, bson.M{"_id": writeupID}, bson.M{"$inc": bson.M{"votecount": -1}})
	return err
}

func IncrementCommentCount(ctx context.Context, writeupIDStr string) error {
	writeupID, err := primitive.ObjectIDFromHex(writeupIDStr)
	if err != nil {
		return err
	}

	_, err = Client.Database(DB).Collection(WriteupCollection).UpdateOne(ctx, bson.M{"_id": writeupID}, bson.M{"$inc": bson.M{"commentcount": 1}})
	return err
}

func DecrementCommentCount(ctx context.Context, writeupIDStr string) error {
	writeupID, err := primitive.ObjectIDFromHex(writeupIDStr)
	if err != nil {
		return err
	}

	_, err = Client.Database(DB).Collection(WriteupCollection).UpdateOne(ctx, bson.M{"_id": writeupID}, bson.M{"$inc": bson.M{"commentcount": -1}})
	return err
}

func UpdateChallengeTags(ctx context.Context, challengeIDStr string, tags []string) error {
	challengeID, err := primitive.ObjectIDFromHex(challengeIDStr)
	if err != nil {
		return err
	}

	_, err = Client.Database(DB).Collection(WriteupCollection).UpdateMany(ctx, bson.M{"challengeid": challengeID}, bson.M{"$addToSet": bson.M{"tags": bson.M{"$each": tags}}})
	if err != nil {
		zap.L().Error("Unable to update challenge tags", zap.Error(err), zap.String("challengeID", challengeIDStr))
		_, err = Client.Database(DB).Collection(WriteupCollection).UpdateMany(ctx, bson.M{"challengeid": challengeID}, bson.M{"$set": bson.M{"tags": tags}})
	}

	return err
}
