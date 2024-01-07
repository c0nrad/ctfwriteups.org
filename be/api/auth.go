package api

import (
	"context"
	"encoding/json"
	"io"

	"net/http"
	"strings"
	"time"

	"github.com/c0nrad/ctfwriteups/config"
	"github.com/c0nrad/ctfwriteups/datastore"
	"github.com/c0nrad/ctfwriteups/models"
	"github.com/google/go-github/github"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var CookieStore *sessions.CookieStore
var OAuthCookieName = "oauth_session"
var SessionCookieName = "session"

func GetUUID() string {
	return uuid.New().String()
}

func SetupCookieStore() {
	secret := config.GetConfig("COOKIE_STORE_SECRET")

	if secret == "" || len(secret) != 32 {
		panic("invalid cookie secret")
	}

	zap.L().Info("cookieStore setup")

	CookieStore = sessions.NewCookieStore([]byte(secret))
}

var GoogleOAuthConf = &oauth2.Config{}
var GithubOAuthConf = &oauth2.Config{}

func SetupOAuthConfigs() {

	configItems := []string{
		"OAUTH_GOOGLE_ID",
		"OAUTH_GOOGLE_SECRET",
		"OAUTH_GITHUB_ID",
		"OAUTH_GITHUB_SECRET",
	}

	for _, i := range configItems {
		if config.GetConfig(i) == "" {
			panic(i + "is not defined")
		}
	}

	GoogleOAuthConf = &oauth2.Config{
		ClientID:     config.GetConfig("OAUTH_GOOGLE_ID"),
		ClientSecret: config.GetConfig("OAUTH_GOOGLE_SECRET"),
		RedirectURL:  config.GetConfig("API_ORIGIN") + "/login/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

	GithubOAuthConf = &oauth2.Config{
		ClientID:     config.GetConfig("OAUTH_GITHUB_ID"),
		ClientSecret: config.GetConfig("OAUTH_GITHUB_SECRET"),
		RedirectURL:  config.GetConfig("API_ORIGIN") + "/login/github/callback",
		Scopes: []string{
			"read:user",
			"user:email",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}
}

type GoogleProfileStruct struct {
	ID    string
	Email string
	Name  string
}

func LoginGoogleHandler(w http.ResponseWriter, r *http.Request) {

	session, err := CookieStore.Get(r, OAuthCookieName)
	if err != nil {
		zap.L().Error("error getting cookie store session", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	state := GetUUID()
	session.Values["state"] = state
	session.Save(r, w)

	url := GoogleOAuthConf.AuthCodeURL(state)

	http.Redirect(w, r, url, http.StatusSeeOther)
}

func LoginGithubHandler(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()

	session, err := CookieStore.Get(r, OAuthCookieName)
	if err != nil {
		zap.L().Error("error getting cookie store session", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	state := GetUUID()
	session.Values["state"] = state
	session.Save(r, w)

	url := GithubOAuthConf.AuthCodeURL(state)

	http.Redirect(w, r, url, http.StatusSeeOther)
}

func GithubOAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	oauthSession, err := CookieStore.Get(r, OAuthCookieName)
	if err != nil {

		zap.L().Error("error getting cookie store session", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	state := r.FormValue("state")
	authcode := r.FormValue("code")

	if len(state) == 0 || len(authcode) == 0 {
		zap.L().Error("invalid state or authcode", zap.String("state", state), zap.String("authcode", authcode), zap.String("provider", "github"))
		http.Error(w, "invalid state or authcode", http.StatusBadRequest)
		return
	}

	sessionState, _ := oauthSession.Values["state"].(string)
	if sessionState != state {
		zap.L().Error("invalid state parameter", zap.String("sessionState", sessionState), zap.String("state", state), zap.String("provider", "github"))
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}

	oauthSession.Values["state"] = ""
	err = oauthSession.Save(r, w)
	if err != nil {
		zap.L().Error("error saving oauth state", zap.Error(err), zap.String("provider", "github"))
		http.Error(w, "error saving oauth state: "+err.Error(), http.StatusBadRequest)
		return
	}

	token, err := GithubOAuthConf.Exchange(context.Background(), authcode)
	if err != nil {
		zap.L().Error("error exchanging code", zap.Error(err), zap.String("provider", "github"), zap.Int("authCode.length", len(authcode)))
		http.Error(w, "error exchanging code: "+err.Error(), http.StatusBadRequest)
		return
	}

	oauthClient := GithubOAuthConf.Client(context.Background(), token)
	client := github.NewClient(oauthClient)
	profile, _, err := client.Users.Get(ctx, "")
	if err != nil {
		zap.L().Error("error getting github user", zap.Error(err), zap.String("provider", "github"))
		http.Error(w, "error getting github user: "+err.Error(), http.StatusBadRequest)
		return
	}

	emails, _, err := client.Users.ListEmails(ctx, nil)
	if err != nil {
		zap.L().Error("error getting github emails", zap.Error(err), zap.String("provider", "github"))
		http.Error(w, "unable to get email: "+err.Error(), http.StatusBadRequest)
		return
	}
	for _, email := range emails {
		if email.GetPrimary() {
			profile.Email = email.Email
		}
	}

	// maybe there is no primary?
	if profile.Email == nil {
		for _, email := range emails {
			if email.GetVerified() {
				profile.Email = email.Email
			}
		}
	}

	user, isNew, err := datastore.GetOrCreateUser(ctx, *profile.Email, "github")
	if err != nil {
		zap.L().Error("error creating user", zap.Error(err), zap.String("provider", "github"), zap.String("email", *profile.Email))
		http.Error(w, "error creating user", http.StatusBadRequest)
		return
	}

	firstName := ""
	lastName := ""
	if profile.Name != nil {
		name := *profile.Name
		if len(strings.Split(name, " ")) >= 2 {
			firstName = strings.Split(name, " ")[0]
			lastName = strings.Split(name, " ")[1]
		} else if len(strings.Split(name, " ")) >= 1 {
			firstName = name
		}
	}

	user, err = datastore.UpdateUser(ctx, user.ID, firstName, lastName)
	if err != nil {
		zap.L().Error("unable to update user", zap.Error(err), zap.String("provider", "github"), zap.String("firstName", firstName), zap.String("lastName", lastName))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	session, err := datastore.CreateSessionForUser(ctx, user.Email, "github")
	if err != nil {
		zap.L().Error("unable to create session", zap.Error(err), zap.String("provider", "github"), zap.String("email", user.Email))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	zap.L().Info("issuing session for user", zap.String("userId", user.ID.Hex()), zap.String("firstName", firstName), zap.String("lastName", lastName), zap.String("email", user.Email))

	StartSession(w, session)

	if isNew {
		http.Redirect(w, r, "/", 307)
	} else {
		http.Redirect(w, r, "/", 307)
	}
}

func GoogleOAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	oauthSession, err := CookieStore.Get(r, OAuthCookieName)
	if err != nil {
		zap.L().Error("error getting cookie store session", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	state := r.FormValue("state")
	authcode := r.FormValue("code")

	sessionState, _ := oauthSession.Values["state"].(string)
	if sessionState != state {
		zap.L().Error("invalid state parameter", zap.String("sessionState", sessionState), zap.String("state", state), zap.String("provider", "google"))
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}

	oauthSession.Values["state"] = ""
	err = oauthSession.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tok, err := GoogleOAuthConf.Exchange(oauth2.NoContext, authcode)
	if err != nil {
		zap.L().Error("error exchanging code", zap.Error(err), zap.String("provider", "google"), zap.Int("authCode.length", len(authcode)))
		http.Error(w, "error exchanging code", http.StatusBadRequest)
		return
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + tok.AccessToken)
	if err != nil {
		zap.L().Error("error grabbing profile", zap.Error(err), zap.Int("authCode.length", len(authcode)), zap.Int("response.statusCode", response.StatusCode))
		http.Error(w, "error grabbing profile", http.StatusBadRequest)
		return
	}

	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		zap.L().Error("error reading google profile", zap.Error(err))
		http.Error(w, "error reading google profile", http.StatusBadRequest)
		return
	}

	var profile GoogleProfileStruct
	err = json.Unmarshal(contents, &profile)
	if err != nil {
		zap.L().Error("error unmarshling profile", zap.Error(err))
		http.Error(w, "error unmarshling profile", http.StatusBadRequest)
		return
	}

	user, isNew, err := datastore.GetOrCreateUser(ctx, profile.Email, "google")
	if err != nil {
		zap.L().Error("error creating user", zap.Error(err), zap.String("provider", "google"), zap.String("email", profile.Email))
		http.Error(w, "error creating user", http.StatusBadRequest)
		return
	}

	firstName := ""
	lastName := ""
	if len(strings.Split(profile.Name, " ")) >= 2 {
		firstName = strings.Split(profile.Name, " ")[0]
		lastName = strings.Split(profile.Name, " ")[1]
	} else if len(strings.Split(profile.Name, " ")) >= 1 {
		firstName = profile.Name
	}

	user, err = datastore.UpdateUser(ctx, user.ID, firstName, lastName)
	if err != nil {
		zap.L().Error("unable to update user", zap.Error(err), zap.String("provider", "google"), zap.String("userId", user.ID.Hex()), zap.String("firstName", firstName), zap.String("lastName", lastName))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	session, err := datastore.CreateSessionForUser(ctx, user.Email, "google")
	if err != nil {
		zap.L().Error("unable to create session", zap.Error(err), zap.String("provider", "google"), zap.String("email", user.Email))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	StartSession(w, session)
	if isNew {
		http.Redirect(w, r, "/", 307)
	} else {
		http.Redirect(w, r, "/", 307)
	}
}

func GetRemoteAddress(r *http.Request) string {
	return r.RemoteAddr
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token, err := GetSessionToken(r)
	if err != nil {
		zap.L().Error("error getting session token", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = datastore.LogoutUser(ctx, token)
	if err != nil {
		zap.L().Error("error logging out user", zap.Error(err), zap.String("token", token))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	EndSession(w)
	http.Redirect(w, r, "/", 307)
}

func StartSession(w http.ResponseWriter, session *models.UserSession) {

	if config.ENV == "dev" || config.ENV == "test" {
		http.SetCookie(w, &http.Cookie{
			Name:    SessionCookieName,
			Value:   session.Token,
			Expires: time.Now().Add(time.Hour * 24 * 7),
			Path:    "/",
		})
	} else {
		http.SetCookie(w, &http.Cookie{
			Name:     SessionCookieName,
			Value:    session.Token,
			HttpOnly: true,
			Secure:   true,
			Expires:  time.Now().Add(time.Hour * 24 * 7),
			Path:     "/",
			SameSite: http.SameSiteStrictMode,
		})
	}
}

func EndSession(w http.ResponseWriter) {

	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		HttpOnly: true,
		// SameSite: http.SameSiteStrictMode,
		// Secure:  true,
		Expires: time.Now().Add(-time.Second * 60 * 60 * 24),
		Path:    "/",
	})
}
