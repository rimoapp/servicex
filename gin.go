package servicex

import (
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

// DefaultGinEngine returns gin engine with needed middleware.
func DefaultGinEngine() *gin.Engine {
	router := gin.Default()

	// Send request infomation to sentry when panic.
	router.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))

	return router
}
