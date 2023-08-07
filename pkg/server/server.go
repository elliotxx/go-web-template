package server

import (
	"fmt"
	"path/filepath"
	"time"

	recovery "github.com/akkuman/gin-logrus-recovery"
	"github.com/elliotxx/go-web-template/pkg/route"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/gorm"
)

type Config struct {
	LoggingDirectory string
	DB               *gorm.DB
}

func NewConfig() *Config {
	return &Config{}
}

type AppServer struct {
	ginEngine *gin.Engine
	route     *route.Route
}

// New creates a new AppServer instance from Config
func (c *Config) New() (*AppServer, error) {
	// Initialize the gin engine and route
	engine := NewGinEngine(c)
	router := &route.Route{
		DB: c.DB,
	}
	err := router.Register(engine)
	if err != nil {
		return nil, err
	}

	return &AppServer{
		ginEngine: engine,
		route:     router,
	}, nil
}

// PreRun is a function that will be called before the server starts to run
func (s *AppServer) PreRun() error {
	_ = logrus.WithFields(logrus.Fields{"func": "PreRun"})

	return nil
}

// Run is a function that will be called when the server starts to run
func (s *AppServer) Run(addr ...string) error {
	return s.ginEngine.Run(addr...)
}

// NewGinEngine creates a new GinEngine instance
func NewGinEngine(c *Config) *gin.Engine {
	// Create the audit writer by lumberjack
	auditLoggingFile := filepath.Join(c.LoggingDirectory, "audit.log")
	auditRotateWriter := &lumberjack.Logger{
		Filename:   auditLoggingFile,
		MaxSize:    1,
		MaxBackups: 10,
		MaxAge:     30,
	}

	// Set global config for gin
	gin.DefaultWriter = auditRotateWriter
	gin.DefaultErrorWriter = auditRotateWriter

	// Create a gin router
	r := gin.New()

	// Use some middlewares
	r.Use(requestid.New())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	// NOTE: cors.Default() allows all origins
	r.Use(cors.Default())
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: CustomLogFormatter,
		Output:    auditRotateWriter,
	}))
	r.Use(recovery.Recovery(logrus.StandardLogger()))

	return r
}

// CustomLogFormatter is a custom formatter for logging messages
func CustomLogFormatter(param gin.LogFormatterParams) string {
	// Custom format:
	// [GIN] | 200 |      14.701µs |       127.0.0.1 | GET      "/healthz"
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency -= param.Latency % time.Second
	}
	return fmt.Sprintf("[GIN] | %s | %s%3d%s | %13v | %15s | %s%-7s%s | %#v\n%s",
		param.TimeStamp.Format("2006-1-2 15:04:05.000"),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.ClientIP,
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,
	)
}
