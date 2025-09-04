package minio

import (
	"net/http"

	"github.com/lm-gautam-bhagat/minio-server/api/packages/router"
)

type Handler struct{}

func (h *Handler) GetHTTPHandler() []*router.HTTPHandler {
	return []*router.HTTPHandler{
		{
			Version: 1,
			Method:  http.MethodGet,
			Path:    "api/new",
			Handler: h.Home,
		},
	}
}

func (*Handler) Home(c *router.SessionContext) {
	c.GinC.JSON(http.StatusAccepted, "Hi")
}
