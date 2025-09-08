package mapiadmin

import (
	"context"

	"github.com/lm-gautam-bhagat/minio-server/log"
	"github.com/lm-gautam-bhagat/minio-server/storage"
)

type MinIOAdminService struct {
	strAdmin  *storage.StorageAdmin
	strClient *storage.StorageClient
}

func NewMiniIOAdminService(strA *storage.StorageAdmin, strC *storage.StorageClient) AdminServiceI {
	return &MinIOAdminService{
		strAdmin:  strA,
		strClient: strC,
	}
}

func (s *MinIOAdminService) SetBucketPolicy(ctx context.Context, bucket string, makePublic bool) error {
	err := s.strClient.SetBucketPublicAccess(ctx, bucket, makePublic)
	if err != nil {
		log.Error("Error while changing bucket policy: ", err.Error())
		return err
	}
	return nil
}

func (s *MinIOAdminService) AddNewUser(ctx context.Context, usr NewUserReq) error {
	err := s.strAdmin.CreateUser(ctx, usr.AccessID, usr.AccessKey)
	if err != nil {
		log.Error("Error while creating a new user: ", err.Error())
		return err
	}
	return nil
}
