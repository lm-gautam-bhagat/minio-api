package storage

import (
	"context"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type StorageClient struct {
	client *minio.Client
}

func NewStorageClient(minIOEndpoint, minIOAccessID, minIOAccessKey string, useSSL bool) (*StorageClient, error) {
	c, err := minio.New(minIOEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minIOAccessID, minIOAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	return &StorageClient{client: c}, nil
}

// EnsureBucket creates a bucket if it doesn't exist
func (m *StorageClient) EnsureBucket(ctx context.Context, bucket string) error {
	exists, err := m.client.BucketExists(ctx, bucket)
	if err != nil {
		return err
	}
	if !exists {
		return m.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
	}
	return nil
}

func (m *StorageClient) UploadFile(ctx context.Context, bucket, objectName, filePath, contentType string) error {
	_, err := m.client.FPutObject(ctx, bucket, objectName, filePath,
		minio.PutObjectOptions{ContentType: contentType})
	return err
}

func (m *StorageClient) DownloadFile(ctx context.Context, bucket, objectName, destPath string) error {
	return m.client.FGetObject(ctx, bucket, objectName, destPath, minio.GetObjectOptions{})
}

func (m *StorageClient) ListObjects(ctx context.Context, bucket string) ([]string, error) {
	var objs []string
	for obj := range m.client.ListObjects(ctx, bucket, minio.ListObjectsOptions{Recursive: true}) {
		if obj.Err != nil {
			return nil, obj.Err
		}
		objs = append(objs, obj.Key)
	}
	return objs, nil
}

func (m *StorageClient) PresignedURL(ctx context.Context, bucket, objectName string, expiry time.Duration) (string, error) {
	u, err := m.client.PresignedGetObject(ctx, bucket, objectName, expiry, nil)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func (m *StorageClient) ListBuckets(ctx context.Context) ([]string, error) {
	buckets, err := m.client.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, bucket := range buckets {
		names = append(names, bucket.Name)
	}
	return names, nil
}

func (m *StorageClient) CreateBucket(ctx context.Context, bucketName string) error {
	ops := minio.MakeBucketOptions{}
	err := m.client.MakeBucket(ctx, bucketName, ops)
	if err != nil {
		return err
	}
	return nil
}
