package storage

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lm-gautam-bhagat/minio-server/constants"
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

func (a *StorageAdmin) CreateUser(ctx context.Context, accessID, secretKey string) error {
	err := a.admin.AddUser(ctx, accessID, secretKey)
	if err != nil {
		return err
	}
	err = a.admin.SetUserStatus(ctx, accessID, madmin.AccountEnabled)
	if err != nil {
		return err
	}
	return nil
}

// DeleteUser removes a user
func (a *StorageAdmin) DeleteUser(ctx context.Context, accessID string) error {
	return a.admin.RemoveUser(ctx, accessID)
}

// EnableUser sets user status to enabled
func (a *StorageAdmin) EnableUser(ctx context.Context, accessID string) error {
	return a.admin.SetUserStatus(ctx, accessID, madmin.AccountEnabled)
}

// DisableUser sets user status to disabled
func (a *StorageAdmin) DisableUser(ctx context.Context, accessID string) error {
	return a.admin.SetUserStatus(ctx, accessID, madmin.AccountDisabled)
}

// ListUsers retrieves all users and their status
func (a *StorageAdmin) ListUsers(ctx context.Context) (map[string]madmin.UserInfo, error) {
	return a.admin.ListUsers(ctx)
}

// GetUserInfo fetches details about a specific user
func (a *StorageAdmin) GetUserInfo(ctx context.Context, accessID string) (*madmin.UserInfo, error) {
	info, err := a.admin.GetUserInfo(ctx, accessID)
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

func (a *StorageAdmin) GetAdmin() *madmin.AdminClient {
	return a.admin
}

func (a *StorageAdmin) AddCannedPolicy(ctx context.Context, policyName constants.PolicyType, policyDoc []byte) error {
	err := a.admin.AddCannedPolicy(ctx, string(policyName), policyDoc)
	if err != nil {
		return err
	}
	return nil
}

func (a *StorageAdmin) AttachPolicy(ctx context.Context, accessID string, policies []string) (madmin.PolicyAssociationResp, error) {
	resp, err := a.admin.AttachPolicy(ctx, madmin.PolicyAssociationReq{
		User:     accessID,
		Policies: policies,
	})

	return resp, err
}

func (a *StorageAdmin) ListPolicies(ctx context.Context) (map[string]json.RawMessage, error) {
	policies, err := a.admin.ListCannedPolicies(ctx)
	return policies, err
}

func (a *StorageAdmin) GetPolicy(ctx context.Context, policyName string) (*madmin.PolicyInfo, error) {
	return a.admin.InfoCannedPolicy(ctx, policyName)
}
