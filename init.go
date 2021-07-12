package servicex

import (
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
)

var closers []func()

// Init initializes servicex
func Init() {
	orDie := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	loadEnv()

	err := initSentry()
	if !IsDevelopment() && !IsTest() {
		orDie(err)
	}
}

// Close cleans up servicex
func Close() {
	for _, f := range closers {
		f()
	}
}

func loadEnv() {
	env := getAppEnv().String()
	godotenv.Load(".env." + env + ".local")
	if env != "test" {
		godotenv.Load(".env.local")
	}
	godotenv.Load(".env." + env)
	godotenv.Load() // The Original .env
}

func initSentry() error {
	dsn, err := GetEnv("SENTRY_DSN")
	if err != nil {
		return err
	}

	err = sentry.Init(sentry.ClientOptions{
		Dsn:         dsn,
		Environment: getAppEnv().String(),
		Release:     os.Getenv("GIT_COMMIT_SHA"), // optional
	})
	if err != nil {
		return err
	}

	closers = append(closers, func() { sentry.Flush(30 * time.Second) })

	return nil
}
