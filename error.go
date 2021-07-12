package servicex

import (
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

func getSentryHub(c *gin.Context) *sentry.Hub {
	hub := sentrygin.GetHubFromContext(c)
	if hub == nil {
		hub = sentry.CurrentHub()
	}
	return hub
}

// RenderError reports error to sentry and renders error json.
func RenderError(c *gin.Context, status int, err error) {
	LogError(c, err.Error())

	hub := getSentryHub(c)
	hub.CaptureException(err)

	c.AbortWithStatusJSON(
		status,
		map[string]string{"error": err.Error(), "message": err.Error()},
	)
}

// NoticeOnlyError sends error to sentry but do nothing else.
func NoticeOnlyError(c *gin.Context, err error) {
	// cf. https://docs.sentry.io/platforms/go/enriching-events/scopes/#local-scopes
	hub := getSentryHub(c)
	hub.WithScope(func(scope *sentry.Scope) {
		scope.SetTag("notice-only", "true")
		scope.SetLevel(sentry.LevelWarning)
		hub.CaptureException(err)
	})
}

// SetUserContextForErrorReporting sets user infomation for sentry.
// You should call this in authentication function.
// cf. https://docs.sentry.io/platforms/go/enriching-events/identify-user/
func SetUserContextForErrorReporting(c *gin.Context, userID, email, name string) {
	hub := getSentryHub(c)
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetUser(sentry.User{
			ID:       userID,
			Email:    email,
			Username: name,
		})
	})
}
