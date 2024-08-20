package models

import "errors"

var (
	ErrLoadEnvFailed = errors.New("failed to load environment")
	ErrServerFailed  = errors.New("failed to connect to server")
)
