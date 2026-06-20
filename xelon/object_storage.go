package xelon

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"net/http"
)

const objectStorageBasePath = "object-storages"

// ObjectStoragesService handles communication with the object storage related methods of the Xelon API.
type ObjectStoragesService service

// ObjectStorageUser represents
type ObjectStorageUser struct {
	ID                     string   `json:"identifier,omitempty"`
	Name                   string   `json:"name,omitempty"`
	QuotaGB                int      `json:"quota,omitempty"`
	S3Endpoints            []string `json:"s3endpoints,omitempty"`
	Tenant                 *Tenant  `json:"tenant,omitempty"`
	UsedGB                 float32  `json:"sizeUsedGb,omitempty"`
	ZoneReplicationEnabled bool     `json:"isReplicated,omitempty"`
}

type ObjectStorageUserCreateRequest struct {
	Name     string `json:"name"`
	QuotaGB  int    `json:"quota"`
	RegionID string `json:"zoneGroup"`
	TenantID string `json:"tenantId,omitempty"`
}

type ObjectStorageUserUpdateRequest struct {
	Name    string `json:"name"`
	QuotaGB int    `json:"quota"`
}

// ObjectStorageUserListOptions specifies the optional parameters to the ObjectStoragesService.ListUsers.
type ObjectStorageUserListOptions struct {
	TenantID string `url:"tenantId,omitempty"`

	ListOptions
}

type objectStorageUsersRoot struct {
	ObjectStorageUsers []ObjectStorageUser `json:"data"`
	Meta               *Meta               `json:"meta,omitempty"`
}

func (v ObjectStorageUser) String() string { return Stringify(v) }

// ListUsers provides a list of available object storage users.
func (s *ObjectStoragesService) ListUsers(ctx context.Context, opts *ObjectStorageUserListOptions) ([]ObjectStorageUser, *Response, error) {
	path, err := addOptions(fmt.Sprintf("%v/users", objectStorageBasePath), opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(objectStorageUsersRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}

	return root.ObjectStorageUsers, resp, nil
}

// AllUsers returns an iterator to paginate over all object storage users.
//
// The return iterator can be used in a for...range loop to easily process all users.
func (s *ObjectStoragesService) AllUsers(ctx context.Context, opts *ListOptions) (iter.Seq2[ObjectStorageUser, *Response], func() error) {
	return newPaginator[ObjectStorageUser](ctx, s.client, fmt.Sprintf("%v/users", objectStorageBasePath), opts)
}

// GetUser provides detailed information for object storage user identified by id.
func (s *ObjectStoragesService) GetUser(ctx context.Context, objectStorageUserID string) (*ObjectStorageUser, *Response, error) {
	if objectStorageUserID == "" {
		return nil, nil, errors.New("failed to get object storage user: id must be supplied")
	}

	path := fmt.Sprintf("%v/users/%v", objectStorageBasePath, objectStorageUserID)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	user := new(ObjectStorageUser)
	resp, err := s.client.Do(ctx, req, user)
	if err != nil {
		return nil, resp, err
	}

	return user, resp, nil
}

// CreateUser makes a new object storage user with given payload.
func (s *ObjectStoragesService) CreateUser(ctx context.Context, createRequest *ObjectStorageUserCreateRequest) (*ObjectStorageUser, *Response, error) {
	if createRequest == nil {
		return nil, nil, errors.New("failed to create object storage user: payload must be supplied")
	}

	path := fmt.Sprintf("%v/users", objectStorageBasePath)
	req, err := s.client.NewRequest(http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	user := new(ObjectStorageUser)
	resp, err := s.client.Do(ctx, req, user)
	if err != nil {
		return nil, resp, err
	}

	return user, resp, nil
}

// UpdateUser changes a object storage user.
func (s *ObjectStoragesService) UpdateUser(ctx context.Context, objectStorageUserID string, updateRequest *ObjectStorageUserUpdateRequest) (*ObjectStorageUser, *Response, error) {
	if objectStorageUserID == "" {
		return nil, nil, errors.New("failed to update object storage user: id must be supplied")
	}
	if updateRequest == nil {
		return nil, nil, errors.New("failed to update object storage user: payload must be supplied")
	}

	path := fmt.Sprintf("%v/users/%v", objectStorageBasePath, objectStorageUserID)
	req, err := s.client.NewRequest(http.MethodPut, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	user := new(ObjectStorageUser)
	resp, err := s.client.Do(ctx, req, user)
	if err != nil {
		return nil, resp, err
	}

	return user, resp, nil
}

// DeleteUser remove sobject storage user identified by id.
func (s *ObjectStoragesService) DeleteUser(ctx context.Context, objectStorageUserID string) (*Response, error) {
	if objectStorageUserID == "" {
		return nil, errors.New("failed to delete object storage user: id must be supplied")
	}

	path := fmt.Sprintf("%v/users/%v", objectStorageBasePath, objectStorageUserID)
	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// ObjectStoragePlan represents a Xelon object storage pricing plan.
type ObjectStoragePlan struct {
	Description          string  `json:"description,omitempty"`
	Name                 string  `json:"title,omitempty"`
	QuotaGB              int     `json:"quota,omitempty"`
	Price                float32 `json:"price,omitempty"`
	PriceWithReplication float32 `json:"priceReplicated,omitempty"`
}

func (v ObjectStoragePlan) String() string { return Stringify(v) }

// ListPlans provides a list of available storage plans.
func (s *ObjectStoragesService) ListPlans(ctx context.Context) ([]ObjectStoragePlan, *Response, error) {
	path := fmt.Sprintf("%v/plans", objectStorageBasePath)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var plans []ObjectStoragePlan
	resp, err := s.client.Do(ctx, req, &plans)
	if err != nil {
		return nil, resp, err
	}

	return plans, resp, nil
}

// ObjectStorageRegion represents a Xelon object storage regions (zone group).
type ObjectStorageRegion struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	ReplicationEnabled bool   `json:"hasReplication"`
}

func (v ObjectStorageRegion) String() string { return Stringify(v) }

// ListRegions provides a list of available regions.
func (s *ObjectStoragesService) ListRegions(ctx context.Context) ([]ObjectStorageRegion, *Response, error) {
	path := fmt.Sprintf("%v/zone-groups", objectStorageBasePath)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var regions []ObjectStorageRegion
	resp, err := s.client.Do(ctx, req, &regions)
	if err != nil {
		return nil, resp, err
	}

	return regions, resp, nil
}
