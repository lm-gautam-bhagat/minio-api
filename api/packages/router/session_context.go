package router

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/gin-gonic/gin"
)

type SessionContext struct {
	ginCtx *gin.Context
}

type ErrResponseObj struct {
	Code    int
	Message string
}

func NewSessionContext(c *gin.Context) *SessionContext {
	return &SessionContext{ginCtx: c}
}

func (sc *SessionContext) Respond(code int, key string, data any) {
	sc.ginCtx.JSON(code, gin.H{
		"data": gin.H{
			key: data,
		},
	})
}

func (sc *SessionContext) RespondError(obj ErrResponseObj) {
	sc.ginCtx.JSON(obj.Code, gin.H{
		"message": obj.Message,
	})
}

// String response
func (sc *SessionContext) String(code int, msg string) {
	sc.ginCtx.String(code, msg)
}

// Get path param
func (sc *SessionContext) Param(key string) string {
	return sc.ginCtx.Param(key)
}

// Get query param
func (sc *SessionContext) Query(key string) string {
	return sc.ginCtx.Query(key)
}

// Get form value (POST form-data/x-www-form-urlencoded)
func (sc *SessionContext) PostForm(key string) string {
	return sc.ginCtx.PostForm(key)
}

// Bind JSON body to struct
func (sc *SessionContext) BindJSON(obj any) error {
	return sc.ginCtx.ShouldBindJSON(obj)
}

// Set value in context
func (sc *SessionContext) Set(key string, val any) {
	sc.ginCtx.Set(key, val)
}

// Get value from context
func (sc *SessionContext) Get(key string) (any, bool) {
	return sc.ginCtx.Get(key)
}

func (sc *SessionContext) GetContextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(sc.ginCtx.Request.Context(), timeout)
}

func (sc *SessionContext) GetContext() (context.Context, context.CancelFunc) {
	return context.WithCancel(sc.ginCtx.Request.Context())
}

func (sc *SessionContext) BindForm(obj any) error {
	return sc.ginCtx.ShouldBind(obj)
}

func (sc *SessionContext) FormFile(name string) (*multipart.FileHeader, error) {
	return sc.ginCtx.FormFile(name)
}
