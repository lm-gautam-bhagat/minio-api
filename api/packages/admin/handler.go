package mapiadmin

import (
	"net/http"

	"github.com/lm-gautam-bhagat/minio-server/api/router"
	"github.com/lm-gautam-bhagat/minio-server/log"
)

type AdminHandler struct {
	service AdminServiceI
}

func NewAdminHandler(ser AdminServiceI) *AdminHandler {
	return &AdminHandler{
		service: ser,
	}
}

func (h *AdminHandler) GetHTTPHandler() []*router.HTTPHandler {
	return []*router.HTTPHandler{
		{
			Version: 1,
			Method:  http.MethodPut,
			Path:    "buckets/:bucket",
			Handler: h.MakeBucketPublicPrivate,
		},
		{
			Version: 1,
			Method:  http.MethodPost,
			Path:    "users/create",
			Handler: h.CreateUser,
		},
		{
			Version: 1,
			Method:  http.MethodGet,
			Path:    "users",
			Handler: h.ListUsers,
		},
		{
			Version: 1,
			Method:  http.MethodDelete,
			Path:    "users/:username",
			Handler: h.DeleteUser,
		},
	}
}

func (h *AdminHandler) MakeBucketPublicPrivate(c *router.SessionContext) {
	ctx, cancel := c.GetContext()
	defer cancel()

	p := c.Query("public")
	bucket := c.Param("bucket")

	var isPublic bool
	switch p {
	case "true", "1":
		isPublic = true
	case "false", "0":
		isPublic = false
	default:
		c.RespondError(router.ErrResponseObj{
			Code:    http.StatusBadRequest,
			Message: "invalid value for 'public' query param; must be 'true', 'false', '1', or '0'",
		})
		return
	}

	err := h.service.SetBucketPolicy(ctx, bucket, isPublic)
	if err != nil {
		c.RespondError(router.ErrResponseObj{
			Code:    http.StatusInternalServerError,
			Message: "failed to change bucket policy",
		})
		return
	}

	c.Respond(http.StatusNoContent, "", "")
}

func (h *AdminHandler) CreateUser(c *router.SessionContext) {
	ctx, cancel := c.GetContext()
	defer cancel()

	var req NewUserReq

	err := c.BindJSON(&req)
	if err != nil {
		log.Error("Error while parsing request body: ", err.Error())
		c.RespondError(router.ErrResponseObj{
			Code:    http.StatusBadRequest,
			Message: "invalid request body",
		})
		return
	}

	valid := req.Policy.IsValid()
	if !valid {
		c.RespondError(router.ErrResponseObj{
			Code:    http.StatusBadRequest,
			Message: "invalid policy",
		})
		return
	}

	err = h.service.AddNewUser(ctx, req)
	if err != nil {
		c.RespondError(router.ErrResponseObj{
			Code:    http.StatusInternalServerError,
			Message: "failed to add user",
		})
		return
	}

	c.Respond(http.StatusNoContent, "", "")
}

func (h *AdminHandler) ListUsers(c *router.SessionContext) {
	ctx, cancel := c.GetContext()
	defer cancel()

	users, err := h.service.ListUsers(ctx)
	if err != nil {
		c.RespondError(router.ErrResponseObj{
			Code:    http.StatusInternalServerError,
			Message: "failed to get users",
		})
	}
	c.Respond(http.StatusAccepted, "users", users)
}

func (h *AdminHandler) DeleteUser(c *router.SessionContext) {
	ctx, cancel := c.GetContext()
	defer cancel()

	accessID := c.Param("username")
	err := h.service.DeleteUser(ctx, accessID)
	if err != nil {
		c.RespondError(router.ErrResponseObj{
			Code:    http.StatusInternalServerError,
			Message: "failed to get users",
		})
	}
	c.Respond(http.StatusAccepted, "", "")
}
