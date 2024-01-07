package datastore

import (
	"context"

	"github.com/c0nrad/ctfwriteups/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SaveComment(ctx context.Context, comment models.Comment) error {
	_, err := Client.Database(DB).Collection(CommentCollection).InsertOne(ctx, comment)
	return err
}

func GetComment(ctx context.Context, commentIDStr string) (*models.Comment, error) {
	commentID, err := primitive.ObjectIDFromHex(commentIDStr)
	if err != nil {
		return nil, err
	}

	var comment models.Comment
	err = Client.Database(DB).Collection(CommentCollection).FindOne(ctx, bson.M{"_id": commentID}).Decode(&comment)
	return &comment, err
}

func GetCommentsForUser(ctx context.Context, userIDStr string) ([]models.Comment, error) {
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return nil, err
	}

	var comments []models.Comment
	cursor, err := Client.Database(DB).Collection(CommentCollection).Find(ctx, bson.M{"userid": userID})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &comments)
	return comments, err
}

func GetCommentsForWriteup(ctx context.Context, writeupIDStr string) ([]models.Comment, error) {
	writeupID, err := primitive.ObjectIDFromHex(writeupIDStr)
	if err != nil {
		return nil, err
	}

	var comments []models.Comment
	cursor, err := Client.Database(DB).Collection(CommentCollection).Find(ctx, bson.M{"writeupid": writeupID})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &comments)
	return comments, err
}

func DeleteComment(ctx context.Context, commentIDStr string) error {
	commentID, err := primitive.ObjectIDFromHex(commentIDStr)
	if err != nil {
		return err
	}

	_, err = Client.Database(DB).Collection(CommentCollection).DeleteOne(ctx, bson.M{"_id": commentID})
	return err
}
