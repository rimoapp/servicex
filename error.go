package servicex

import (
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

// RenderError reports error to sentry and renders error json
func RenderError(c *gin.Context, status int, err error) {
	LogError(c, err.Error())

	hub := sentrygin.GetHubFromContext(c)
	if hub == nil {
		hub = sentry.CurrentHub()
	}

	hub.CaptureException(err)

	c.AbortWithStatusJSON(
		status,
		map[string]string{"error": err.Error(), "message": err.Error()},
	)
}

// NoticeOnlyError sends error to sentry but do nothing else.
func NoticeOnlyError(c *gin.Context, err error) {
	hub := sentrygin.GetHubFromContext(c)
	if hub == nil {
		hub = sentry.CurrentHub()
	}
	hub.WithScope(func(scope *sentry.Scope) {
		scope.SetTag("notice-only", "true")
		scope.SetLevel(sentry.LevelWarning)
		hub.CaptureException(err)
	})
}
