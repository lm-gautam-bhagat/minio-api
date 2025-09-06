package minioclient

import (
	"net/http"

	"github.com/lm-gautam-bhagat/minio-server/api/packages/router"
)

type Handler struct {
	service ServiceI
}

func NewHandler(ser ServiceI) *Handler {
	return &Handler{
		service: ser,
	}
}

func (h *Handler) GetHTTPHandler() []*router.HTTPHandler {
	return []*router.HTTPHandler{
		{
			Version: 1,
			Method:  http.MethodGet,
			Path:    "home",
			Handler: h.Home,
		},
		{
			Version: 1,
			Method:  http.MethodGet,
			Path:    "buckets",
			Handler: h.GetBuckets,
		},
		{
			Version: 1,
			Method:  http.MethodPost,
			Path:    "buckets/new",
			Handler: h.CreateBucket,
		},
	}
}

func (*Handler) Home(c *router.SessionContext) {
	c.Respond(http.StatusAccepted, "msg", "Hi")
}

func (h *Handler) GetBuckets(c *router.SessionContext) {
	ctx, cancel := c.GetContext()
	defer cancel()

	buckets, err := h.service.GetAllBuckets(ctx)
	if err != nil {
		c.RespondError(
			router.ErrResponseObj{
				Code:    http.StatusInternalServerError,
				Message: "failed to get buckets",
			})
		return
	}
	c.Respond(http.StatusAccepted, "buckets", buckets)
}

func (h *Handler) CreateBucket(c *router.SessionContext) {
	ctx, cancel := c.GetContext()
	defer cancel()

	var req CreateBucketReq
	err := c.BindJSON(&req)
	if err != nil {
		c.RespondError(
			router.ErrResponseObj{
				Code:    http.StatusBadRequest,
				Message: "failed to bind request",
			})
		return
	}

	err = h.service.CreateBucket(ctx, req.Name)
	if err != nil {
		c.RespondError(
			router.ErrResponseObj{
				Code:    http.StatusInternalServerError,
				Message: "failed to create bucket",
			})
		return
	}
	c.Respond(http.StatusCreated, "", "")
}
