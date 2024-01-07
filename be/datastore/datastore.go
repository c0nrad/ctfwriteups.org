package datastore

import (
	"context"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.uber.org/zap"
)

var Client *mongo.Client
var DB = ""

const (
	WriteupCollection   = "writeups"
	CommentCollection   = "comments"
	VoteCollection      = "votes"
	TagCollection       = "tags"
	SeenCollection      = "seens"
	UserCollection      = "users"
	CTFCollection       = "ctfs"
	ChallengeCollection = "challenges"

	UserSessionCollection = "sessions"

	NewsletterSubscriptionCollection = "newslettersubscriptions"

	NewsletterCollection = "newsletters"
)

func InitDatabase() {

	uri := os.Getenv("MONGODB_URI")
	DB = os.Getenv("MONGODB_DB")

	if uri == "" {
		zap.L().Fatal("MONGODB_URI not set")
	}
	var err error
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	Client, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		zap.L().Fatal("Error connecting to MongoDB", zap.Error(err), zap.String("uri", uri))
	}

	if err := Client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}

	zap.L().Info("Connected to MongoDB!", zap.String("uri", strings.Split(uri, "@")[1]))

	EnsureIndexes()
}

func EnsureIndexes() {
	_, err := Client.Database(DB).Collection(WriteupCollection).Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{"ctfid", 1}, {"challengeid", 1}, {"url", 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		zap.L().Error("Error creating index", zap.Error(err))
	}

	_, err = Client.Database(DB).Collection(CTFCollection).Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{"name", 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		zap.L().Error("Error creating index", zap.Error(err))
	}

	_, err = Client.Database(DB).Collection(NewsletterSubscriptionCollection).Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{"email", 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		zap.L().Error("Error creating index", zap.Error(err))
	}

	_, err = Client.Database(DB).Collection(ChallengeCollection).Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{"name", 1}, {"category", 1}, {"ctfid", 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		zap.L().Error("Error creating index", zap.Error(err))
	}

	_, err = Client.Database(DB).Collection(VoteCollection).Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{"userid", 1}, {"writeupid", 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		zap.L().Error("Error creating index", zap.Error(err))
	}

	_, err = Client.Database(DB).Collection(SeenCollection).Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{"userid", 1}, {"writeupid", 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		zap.L().Error("Error creating index", zap.Error(err))
	}
}
