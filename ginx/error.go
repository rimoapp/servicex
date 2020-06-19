package ginx

import (
	"net/http"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/rimoapp/rimo-backend/lib/logx"
)

var closers []func()

func Init(dsn string) {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:         dsn,
		Environment: os.Getenv("APP_ENV"),
	}); err != nil {
		panic(err)
	}

	closers = append(closers, func() { sentry.Flush(30 * time.Second) })
}

func Close() {
	for _, f := range closers {
		f()
	}
}

func ErrorSetStatus(ctx *gin.Context, status int, err error) {
	logx.LogError(ctx, err.Error())

	hub := sentrygin.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub()
	}

	hub.CaptureException(err)

	ctx.AbortWithStatusJSON(
		status,
		map[string]string{"message": err.Error()},
	)
}

func Error(ctx *gin.Context, err error) {
	ErrorSetStatus(ctx, http.StatusInternalServerError, err)
}

func NoticeOnlyError(ctx *gin.Context, msg string) {
	hub := sentrygin.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub()
	}

	hub.CaptureException(errors.New(msg))
}
