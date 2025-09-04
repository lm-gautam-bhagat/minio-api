package router

import "github.com/gin-gonic/gin"

type SessionContext struct {
	GinC *gin.Context
}

func NewSessionContext(c *gin.Context) *SessionContext {
	return &SessionContext{GinC: c}
}
