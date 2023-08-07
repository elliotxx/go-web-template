package healthz

import (
	"github.com/elliotxx/healthcheck"
	"github.com/elliotxx/healthcheck/checks"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Register registers the livez and readyz handlers to the specified
// router.
func Register(r *gin.RouterGroup, db *gorm.DB) {
	r.GET("/livez", NewLivezHandler())
	r.GET("/readyz", NewReadyzHandler(db))
}

// NewLivezHandler creates a new liveness check handler that can be
// used to check if the application is running.
func NewLivezHandler() gin.HandlerFunc {
	conf := healthcheck.HandlerConfig{
		Verbose: false,
		// checkList is a list of healthcheck to run.
		Checks: []checks.Check{
			checks.NewPingCheck(),
		},
		FailureNotification: healthcheck.FailureNotification{Threshold: 1},
	}

	return healthcheck.NewHandler(conf)
}

// NewReadyzHandler creates a new readiness check handler that can be
// used to check if the application is ready to serve traffic.
func NewReadyzHandler(db *gorm.DB) gin.HandlerFunc {
	conf := healthcheck.HandlerConfig{
		Verbose: true,
		// checkList is a list of healthcheck to run.
		Checks: []checks.Check{
			checks.NewPingCheck(),
			NewGormDBCheck(db),
		},
		FailureNotification: healthcheck.FailureNotification{Threshold: 1},
	}

	return healthcheck.NewHandler(conf)
}
