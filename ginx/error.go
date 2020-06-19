package ginx

import (
	"net/http"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/kpango/glg"
)

func ErrorSetStatus(ctx *gin.Context, status int, err error) {
	glg.Error(err)

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

func NoticeOnlyError(ctx *gin.Context, err error) {
	glg.Error(err)

	hub := sentrygin.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub()
	}

	hub.CaptureException(err)
}
