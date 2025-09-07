package minioclient

import (
	"net/http"

	"github.com/lm-gautam-bhagat/minio-server/api/packages/router"
	"github.com/lm-gautam-bhagat/minio-server/log"
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
		{
			Version: 1,
			Method:  http.MethodPost,
			Path:    "buckets/:bucket/upload",
			Handler: h.UploadFile,
		},
		{
			Version: 1,
			Method:  http.MethodPost,
			Path:    "buckets/:bucket/upload/image/string",
			Handler: h.UploadImageBase64,
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

func (h *Handler) UploadFile(c *router.SessionContext) {
	ctx, cancel := c.GetContext()
	defer cancel()

	bucket := c.Param("bucket")

	fileHeader, err := c.FormFile("file")
	if err != nil {
		log.Error("Error while getting file: ", err.Error())
		c.RespondError(router.ErrResponseObj{Code: http.StatusBadRequest, Message: "file is required"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		log.Error("Error while opening file: ", err.Error())
		c.RespondError(router.ErrResponseObj{Code: http.StatusInternalServerError, Message: "cannot open uploaded file"})
		return

	}
	defer file.Close()

	stream := UploadFile{
		bucket:      bucket,
		fileName:    fileHeader.Filename,
		file:        file,
		fileSize:    fileHeader.Size,
		contentType: fileHeader.Header.Get("Content-Type"),
	}

	url, err := h.service.UploadStream(
		ctx,
		stream,
	)
	if err != nil {
		c.RespondError(router.ErrResponseObj{Code: http.StatusInternalServerError, Message: "failed to upload file"})
		return
	}

	c.Respond(http.StatusCreated, "public_id", *url)
}

func (h *Handler) UploadImageBase64(c *router.SessionContext) {
	ctx, cancel := c.GetContext()
	defer cancel()

	bucket := c.Param("bucket")

	var req UploadImageStringReq

	err := c.BindJSON(&req)
	if err != nil {
		c.RespondError(router.ErrResponseObj{
			Code:    http.StatusBadRequest,
			Message: "invalid body request",
		})
	}

	stream := UploadFile{
		bucket: bucket,
		Base64: req.Base64,
	}

	url, err := h.service.UploadImageString(
		ctx,
		stream,
	)
	if err != nil {
		c.RespondError(router.ErrResponseObj{Code: http.StatusInternalServerError, Message: "failed to upload file"})
		return
	}

	c.Respond(http.StatusCreated, "public_id", *url)
}
