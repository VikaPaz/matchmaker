package app

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/VikaPaz/matchmaker/internal/models"
	"github.com/VikaPaz/matchmaker/internal/repository"
	"github.com/VikaPaz/matchmaker/internal/server"
	"github.com/VikaPaz/matchmaker/internal/service"
	"github.com/sirupsen/logrus"
)

type RedisConfig struct {
	DB repository.Config
}

func Run(logger *logrus.Logger) error {
	group_size, err := strconv.Atoi(os.Getenv("GROUP_SIZE"))
	if err != nil {
		logger.Errorf("Error loading .env file: %v", err)
		return models.ErrLoadEnvFailed
	}

	group_wait, err := strconv.Atoi(os.Getenv("GROUP_WAIT"))
	if err != nil {
		logger.Errorf("Error loading .env file: %v", err)
		return models.ErrLoadEnvFailed
	}

	maxSkill, err := strconv.ParseFloat(os.Getenv("MAX_SKILL"), 64)
	if err != nil {
		logger.Errorf("Error loading .env file: %v", err)
		return models.ErrLoadEnvFailed
	}

	maxLatency, err := strconv.ParseFloat(os.Getenv("MAX_LATENCY"), 64)
	if err != nil {
		logger.Errorf("Error loading .env file: %v", err)
		return models.ErrLoadEnvFailed
	}

	confRedis := RedisConfig{
		DB: repository.Config{
			Host:     os.Getenv("HOST"),
			Port:     os.Getenv("REDIS_PORT"),
			Password: os.Getenv("PASSWORD"),
		},
	}

	redisConn, err := repository.Connection(confRedis.DB)
	if err != nil {
		return err
	}
	repo := repository.NewRepo(redisConn, logger)

	matcher := service.NewService(repo, maxSkill, maxLatency, logger)

	go func() {
		for {
			err := matcher.Matching(group_size)
			if err != nil {
				logger.Fatal("Error Matching")
			}
			time.Sleep(time.Duration(group_wait * int(time.Millisecond)))
		}
	}()

	srv := server.NewServer(matcher, logger)

	logger.Infof("Running server on port %s", os.Getenv("PORT"))
	err = http.ListenAndServe(":"+os.Getenv("PORT"), srv.Handlers())
	if err != nil {
		logger.Errorf("Error starting server")
		return models.ErrServerFailed
	}

	return err

}
