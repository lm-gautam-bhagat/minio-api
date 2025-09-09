package mapiadmin

import (
	"context"

	"github.com/lm-gautam-bhagat/minio-server/constants"
	"github.com/lm-gautam-bhagat/minio-server/log"
	"github.com/lm-gautam-bhagat/minio-server/storage"
	"github.com/minio/madmin-go/v4"
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
	err = s.SetUserPolicy(ctx, usr.AccessID, usr.Policy)
	if err != nil {
		log.Error("Error while assigning policy to user: ", err.Error())
		return err
	}
	return nil
}

func (s *MinIOAdminService) DeleteUser(ctx context.Context, accessID string) error {
	err := s.strAdmin.DeleteUser(ctx, accessID)
	if err != nil {
		log.Error("Error while creating a new user: ", err.Error())
		return err
	}

	return nil
}

func (s *MinIOAdminService) SetUserPolicy(ctx context.Context, username string, policy constants.PolicyType) error {
	_, err := s.strAdmin.AttachPolicy(ctx, username, []string{string(policy)})
	if err != nil {
		return err
	}
	return nil
}

func (s *MinIOAdminService) ListUsers(ctx context.Context) (map[string]madmin.UserInfo, error) {
	users, err := s.strAdmin.ListUsers(ctx)
	if err != nil {
		log.Error("Error while getting users: ", err.Error())
		return nil, err
	}
	return users, nil
}
