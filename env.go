package servicex

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	// EnvUndefined represents the environment is not set.
	EnvUndefined AppEnv = iota
	// EnvUnknown represents the unknown environment.
	EnvUnknown
	// EnvTest represents the test environment.
	EnvTest
	// EnvDev represents the development environment.
	EnvDev
	// EnvStg represents the staging environment.
	EnvStg
	// EnvProd represents the production environment.
	EnvProd
)

var (
	appEnv AppEnv = EnvUndefined
)

// AppEnv represents an application runtime environment.
type AppEnv int

// String returns string environment.
func (e AppEnv) String() string {
	switch e {
	case EnvUndefined:
		return "undefined"
	case EnvUnknown:
		return "unknown"
	case EnvTest:
		return "test"
	case EnvDev:
		return "development"
	case EnvStg:
		return "staging"
	case EnvProd:
		return "production"
	default:
		return "error"
	}
}

func toAppEnv(str string) AppEnv {
	str = strings.ToLower(str)
	switch str {
	case EnvTest.String():
		return EnvTest
	case EnvDev.String():
		return EnvDev
	case EnvStg.String():
		return EnvStg
	case EnvProd.String():
		return EnvProd
	default:
		return EnvUnknown
	}
}

// EnvVarNotSet representes the error that environment variable is not set.
type EnvVarNotSet struct{ key string }

func (e *EnvVarNotSet) Error() string { return fmt.Sprintf("%s is not set", e.key) }

// EnvVarEmpty representes the error that environment variable is empty.
type EnvVarEmpty struct{ key string }

func (e *EnvVarEmpty) Error() string { return fmt.Sprintf("%s is empty", e.key) }

// GetEnv returns os environment variables.
func GetEnv(key string) (string, error) {
	v, ok := os.LookupEnv(key)
	if !ok {
		return "", &EnvVarNotSet{key: key}
	}
	if v == "" {
		return "", &EnvVarEmpty{key: key}
	}
	return v, nil
}

func getAppEnv() AppEnv {
	if appEnv != EnvUndefined {
		return appEnv
	}
	if flag.Lookup("test.v") != nil {
		appEnv = EnvTest
	}
	envStr, err := GetEnv("APP_ENV")
	if err != nil {
		log.Println("You should explicitly indicate APP_ENV environment. AppEnv becomes Dev")
		appEnv = EnvDev
	} else {
		appEnv = toAppEnv(envStr)
		if appEnv == EnvUnknown {
			panic(errors.New("the unknown APP_ENV environment is set."))
		}
	}
	return appEnv
}

// IsTest returns whether app is in test environment or not.
func IsTest() bool {
	return getAppEnv() == EnvTest
}

// IsDevelopment returns whether app is in development environment or not.
func IsDevelopment() bool {
	return getAppEnv() == EnvDev
}

// IsStaging returns whether app is in staging environment or not.
func IsStaging() bool {
	return getAppEnv() == EnvStg
}

// IsProduction returns whether app is in production environment or not.
func IsProduction() bool {
	return getAppEnv() == EnvProd
}
