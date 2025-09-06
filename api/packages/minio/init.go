package minioclient

import "github.com/lm-gautam-bhagat/minio-server/storage"

func NewMinioModule(str *storage.StorageClient) *Handler {
	ser := NewMiniIOClientService(str)
	hldr := NewHandler(ser)
	return hldr
}
