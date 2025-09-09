package mapiadmin

import (
	"context"

	"github.com/lm-gautam-bhagat/minio-server/constants"
	"github.com/minio/madmin-go/v4"
)

type AdminServiceI interface {
	SetBucketPolicy(ctx context.Context, bucket string, makePublic bool) error
	AddNewUser(ctx context.Context, usr NewUserReq) error
	ListUsers(ctx context.Context) (map[string]madmin.UserInfo, error)
	SetUserPolicy(ctx context.Context, accessID string, policy constants.PolicyType) error
	DeleteUser(ctx context.Context, accessID string) error
}
