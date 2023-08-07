package route

import (
	"github.com/elliotxx/expvar"
	docs "github.com/elliotxx/go-web-template/api/openapispec"
	"github.com/elliotxx/go-web-template/pkg/handler"
	"github.com/elliotxx/go-web-template/pkg/handler/api/v1/systemconfig"
	"github.com/elliotxx/go-web-template/pkg/handler/debug/statsviz"
	"github.com/elliotxx/go-web-template/pkg/handler/endpoints"
	"github.com/elliotxx/go-web-template/pkg/handler/healthz"
	"github.com/elliotxx/go-web-template/pkg/infrastructure/persistence"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type Route struct {
	DB *gorm.DB
}

// Register registers some api to the route
func (r *Route) Register(engine *gin.Engine) error {
	// Create the workspace domain service
	systemConfigHandler := systemconfig.NewHandler(persistence.NewSystemConfigRepository(r.DB))

	// Registers some api to the route
	docs.SwaggerInfo.BasePath = "/"
	root := engine.Group("/")
	{
		root.GET("/livez", healthz.NewLivezHandler())
		root.GET("/readyz", healthz.NewReadyzHandler(r.DB))
		root.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
	debug := engine.Group("/debug")
	{
		// Add a expvar handler for gin framework, expvar provides
		// a standardized interface to public variables.
		// You can visit http://localhost/debug/vars to view all
		// public variables.
		debug.GET("/vars", expvar.Handler(expvar.WithFilters("memstats")))
		// Register statsviz handler
		debug.GET("/statsviz/*filepath", statsviz.StatsvizHandler)
		// The default pprof router is /debug/pprof
		pprof.RouteRegister(debug, "/pprof")
	}
	apiv1 := engine.Group("/api/v1")
	{
		// Register system config handler
		apiv1.POST("/systemconfig", handler.WrapFD(systemConfigHandler.CreateSystemConfig))
		apiv1.DELETE("/systemconfig/:id", handler.WrapFD(systemConfigHandler.DeleteSystemConfig))
		apiv1.PUT("/systemconfig", handler.WrapFD(systemConfigHandler.UpdateSystemConfig))
		apiv1.GET("/systemconfig/:id", handler.WrapFD(systemConfigHandler.GetSystemConfig))
		apiv1.GET("/systemconfigs", handler.WrapFD(systemConfigHandler.FindSystemConfigs))
		apiv1.GET("/systemconfig/count", handler.WrapFD(systemConfigHandler.CountSystemConfigs))
	}

	engine.GET("/endpoints", endpoints.NewEndpointsGETHandler(engine.Routes()))
	engine.OPTIONS("/endpoints", endpoints.NewEndpointsOPTIONSHandler(engine.Routes()))

	return nil
}
