package statsviz

import (
	"github.com/arl/statsviz"
	"github.com/gin-gonic/gin"
)

func StatsvizHandler(context *gin.Context) {
	if context.Param("filepath") == "/ws" {
		statsviz.Ws(context.Writer, context.Request)
		return
	}
	statsviz.IndexAtRoot("/debug/statsviz").ServeHTTP(context.Writer, context.Request)
}
