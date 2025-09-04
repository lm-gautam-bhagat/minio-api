package router

import "github.com/gin-gonic/gin"

type HTTPHandler struct {
	Version    int
	Method     string
	Path       string
	Middleware []gin.HandlerFunc
	Handler    func(*SessionContext)
}

type HTTPHandlerProvider interface {
	GetHTTPHandler() []*HTTPHandler
}
