package api

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func BuildRouter() http.Handler {
	SetupCookieStore()
	SetupOAuthConfigs()

	r := mux.NewRouter()

	// Writeup
	r.HandleFunc("/api/v1/writeups", GetWriteupsHandler).Methods("GET")
	r.HandleFunc("/api/v1/writeups", CreateWriteupHandler).Methods("POST")
	r.HandleFunc("/api/v1/writeups/{writeupID}", GetWriteupHandler).Methods("GET")
	r.HandleFunc("/api/v1/writeups/{writeupID}", UpdateWriteupHandler).Methods("PUT")
	r.HandleFunc("/api/v1/writeups/{writeupID}", DeleteWriteupHandler).Methods("DELETE")
	r.HandleFunc("/api/v1/ctfs/{ctfID}/writeups", GetWriteupsForCTFHandler).Methods("GET")

	r.HandleFunc("/api/v1/ctfs", GetCTFsHandler).Methods("GET")
	r.HandleFunc("/api/v1/ctfs", CreateCTFHandler).Methods("POST")
	r.HandleFunc("/api/v1/ctfs/{ctfID}", GetCTFHandler).Methods("GET")
	r.HandleFunc("/api/v1/ctfs/{ctfID}", UpdateCTFHandler).Methods("PUT")
	r.HandleFunc("/api/v1/ctfs/{ctfID}", DeleteCTFHandler).Methods("DELETE")

	r.HandleFunc("/api/v1/ctfs/{ctfID}/challenges", GetChallengesHandler).Methods("GET")
	r.HandleFunc("/api/v1/challenges", CreateChallengeHandler).Methods("POST")
	r.HandleFunc("/api/v1/challenges/{challengeID}", GetChallengeHandler).Methods("GET")
	r.HandleFunc("/api/v1/challenges/{challengeID}", UpdateChallengeHandler).Methods("PUT")
	r.HandleFunc("/api/v1/challenges/{challengeID}", DeleteChallengeHandler).Methods("DELETE")

	// votes
	r.HandleFunc("/api/v1/users/me/votes", GetVotesForUserHandler).Methods("GET")
	r.HandleFunc("/api/v1/votes", CreateVoteForWriteupHandler).Methods("POST")
	r.HandleFunc("/api/v1/votes/{voteID}", DeleteVoteForWriteupHandler).Methods("DELETE")

	// seens
	r.HandleFunc("/api/v1/users/me/seens", GetSeensForUserHandler).Methods("GET")
	r.HandleFunc("/api/v1/seens", IncrementSeenHandler).Methods("POST")
	r.HandleFunc("/api/v1/seens/{seenID}", DeleteSeenHandler).Methods("DELETE")

	// comments
	r.HandleFunc("/api/v1/comments", CreateCommentHandler).Methods("POST")
	r.HandleFunc("/api/v1/comments/{commentID}", DeleteCommentHandler).Methods("DELETE")
	r.HandleFunc("/api/v1/writeups/{writeupID}/comments", GetCommentsForWriteupHandler).Methods("GET")

	// newsletter
	r.HandleFunc("/api/v1/newsletter", CreateNewsletterSubscriptionHandler).Methods("POST")
	r.HandleFunc("/api/v1/newsletter/unsubscribe", RemoveNewsletterSubscriptionHandler).Methods("POST")

	// user
	r.HandleFunc("/api/v1/users/me", GetMeHandler).Methods("GET")
	r.HandleFunc("/login/google", LoginGoogleHandler)
	r.HandleFunc("/login/google/callback", GoogleOAuthCallbackHandler)
	r.HandleFunc("/login/github", LoginGithubHandler)
	r.HandleFunc("/login/github/callback", GithubOAuthCallbackHandler)
	r.HandleFunc("/logout", LogoutHandler)

	r.HandleFunc("/health", HealthEndpoint).Methods("GET")

	r.PathPrefix("/").HandlerFunc(AngularFileHandler)
	r.Use(LoggingMiddleware)

	return r
}

func HealthEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "alive")
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().UTC()

		recorder := NewResponseRecorder(w)
		next.ServeHTTP(recorder, r)
		// body := recorder.Body

		end := time.Now().UTC()
		latency := end.Sub(start)

		fields := []zapcore.Field{
			zap.String("status", fmt.Sprintf("%d", recorder.Status)),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Duration("duration", latency),
			zap.Int("size", recorder.Size),
		}

		message := fmt.Sprintf("%-6s %s", r.Method, r.URL.String())

		if recorder.Status < 400 {
			zap.L().Info(message, fields...)
		} else if recorder.Status < 500 {
			zap.L().Warn(message, fields...)
		} else {
			zap.L().Error(message, fields...)
		}
	})
}

func NewResponseRecorder(w http.ResponseWriter) *ResponseRecorder {
	return &ResponseRecorder{
		ResponseWriter: w,
		Status:         http.StatusOK,
	}
}

type ResponseRecorder struct {
	http.ResponseWriter
	Status int
	Size   int
	Body   []byte
}

func (r *ResponseRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *ResponseRecorder) Write(buf []byte) (int, error) {
	n, err := r.ResponseWriter.Write(buf)
	if err == nil {
		r.Size += n
		r.Body = append(r.Body, buf...)
	}
	return n, err
}

func (r *ResponseRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := r.ResponseWriter.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, errors.New("no hijack on original responsewriter")
}

func ToObjectID(str string) primitive.ObjectID {
	oid, err := primitive.ObjectIDFromHex(str)
	if err != nil {
		return primitive.NilObjectID
	}
	return oid
}

func AngularFileHandler(w http.ResponseWriter, r *http.Request) {
	dir, err := os.Getwd()
	baseDist := dir + "/dist"
	requestedFilePath := filepath.Clean(filepath.Join(baseDist, r.URL.Path))

	if !strings.HasPrefix(requestedFilePath, baseDist) {
		zap.L().Warn("someone is performing directory traversal?", zap.String("url", r.URL.Path), zap.String("requestedFilePath", requestedFilePath), zap.String("clean", filepath.Clean(filepath.Join(baseDist, r.URL.Path))), zap.String("join", filepath.Join(baseDist, r.URL.Path)))
		http.Error(w, "invalid filepath", 400)
		return
	}

	_, err = os.Stat(requestedFilePath)
	if err == nil {
		if !strings.Contains(requestedFilePath, "index.html") && requestedFilePath != baseDist {
			cacheSince := time.Now().Format(http.TimeFormat)
			cacheUntil := time.Now().AddDate(60, 0, 0).Format(http.TimeFormat)

			w.Header().Set("Cache-Control", "max-age:290304000, public")
			w.Header().Set("Last-Modified", cacheSince)
			w.Header().Set("Expires", cacheUntil)
		}

		http.ServeFile(w, r, requestedFilePath)
		return
	}

	// default to index.html
	http.ServeFile(w, r, baseDist+"/index.html")
}
