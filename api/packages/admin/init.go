package mapiadmin

import "github.com/lm-gautam-bhagat/minio-server/storage"

func NewMinioAdminModule(strA *storage.StorageAdmin, strC *storage.StorageClient) *AdminHandler {
	ser := NewMiniIOAdminService(strA, strC)
	hldr := NewAdminHandler(ser)
	return hldr
}
