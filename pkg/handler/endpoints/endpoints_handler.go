package endpoints

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/elliotxx/go-web-template/pkg/util/misc"

	"github.com/gin-gonic/gin"
)

// NewEndpointsGETHandler returns a handler that lists all API
// endpoints for the GET method.
func NewEndpointsGETHandler(routes gin.RoutesInfo) gin.HandlerFunc {
	return func(c *gin.Context) {
		// endpointMethodPathPattern is the format string used to
		// output each endpoint.
		endpointMethodPathPattern := "%s\t%s"

		// Iterate over all the routes and add each endpoint to the
		// endpoints slice.
		endpoints := []string{}
		for _, route := range routes {
			endpoints = append(
				endpoints,
				fmt.Sprintf(endpointMethodPathPattern, route.Method, route.Path))
		}

		sort.Strings(endpoints)
		c.String(http.StatusOK, strings.Join(endpoints, "\n"))
	}
}

// NewEndpointsOPTIONSHandler returns a handler that lists all API
// endpoints for the OPTIONS method.
func NewEndpointsOPTIONSHandler(routes gin.RoutesInfo) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Iterate over all the routes and add each endpoint to the
		// endpoints slice and add the HTTP method to the allow sets
		// map.
		allow := misc.NewSet()
		endpoints := []string{}
		for _, route := range routes {
			allow.Set(route.Method)
			endpoints = append(endpoints, route.Path)
		}

		// Sort the endpoints in alphabetical orders.
		sort.Strings(endpoints)

		// Set the response headers and status code.
		c.Header("Allow", allow.String())
		c.Header("API-Endpoints", strings.Join(endpoints, ","))
		c.Header("Content-Length", "0")
		c.Status(http.StatusOK)
	}
}
