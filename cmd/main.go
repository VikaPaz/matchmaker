package main

import (
	"log"
	"os"

	"github.com/VikaPaz/matchmaker/internal/app"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	_ "github.com/VikaPaz/matchmaker/docs"
)

// @title Matchmaker API
// @description This is matchmaker server.
// @host localhost:8900
func main() {
	if err := godotenv.Overload("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	level, err := logrus.ParseLevel(os.Getenv("LOGS_LEVEL"))
	if err != nil {
		log.Fatal(err)
	}

	logger := NewLogger(level, &logrus.TextFormatter{
		FullTimestamp: true,
	})

	err = app.Run(logger)
	if err != nil {
		logger.Fatal(err)
	}
}

func NewLogger(level logrus.Level, formatter logrus.Formatter) *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(level)
	logger.SetFormatter(formatter)
	return logger
}
