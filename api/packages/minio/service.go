package minioclient

import (
	"context"
	"fmt"

	"github.com/lm-gautam-bhagat/minio-server/log"
	"github.com/lm-gautam-bhagat/minio-server/storage"
)

type MinioClientService struct {
	str *storage.StorageClient
}

func NewMiniIOClientService(str *storage.StorageClient) ServiceI {
	return &MinioClientService{
		str: str,
	}
}

func (s *MinioClientService) GetAllBuckets(ctx context.Context) ([]string, error) {
	buckets, err := s.str.ListBuckets(ctx)
	if err != nil {
		log.Error(fmt.Sprintf("Error while getting buckets: %s", err.Error()))
		return nil, err
	}
	return buckets, nil
}

func (s *MinioClientService) CreateBucket(ctx context.Context, bucketName string) error {
	err := s.str.CreateBucket(ctx, bucketName)
	if err != nil {
		log.Error(fmt.Sprintf("Error while getting buxkets: %s", err.Error()))
		return err
	}
	return nil
}
