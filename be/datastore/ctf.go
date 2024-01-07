package datastore

import (
	"context"

	"github.com/c0nrad/ctfwriteups/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetCTFByName(ctx context.Context, name string) (*models.CTF, error) {
	var ctf models.CTF
	err := Client.Database(DB).Collection(CTFCollection).FindOne(ctx, bson.M{"name": name}).Decode(&ctf)
	return &ctf, err
}

func SaveCTF(ctx context.Context, ctf models.CTF) error {
	_, err := Client.Database(DB).Collection(CTFCollection).InsertOne(ctx, ctf)
	return err
}

func GetCTFs(ctx context.Context) ([]models.CTF, error) {
	var ctfs []models.CTF
	cursor, err := Client.Database(DB).Collection(CTFCollection).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &ctfs)
	return ctfs, err
}

func GetCTF(ctx context.Context, ctfIDStr string) (*models.CTF, error) {
	ctfID, err := primitive.ObjectIDFromHex(ctfIDStr)
	if err != nil {
		return nil, err
	}

	var ctf models.CTF
	err = Client.Database(DB).Collection(CTFCollection).FindOne(ctx, bson.M{"_id": ctfID}).Decode(&ctf)
	return &ctf, err
}

func UpdateCTF(ctx context.Context, ctf models.CTF) error {
	_, err := Client.Database(DB).Collection(CTFCollection).UpdateOne(ctx, bson.M{"_id": ctf.ID}, bson.M{"$set": ctf})
	return err
}

func DeleteCTF(ctx context.Context, ctfIDStr string) error {
	ctfID, err := primitive.ObjectIDFromHex(ctfIDStr)
	if err != nil {
		return err
	}

	_, err = Client.Database(DB).Collection(CTFCollection).DeleteOne(ctx, bson.M{"_id": ctfID})
	return err
}

func IncrementChallengeCount(ctx context.Context, ctfIDStr string) error {
	ctfID, err := primitive.ObjectIDFromHex(ctfIDStr)
	if err != nil {
		return err
	}

	_, err = Client.Database(DB).Collection(CTFCollection).UpdateOne(ctx, bson.M{"_id": ctfID}, bson.M{"$inc": bson.M{"challengecount": 1}})
	return err
}

func DecrementChallengeCount(ctx context.Context, ctfIDStr string) error {
	ctfID, err := primitive.ObjectIDFromHex(ctfIDStr)
	if err != nil {
		return err
	}

	_, err = Client.Database(DB).Collection(CTFCollection).UpdateOne(ctx, bson.M{"_id": ctfID}, bson.M{"$inc": bson.M{"challengecount": -1}})
	return err
}

func IncrementWriteupCount(ctx context.Context, ctfIDStr string) error {
	ctfID, err := primitive.ObjectIDFromHex(ctfIDStr)
	if err != nil {
		return err
	}

	_, err = Client.Database(DB).Collection(CTFCollection).UpdateOne(ctx, bson.M{"_id": ctfID}, bson.M{"$inc": bson.M{"writeupcount": 1}})
	return err
}

func DecrementWriteupCount(ctx context.Context, ctfIDStr string) error {
	ctfID, err := primitive.ObjectIDFromHex(ctfIDStr)
	if err != nil {
		return err
	}

	_, err = Client.Database(DB).Collection(CTFCollection).UpdateOne(ctx, bson.M{"_id": ctfID}, bson.M{"$inc": bson.M{"writeupcount": -1}})
	return err
}
