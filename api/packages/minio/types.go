package minioclient

import (
	"io"
)

type CreateBucketReq struct {
	Name string `json:"name"`
}

type UploadImageStringReq struct {
	Base64 string `json:"base64"`
}

type UploadFile struct {
	bucket      string
	fileName    string
	file        io.Reader
	fileSize    int64
	contentType string
	Base64      string
}
