package storage

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type StorageClient struct {
	client   *minio.Client
	endpoint string
	secure   bool
}

type ClientMetadata struct {
	Endpoint string
	Secure   bool
}

func NewStorageClient(minIOEndpoint, minIOAccessID, minIOAccessKey string, useSSL bool) (*StorageClient, error) {
	c, err := minio.New(minIOEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minIOAccessID, minIOAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	return &StorageClient{client: c, endpoint: minIOEndpoint, secure: useSSL}, nil
}

func (s *StorageClient) GetClientMetadata() ClientMetadata {
	return ClientMetadata{
		Endpoint: s.endpoint,
		Secure:   s.secure,
	}
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

func (m *StorageClient) UploadStream(
	ctx context.Context,
	bucket, objectName string,
	file io.Reader,
	objectSize int64,
	contentType string,
) error {
	// ensure bucket exists before upload
	if err := m.EnsureBucket(ctx, bucket); err != nil {
		return fmt.Errorf("failed to ensure bucket: %w", err)
	}
	_, err := m.client.PutObject(
		ctx,
		bucket,
		objectName,
		file,
		objectSize,
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return fmt.Errorf("failed to upload object: %w", err)
	}

	return nil
}

func (s *StorageClient) UploadBase64Image(ctx context.Context, bucketName, objectName, base64Image string) error {
	var contentType string
	if idx := strings.Index(base64Image, ","); idx != -1 {
		prefix := base64Image[:idx]
		if strings.Contains(prefix, "data:") && strings.Contains(prefix, ";base64") {
			contentType = strings.TrimPrefix(strings.Split(prefix, ";")[0], "data:")
		}
		base64Image = base64Image[idx+1:]
	}

	imageData, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		return fmt.Errorf("failed to decode base64 image: %w", err)
	}

	if contentType == "" {
		contentType = http.DetectContentType(imageData)
	}

	reader := bytes.NewReader(imageData)

	exists, err := s.client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket: %w", err)
	}
	if !exists {
		if err := s.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	// Upload
	_, err = s.client.PutObject(
		ctx,
		bucketName,
		objectName,
		reader,
		int64(len(imageData)),
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to upload image: %w", err)
	}

	return nil
}

func (a *StorageClient) SetBucketPublicAccess(ctx context.Context, bucketName string, makePublic bool) error {
	var policy string

	if makePublic {
		// Public Read-only Policy
		policy = `{
			"Version":"2012-10-17",
			"Statement":[{
				"Effect":"Allow",
				"Principal":"*",
				"Action":["s3:GetObject"],
				"Resource":["arn:aws:s3:::` + bucketName + `/*"]
			}]
		}`
	} else {
		policy = ``
	}
	return a.client.SetBucketPolicy(ctx, bucketName, policy)
}
