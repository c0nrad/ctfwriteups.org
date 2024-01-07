package datastore

import (
	"context"
	"time"

	"github.com/c0nrad/ctfwriteups/models"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var BannedUsers = []string{}

func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {

	var out models.User
	err := Client.Database(DB).Collection(UserCollection).FindOne(ctx, bson.M{"email": email}).Decode(&out)

	if err != nil {
		return nil, err
	}

	return &out, err
}

func GetUserByID(ctx context.Context, userID primitive.ObjectID) (*models.User, error) {
	var out models.User
	err := Client.Database(DB).Collection(UserCollection).FindOne(ctx, bson.M{"_id": userID}).Decode(&out)
	if err != nil {
		return nil, err
	}

	return &out, err
}

func GetOrCreateUser(ctx context.Context, email, source string) (*models.User, bool, error) {
	isNew := false

	user, err := GetUserByEmail(ctx, email)
	if err != nil {
		Client.Database(DB).Collection(UserCollection).InsertOne(ctx, models.User{
			TS:    time.Now(),
			ID:    primitive.NewObjectID(),
			Email: email,
		})
		isNew = true

		user, err = GetUserByEmail(ctx, email)
	}

	return user, isNew, err
}

func LogoutUser(ctx context.Context, token string) error {
	_, err := Client.Database(DB).Collection(UserSessionCollection).DeleteOne(ctx, bson.M{"token": token})

	return err
}

func GetUserBySession(ctx context.Context, token string) (*models.User, error) {

	var session models.UserSession
	err := Client.Database(DB).Collection(UserSessionCollection).FindOne(ctx, bson.M{"token": token}).Decode(&session)
	if err != nil {
		return nil, err
	}

	return GetUserByEmail(ctx, session.Email)
}

func GetSessionByToken(ctx context.Context, token string) (*models.UserSession, error) {

	var session models.UserSession
	err := Client.Database(DB).Collection(UserSessionCollection).FindOne(ctx, bson.M{"token": token}).Decode(&session)
	return &session, err
}

func GetSessionsForUser(ctx context.Context, email string) ([]models.UserSession, error) {

	cursor, err := Client.Database(DB).Collection(UserSessionCollection).Find(ctx, bson.M{"email": email}, options.Find().SetSort(bson.M{"ts": -1}))
	if err != nil {
		return nil, err
	}

	sessions := []models.UserSession{}
	err = cursor.All(context.TODO(), &sessions)
	return sessions, err
}

func GetUsers(ctx context.Context) ([]models.User, error) {

	cursor, err := Client.Database(DB).Collection(UserCollection).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	users := []models.User{}
	err = cursor.All(context.TODO(), &users)
	return users, err
}

func GetUsersSearch(ctx context.Context, searchStr string, page int) ([]models.User, error) {

	pageSize := int64(5)

	findOptions := options.Find()
	findOptions.SetLimit(pageSize)
	findOptions.SetSkip(pageSize * int64(page))

	query := bson.M{}
	if searchStr != "" {
		query = bson.M{"$text": bson.M{"$search": searchStr}}
	}

	cursor, err := Client.Database(DB).Collection(UserCollection).Find(ctx, query, findOptions)
	if err != nil {
		return nil, err
	}

	users := []models.User{}
	err = cursor.All(context.TODO(), &users)
	return users, err

}

func CreateSessionForUser(ctx context.Context, email string, authMechanism string) (*models.UserSession, error) {

	session := models.UserSession{
		TS: time.Now(),
		ID: primitive.NewObjectID(),

		Email: email,
		Token: uuid.New().String(),
	}

	_, err := Client.Database(DB).Collection(UserSessionCollection).InsertOne(ctx, session)

	return &session, err
}

func UpdateUser(ctx context.Context, userID primitive.ObjectID, firstName, lastName string) (*models.User, error) {

	_, err := Client.Database(DB).Collection(UserCollection).UpdateOne(ctx, bson.M{"_id": userID}, bson.M{
		"$set": bson.M{
			"firstname": firstName,
			"lastname":  lastName,
		}})

	if err != nil {
		return nil, err
	}

	return GetUserByID(ctx, userID)
}

func PurgeUser(ctx context.Context, userID primitive.ObjectID) error {

	zap.L().Warn("PurgeUser", zap.String("userID", userID.Hex()))

	_, err := Client.Database(DB).Collection(UserCollection).DeleteOne(ctx, bson.M{"_id": userID})
	return err
}
