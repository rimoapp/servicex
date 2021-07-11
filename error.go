package servicex

import (
	"net/http"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// RenderError reports error to sentry and renders error json
func RenderError(ctx *gin.Context, status int, err error) {
	LogError(ctx, err.Error())

	hub := sentrygin.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub()
	}

	hub.CaptureException(err)

	ctx.AbortWithStatusJSON(
		status,
		map[string]string{"error": err.Error(), "message": err.Error()},
	)
}

// NoticeOnlyError sends error to sentry but do nothing else.
func NoticeOnlyError(c *gin.Context, err error) {
	hub := sentrygin.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub()
	}

	hub.CaptureException(err))
}
