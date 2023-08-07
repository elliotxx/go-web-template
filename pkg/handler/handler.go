package handler

import (
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/elliotxx/go-web-template/pkg/errcode"
	"github.com/elliotxx/go-web-template/pkg/util/ctxutil"
	"github.com/elliotxx/errors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// The Handler is an abstract of handle request
type Handler interface {
	Validate(c *gin.Context) error
	Handle(c *gin.Context, log logrus.FieldLogger) (any, error)
}

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as HTTP handlers. If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler that calls f.
type HandlerFunc func(c *gin.Context, log logrus.FieldLogger) error

// The HandlerDataFunc type is an adapter to allow the use of
// ordinary functions as HTTP handlers. If f is a function
// with the appropriate signature, HandlerDataFunc(f) is a
// Handler that calls f.
type HandlerDataFunc func(c *gin.Context, log logrus.FieldLogger) (any, error)

// WrapF is a helper function for wrapping handler.HandlerFunc and returns a Gin middleware.
func WrapF(f HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a logger for current request
		log := getRequestLogger(c, f)
		// Inject the logger to context
		c.Request = c.Request.WithContext(ctxutil.WithLogger(log))

		// Handle and calculate the cost time for request
		handleRequest(c, log, f)
	}
}

// WrapFD is a helper function for wrapping handler.HandlerDataFunc and returns a Gin middleware.
func WrapFD(f HandlerDataFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a logger for current request
		log := getRequestDataLogger(c, f)
		// Inject the logger to context
		c.Request = c.Request.WithContext(ctxutil.WithLogger(log))

		// Handle and calculate the cost time for request
		handleDataRequest(c, log, f)
	}
}

// WrapH is a helper function for wrapping handler.Handler and returns a Gin middleware.
func WrapH(h Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a logger for current request
		log := getRequestDataLogger(c, h.Handle)
		// Inject the logger to context
		c.Request = c.Request.WithContext(ctxutil.WithLogger(log))

		// Logging the request start message
		loggingStartMsg(log)

		// Calculate the total cost time for request
		requestStartTime := time.Now()
		defer func() {
			requestCostTime := time.Since(requestStartTime)
			log.Debugf("Request total took [%v]", requestCostTime)
		}()

		// Validate and calculate the cost time for request
		if ok := validateRequest(c, log, h); ok {
			// Handle and calculate the cost time for request
			handleDataRequest(c, log, h.Handle)
		}
	}
}

// Create a logger for current request
func getRequestDataLogger(c *gin.Context, f HandlerDataFunc) logrus.FieldLogger {
	// Get handler name
	fullName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	fullNames := strings.Split(fullName, "/")
	var name string
	if len(fullNames) > 0 {
		name = fullNames[len(fullNames)-1]
	} else {
		name = fullName
	}

	// Create a logger with fields
	return logrus.WithFields(
		logrus.Fields{
			"traceID": requestid.Get(c),
			"handler": name,
		},
	)
}

// Create a logger for current request
func getRequestLogger(c *gin.Context, f HandlerFunc) logrus.FieldLogger {
	// Get handler name
	fullName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	fullNames := strings.Split(fullName, "/")
	var name string
	if len(fullNames) > 0 {
		name = fullNames[len(fullNames)-1]
	} else {
		name = fullName
	}

	// Create a logger with fields
	return logrus.WithFields(
		logrus.Fields{
			"traceID": requestid.Get(c),
			"handler": name,
		},
	)
}

func loggingStartMsg(log logrus.FieldLogger) {
	log.Debug("Start processing request ...")
}

// Validate and calculate the cost time for current request
func validateRequest(c *gin.Context, log logrus.FieldLogger, h Handler) bool {
	log.Debug("Start validating request ...")
	validStartTime := time.Now()
	err := h.Validate(c)
	validCostTime := time.Since(validStartTime)
	log.Debugf("Validate request took [%v]", validCostTime)

	if err != nil {
		log.Errorf("Failed to validate request: %v", err)
		response := Response{
			Success: false,
			TraceID: requestid.Get(c),
		}
		switch e := err.(type) {
		case errors.DetailError:
			response.Code = e.GetCode()
			response.Message = errors.Wrap(e.GetCause(), e.GetMsg()).Error()
			c.AbortWithStatusJSON(errcode.StatusCode(e), response)
		case errors.ErrorCode:
			response.Code = e.GetCode()
			response.Message = e.GetMsg()
			c.AbortWithStatusJSON(errcode.StatusCode(e), response)
		default:
			response.Code = errcode.InvalidParams.GetCode()
			response.Message = e.Error()
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
		}
		return false
	}

	return true
}

// Handle and calculate the cost time for current request
func handleDataRequest(c *gin.Context, log logrus.FieldLogger, f HandlerDataFunc) {
	log.Debug("Start handling request ...")
	startTime := time.Now()
	data, err := f(c, log)
	endTime := time.Now()
	costTime := endTime.Sub(startTime)
	log.Debugf("Handle request took [%v]", costTime)

	if c.IsAborted() {
		return
	}

	response := Response{
		TraceID:   requestid.Get(c),
		StartTime: startTime,
		EndTime:   endTime,
		CostTime:  Duration(costTime),
	}

	if err == nil {
		response.Success = true
		response.Message = http.StatusText(http.StatusOK)
		response.Code = errcode.Success.GetCode()
		response.Data = data
		c.AbortWithStatusJSON(http.StatusOK, response)
	} else {
		log.Errorf("Failed to handle request: %+v", err)

		response.Success = false
		switch e := err.(type) {
		case errors.DetailError:
			response.Code = e.GetCode()
			response.Message = errors.Wrap(e.GetCause(), e.GetMsg()).Error()
			c.AbortWithStatusJSON(errcode.StatusCode(e), response)
		case errors.ErrorCode:
			response.Code = e.GetCode()
			response.Message = e.GetMsg()
			c.AbortWithStatusJSON(errcode.StatusCode(e), response)
		default:
			response.Code = errcode.InternalError.GetCode()
			response.Message = e.Error()
			c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		}
	}
}

// Handle and calculate the cost time for current request
func handleRequest(c *gin.Context, log logrus.FieldLogger, f HandlerFunc) {
	log.Debug("Start handling request ...")
	handleStartTime := time.Now()
	err := f(c, log)
	handleCostTime := time.Since(handleStartTime)
	log.Debugf("Handle request took [%v]", handleCostTime)

	if err != nil {
		log.Errorf("Failed to handle request: %+v", err)
		switch e := err.(type) {
		case errors.DetailError:
			c.JSON(errcode.StatusCode(e), gin.H{"code": e.GetCode(), "msg": e.GetMsg(), "cause": e.GetCause().Error()})
		case errors.ErrorCode:
			c.JSON(errcode.StatusCode(e), gin.H{"code": e.GetCode(), "msg": e.GetMsg()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"code": errcode.InternalError.GetCode(), "cause": e.Error()})
		}
		return
	}

	c.AbortWithStatus(http.StatusOK)
}
