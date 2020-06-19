package servicex

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
)

func init() {
	orDie := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	loadEnv()

	err := initSentry()
	if appEnv != "development" && appEnv != "test" {
		orDie(err)
	}
}

func loadEnv() {
	if flag.Lookup("test.v") != nil {
		appEnv = Env("test")
	} else {
		var v string
		_ = setFromEnv("APP_ENV", &v)
		if v == "" {
			v = "development"
		}
		appEnv = Env(v)
	}

	env := string(appEnv)
	godotenv.Load(".env." + env + ".local")
	if "test" != env {
		godotenv.Load(".env.local")
	}
	godotenv.Load(".env." + env)
	godotenv.Load() // The Original .env
}

func initSentry() error {
	var dsn string
	err := setFromEnv("SENTRY_DSN", &dsn)
	if err != nil {
		return err
	}

	err = sentry.Init(sentry.ClientOptions{
		Dsn:         dsn,
		Environment: string(AppEnv()),
	})
	if err != nil {
		return err
	}

	closers = append(closers, func() { sentry.Flush(30 * time.Second) })

	return nil
}

var closers []func()

func Close() {
	for _, f := range closers {
		f()
	}
}

var (
	appEnv Env
)

type Env string

func AppEnv() Env { return appEnv }

type EnvVarNotSet struct{ key string }

func (e *EnvVarNotSet) Error() string { return fmt.Sprintf("%s is not set", e.key) }

type EnvVarEmpty struct{ key string }

func (e *EnvVarEmpty) Error() string { return fmt.Sprintf("%s is empty", e.key) }

func setFromEnv(key string, dest *string) error {
	v, ok := os.LookupEnv(key)
	if !ok {
		return &EnvVarNotSet{key: key}
	}

	if v == "" {
		return &EnvVarEmpty{key: key}
	}

	*dest = v

	return nil
}
