package minioapi

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
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

func (s *MinioClientService) UploadStream(ctx context.Context, stream UploadFile) (*string, error) {
	stream.fileName = uuid.New().String() + "_" + stream.fileName
	err := s.str.UploadStream(ctx, stream.bucket, stream.fileName, stream.file, stream.fileSize, stream.contentType)
	if err != nil {
		log.Error("failed to upload file: ", err.Error())
		return nil, err
	}
	mt := s.str.GetClientMetadata()
	scheme := "http"
	if mt.Secure {
		scheme = "https"
	}

	publicURL := fmt.Sprintf("%s://%s/%s/%s", scheme, mt.Endpoint, stream.bucket, stream.fileName)

	return &publicURL, nil
}

func (s *MinioClientService) UploadImageString(ctx context.Context, stream UploadFile) (*string, error) {
	stream.fileName = uuid.New().String() + "_" + time.Now().Format("20060102_150405")
	err := s.str.UploadBase64Image(ctx, stream.bucket, stream.fileName, stream.Base64)
	if err != nil {
		log.Error("failed to upload file: ", err.Error())
		return nil, err
	}
	mt := s.str.GetClientMetadata()
	scheme := "http"
	if mt.Secure {
		scheme = "https"
	}

	publicURL := fmt.Sprintf("%s://%s/%s/%s", scheme, mt.Endpoint, stream.bucket, stream.fileName)

	return &publicURL, nil
}

func (s *MinioClientService) Presigned(ctx context.Context, bucket, object string) (*string, error) {
	url, err := s.str.PresignedURL(ctx, bucket, object, 2*time.Minute)
	if err != nil {
		log.Error("Error while getting presigned URL: ", err.Error())
		return nil, err
	}
	return &url, nil
}
