package xelon

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

const (
	backupPlanBasePath       = "backups/tags"
	deviceBackupPlanBasePath = deviceBasePath + "/%s/" + backupPlanBasePath
)

// BackupService handles communication with the backup-related methods of the Xelon API.
type BackupService service

// BackupPlan represents a backup plan available in a Xelon cloud.
type BackupPlan struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type backupPlanListOptions struct {
	CloudID string `url:"cloudIdentifier"`
}

type setDeviceBackupPlanRequest struct {
	BackupPlanID int `json:"backupId"`
}

type backupPlansRoot struct {
	BackupPlans []BackupPlan `json:"data"`
}

func (v BackupPlan) String() string { return Stringify(v) }

// ListPlans lists backup plans available in a cloud.
func (s *BackupService) ListPlans(ctx context.Context, cloudID string) ([]BackupPlan, *Response, error) {
	if cloudID == "" {
		return nil, nil, fmt.Errorf("cloud id: %w", ErrEmptyArgument)
	}

	path, err := addOptions(backupPlanBasePath, &backupPlanListOptions{CloudID: cloudID})
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(backupPlansRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if root.BackupPlans == nil {
		return nil, resp, errors.New("failed to list backup plans: response data is empty")
	}

	return root.BackupPlans, resp, nil
}

// GetDevicePlan gets the backup plan assigned to a device. It returns a nil plan
// without an error when backups are disabled for the device.
func (s *BackupService) GetDevicePlan(ctx context.Context, deviceID string) (*BackupPlan, *Response, error) {
	if deviceID == "" {
		return nil, nil, fmt.Errorf("device id: %w", ErrEmptyArgument)
	}

	path := fmt.Sprintf(deviceBackupPlanBasePath, deviceID)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(backupPlansRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if root.BackupPlans == nil {
		return nil, resp, errors.New("failed to get device backup plan: response data is empty")
	}

	switch len(root.BackupPlans) {
	case 0:
		return nil, resp, nil
	case 1:
		return &root.BackupPlans[0], resp, nil
	default:
		return nil, resp, fmt.Errorf("%w: expected at most one plan, got %d", ErrAmbiguousBackupPlanAssignment, len(root.BackupPlans))
	}
}

// SetDevicePlan assigns a backup plan to a device, replacing any current assignment.
func (s *BackupService) SetDevicePlan(ctx context.Context, deviceID string, backupPlanID int) (*Response, error) {
	if deviceID == "" {
		return nil, fmt.Errorf("device id: %w", ErrEmptyArgument)
	}
	if backupPlanID <= 0 {
		return nil, fmt.Errorf("backup plan id: %w", ErrEmptyArgument)
	}

	path := fmt.Sprintf(deviceBackupPlanBasePath, deviceID)
	setRequest := &setDeviceBackupPlanRequest{BackupPlanID: backupPlanID}
	req, err := s.client.NewRequest(http.MethodPut, path, setRequest)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// DisableDeviceBackup disables backups for a device by removing its backup plan assignment.
func (s *BackupService) DisableDeviceBackup(ctx context.Context, deviceID string) (*Response, error) {
	if deviceID == "" {
		return nil, fmt.Errorf("device id: %w", ErrEmptyArgument)
	}

	path := fmt.Sprintf(deviceBackupPlanBasePath, deviceID)
	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
