package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/c0nrad/ctfwriteups/datastore"
	"github.com/c0nrad/ctfwriteups/models"
)

var ADMINS = map[string]bool{"stuart@ctfwriteups.org": true}
var MODERATORS = map[string]bool{}

func GetUser(ctx context.Context, r *http.Request) (*models.User, error) {
	token, err := GetSessionToken(r)
	if err != nil {
		return nil, err
	}

	return datastore.GetUserBySession(ctx, token)
}

func GetEmail(ctx context.Context, r *http.Request) (string, error) {
	token, err := GetSessionToken(r)
	if err != nil {
		return "", err
	}

	sess, err := datastore.GetSessionByToken(ctx, token)
	if err != nil {
		return "", err
	}

	return sess.Email, nil
}

func GetSessionToken(r *http.Request) (string, error) {
	c, err := r.Cookie(SessionCookieName)
	if err != nil {
		return "", err
	}

	return c.Value, nil
}

func CanEditWriteup(ctx context.Context, r *http.Request, writeupID string) error {
	user, err := GetUser(ctx, r)
	if err != nil {
		return err
	}

	writeup, err := datastore.GetWriteup(ctx, writeupID)
	if err != nil {
		return err
	}

	if ADMINS[user.Email] {
		return nil
	}

	if user.ID.Hex() == writeup.SubmitterID.Hex() {
		return nil
	}

	return errors.New("invalid access")
}

// func HasAccessViaAPIToken(ctx context.Context, id, tokenStr, projectIDStr, role string) error {

// 	projectID, err := primitive.ObjectIDFromHex(projectIDStr)
// 	if err != nil {
// 		return err
// 	}

// 	token, err := datastore.GetAPITokenByID(ctx, id)
// 	if err != nil {
// 		return err
// 	}

// 	if !models.IsValidPassword(tokenStr, token.Token) {
// 		return errors.New("invalid token for id")
// 	}

// 	if token.ProjectID.Hex() != projectID.Hex() {
// 		return errors.New("does not have access to project")
// 	}

// 	if !models.HasRequiredRole(role, token.Role) {
// 		return errors.New("does not meet role requirements. has: " + token.Role + " requires: " + role)
// 	}

// 	return nil
// }

// func GetAPIToken(r *http.Request) (string, string, error) {
// 	tokenStr := r.Header.Get("authorization")
// 	if tokenStr == "" {
// 		return "", "", errors.New("missing header")
// 	}

// 	if !strings.HasPrefix(tokenStr, "bearer ") && !strings.HasPrefix(tokenStr, "Bearer ") {
// 		return "", "", errors.New("not a bearer token")
// 	}

// 	if len(strings.Split(tokenStr, " ")) != 2 {
// 		return "", "", errors.New("invalid bearer token format (multiple spaces)")
// 	}

// 	tokenStr = strings.Split(tokenStr, " ")[1]

// 	if len(strings.Split(tokenStr, ":")) != 2 {
// 		return "", "", errors.New("invalid bearer token format (multiple or missing : split)")
// 	}

// 	return strings.Split(tokenStr, ":")[0], strings.Split(tokenStr, ":")[1], nil

// }
