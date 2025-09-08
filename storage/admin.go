package storage

import (
	"context"
	"fmt"

	"github.com/minio/madmin-go/v4"
)

type StorageAdmin struct {
	admin *madmin.AdminClient
}

func NewStorageAdmin(minIOEndpoint, minIOAccessID, minIOAccessKey string, useSSL bool) (*StorageAdmin, error) {
	c, err := madmin.New(minIOEndpoint, minIOAccessID, minIOAccessKey, useSSL)
	// minIOAccessID, minIOAccessKey, useSSL)
	if err != nil {
		return nil, err
	}
	return &StorageAdmin{admin: c}, nil
}

func (a *StorageAdmin) CreateUser(ctx context.Context, accessKey, secretKey string) error {
	err := a.admin.AddUser(ctx, accessKey, secretKey)
	if err != nil {
		return err
	}
	err = a.admin.SetUserStatus(ctx, accessKey, madmin.AccountEnabled)
	if err != nil {
		return err
	}
	return nil
}

// DeleteUser removes a user
func (a *StorageAdmin) DeleteUser(ctx context.Context, accessKey string) error {
	return a.admin.RemoveUser(ctx, accessKey)
}

// EnableUser sets user status to enabled
func (a *StorageAdmin) EnableUser(ctx context.Context, accessKey string) error {
	return a.admin.SetUserStatus(ctx, accessKey, madmin.AccountEnabled)
}

// DisableUser sets user status to disabled
func (a *StorageAdmin) DisableUser(ctx context.Context, accessKey string) error {
	return a.admin.SetUserStatus(ctx, accessKey, madmin.AccountDisabled)
}

// ListUsers retrieves all users and their status
func (a *StorageAdmin) ListUsers(ctx context.Context) (map[string]madmin.UserInfo, error) {
	return a.admin.ListUsers(ctx)
}

// GetUserInfo fetches details about a specific user
func (a *StorageAdmin) GetUserInfo(ctx context.Context, accessKey string) (*madmin.UserInfo, error) {
	info, err := a.admin.GetUserInfo(ctx, accessKey)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

// HealthCheck pings the MinIO server to check status
func (a *StorageAdmin) HealthCheck(ctx context.Context) error {
	_, err := a.admin.ServerInfo(ctx)
	if err != nil {
		return fmt.Errorf("server not healthy: %w", err)
	}
	return nil
}
