package mapiadmin

import (
	"context"

	"github.com/minio/madmin-go/v4"
)

type AdminServiceI interface {
	SetBucketPolicy(ctx context.Context, bucket string, makePublic bool) error
	AddNewUser(ctx context.Context, usr NewUserReq) error
	ListUsers(ctx context.Context) (map[string]madmin.UserInfo, error)
	SetUserPolicy(ctx context.Context, username, policy string) error
}
