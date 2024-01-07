package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var ENV = "test"

func GetEnv() string {
	return ENV
}

func GetConfig(v string) string {
	return os.Getenv(v)
}

func InitLogger() {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	logger, err := config.Build()

	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)

}

func InitEnv(env string) {
	ENV = env

	dir, _ := os.Getwd()

	config_path := dir + "/../../config/.env." + strings.ToLower(env)
	config_path = filepath.Clean(config_path)

	if err := godotenv.Load(config_path); err != nil {

		config_path := dir + "/../config/.env." + strings.ToLower(env)
		config_path = filepath.Clean(config_path)
		if err := godotenv.Load(config_path); err != nil {
			zap.L().Fatal("Error loading .env file", zap.Error(err), zap.String("path", config_path))
			log.Println("No .env file found")
		}
	}
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		zap.L().Fatal("MONGODB_URI not set")
	}
}
