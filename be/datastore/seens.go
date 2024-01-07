package datastore

import (
	"context"
	"time"

	"github.com/c0nrad/ctfwriteups/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SaveSeen(ctx context.Context, upseen models.Seen) error {
	_, err := Client.Database(DB).Collection(SeenCollection).InsertOne(ctx, upseen)
	return err
}

func GetSeensForUser(ctx context.Context, userIDStr string) ([]models.Seen, error) {
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return nil, err
	}

	var seens []models.Seen
	cursor, err := Client.Database(DB).Collection(SeenCollection).Find(ctx, bson.M{"userid": userID})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &seens)
	return seens, err
}

func GetSeen(ctx context.Context, seenIDStr string) (*models.Seen, error) {
	seenID, err := primitive.ObjectIDFromHex(seenIDStr)
	if err != nil {
		return nil, err
	}

	var seen models.Seen
	err = Client.Database(DB).Collection(SeenCollection).FindOne(ctx, bson.M{"_id": seenID}).Decode(&seen)
	return &seen, err
}

func DeleteSeen(ctx context.Context, seenIDStr string) error {
	seenID, err := primitive.ObjectIDFromHex(seenIDStr)
	if err != nil {
		return err
	}

	_, err = Client.Database(DB).Collection(SeenCollection).DeleteOne(ctx, bson.M{"_id": seenID})
	return err
}

func IncrementSeen(ctx context.Context, seen models.Seen) error {
	// upsert
	_, err := Client.Database(DB).Collection(SeenCollection).UpdateOne(
		ctx,
		bson.M{"userid": seen.UserID, "writeupid": seen.WriteupID},
		bson.M{"$inc": bson.M{"count": 1}, "$setOnInsert": bson.M{"ts": time.Now()}},
		options.Update().SetUpsert(true),
	)
	return err
}
