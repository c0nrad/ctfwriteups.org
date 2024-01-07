package main

import (
	"net/http"
	"os"

	"github.com/c0nrad/ctfwriteups/api"
	"github.com/c0nrad/ctfwriteups/config"
	"github.com/c0nrad/ctfwriteups/datastore"
	"go.uber.org/zap"
)

func main() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}

	config.InitLogger()
	config.InitEnv(env)
	datastore.InitDatabase()

	router := api.BuildRouter()

	mux := http.NewServeMux()
	mux.Handle("/", router)

	// if env == "dev" || config.GetConfig("DAEMON") == true {
	// go lib.RunDaemon()
	// }

	zap.L().Info("Starting server", zap.String("port", "8080"), zap.String("env", config.GetEnv()))
	err := http.ListenAndServe(":8080", mux)

	zap.L().Fatal("Error starting server", zap.Error(err))
}
