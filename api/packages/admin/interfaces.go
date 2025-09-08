package mapiadmin

import "context"

type AdminServiceI interface {
	SetBucketPolicy(ctx context.Context, bucket string, makePublic bool) error
	AddNewUser(ctx context.Context, usr NewUserReq) error
}
