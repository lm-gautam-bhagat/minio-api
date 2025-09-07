package minioclient

import "context"

type ServiceI interface {
	GetAllBuckets(ctx context.Context) ([]string, error)
	CreateBucket(ctx context.Context, bucketName string) error
	UploadStream(ctx context.Context, stream UploadFile) (*string, error)
	UploadImageString(ctx context.Context, stream UploadFile) (*string, error)
}
