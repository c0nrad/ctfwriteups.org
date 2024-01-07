package datastore

import (
	"context"

	"github.com/c0nrad/ctfwriteups/models"
	"go.mongodb.org/mongo-driver/bson"
)

func SaveNewsletterSubscription(ctx context.Context, newsletterSubscription models.NewsletterSubscription) error {
	_, err := Client.Database(DB).Collection(NewsletterSubscriptionCollection).InsertOne(ctx, newsletterSubscription)
	return err
}

func RemoveNewsletterSubscription(ctx context.Context, email string) error {
	_, err := Client.Database(DB).Collection(NewsletterSubscriptionCollection).DeleteOne(ctx, bson.M{"email": email})
	return err
}

func GetNewsletterSubscriptions(ctx context.Context) ([]models.NewsletterSubscription, error) {
	var newsletterSubscriptions []models.NewsletterSubscription
	cursor, err := Client.Database(DB).Collection(NewsletterSubscriptionCollection).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &newsletterSubscriptions)
	return newsletterSubscriptions, err
}
