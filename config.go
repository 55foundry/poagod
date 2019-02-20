package main

import (
	"os"

	"github.com/Sirupsen/logrus"
)

const (
	Name = "poagod"
	Env  = "env"

	Development = "development"
	Production  = "production"

	Environment = "ENVIRONMENT"
)

func GetEnv(key, fallback string) string {
	value := os.Getenv(key)

	if len(value) == 0 {
		return fallback
	}

	return value
}

func LoadLogger() *logrus.Logger {
	log := logrus.New()
	env := GetEnv(Environment, Development)

	if env == Production {
		log.Formatter = &logrus.JSONFormatter{}
	} else {
		log.Formatter = &logrus.TextFormatter{}
	}

	log.Out = os.Stdout

	log.SetLevel(logrus.InfoLevel)
	log.WithField(Env, env)

	return log
}
