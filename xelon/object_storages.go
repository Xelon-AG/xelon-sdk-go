package xelon

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"net/http"
	"time"
)

const objectStorageBasePath = "object-storages"

// ObjectStoragesService handles communication with the object storage related methods of the Xelon API.
type ObjectStoragesService service

// ObjectStorageUser represents a Xelon users for S3-compatible object storage.
type ObjectStorageUser struct {
	ID                       string                   `json:"identifier,omitempty"`
	Name                     string                   `json:"name,omitempty"`
	QuotaGB                  int                      `json:"quota,omitempty"`
	RegionReplicationEnabled bool                     `json:"isReplicated,omitempty"`
	S3Endpoints              []string                 `json:"s3endpoints,omitempty"`
	Tenant                   *Tenant                  `json:"tenant,omitempty"`
	Tokens                   []ObjectStorageUserToken `json:"tokens,omitempty"`
	UsedGB                   float32                  `json:"sizeUsedGb,omitempty"`
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

type objectStorageUserRoot struct {
	ObjectStorageUser *ObjectStorageUser `json:"data"`
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

	root := new(objectStorageUserRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if root.ObjectStorageUser == nil {
		return nil, resp, errors.New("failed to create object storage user: response data is empty")
	}

	return root.ObjectStorageUser, resp, nil
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

	root := new(objectStorageUserRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if root.ObjectStorageUser == nil {
		return nil, resp, errors.New("failed to update object storage user: response data is empty")
	}

	return root.ObjectStorageUser, resp, nil
}

// DeleteUser removes object storage user identified by id.
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

// ObjectStorageUserToken represents a Xelon object storage user token.
//
// Note! SecretKey is available only after calling ObjectStoragesService.CreateUserToken once.
type ObjectStorageUserToken struct {
	AccessKey string     `json:"accessKey"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	ID        string     `json:"identifier"`
	SecretKey string     `json:"secretKey,omitempty"`
}

type objectStorageUserTokenRoot struct {
	ObjectStorageUserToken *ObjectStorageUserToken `json:"data"`
}

func (v ObjectStorageUserToken) String() string { return Stringify(v) }

// CreateUserToken makes a new object storage user token.
func (s *ObjectStoragesService) CreateUserToken(ctx context.Context, objectStorageUserID string) (*ObjectStorageUserToken, *Response, error) {
	if objectStorageUserID == "" {
		return nil, nil, errors.New("failed to create user token: object storage user id must be supplied")
	}

	path := fmt.Sprintf("%v/users/%v/tokens", objectStorageBasePath, objectStorageUserID)
	req, err := s.client.NewRequest(http.MethodPost, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(objectStorageUserTokenRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if root.ObjectStorageUserToken == nil {
		return nil, resp, errors.New("failed to create user token: response data is empty")
	}

	return root.ObjectStorageUserToken, resp, nil
}

// DeleteUserToken removes object storage user token identified by id.
func (s *ObjectStoragesService) DeleteUserToken(ctx context.Context, objectStorageUserID, objectStorageUserTokenID string) (*Response, error) {
	if objectStorageUserID == "" {
		return nil, errors.New("failed to delete user token: object storage user id must be supplied")
	}
	if objectStorageUserTokenID == "" {
		return nil, errors.New("failed to delete user token: id must be supplied")
	}

	path := fmt.Sprintf("%v/users/%v/tokens/%v", objectStorageBasePath, objectStorageUserID, objectStorageUserTokenID)
	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// ObjectStorageBucket represents a Xelon bucket for S3-compatible object storage.
type ObjectStorageBucket struct {
	CreatedAt                *time.Time `json:"createdAt,omitempty"`
	ID                       string     `json:"identifier"`
	IPRestrictionsEnabled    bool       `json:"isWithRestrictedIps,omitempty"`
	Name                     string     `json:"name,omitempty"`
	ObjectLockEnabled        bool       `json:"isObjectLock,omitempty"`
	ObjectLockRetentionDays  int        `json:"retentionPeriodDays,omitempty"`
	ObjectStorageUserID      string     `json:"s3UserIdentifier,omitempty"`
	ObjectStorageUserName    string     `json:"s3UserName,omitempty"`
	RegionName               string     `json:"zoneGroup,omitempty"`
	RegionReplicationEnabled bool       `json:"isReplicated,omitempty"`
	S3Endpoints              []string   `json:"s3endpoints,omitempty"`
	Tenant                   *Tenant    `json:"tenant,omitempty"`
	VersioningEnabled        bool       `json:"isVersioning,omitempty"`
}

type ObjectStorageBucketCreateRequest struct {
	Name                    string `json:"name"`
	ObjectLockEnabled       bool   `json:"isObjectLock"`
	ObjectLockRetentionDays int    `json:"retentionPeriodDays,omitempty"`
	ObjectStorageUserID     string `json:"s3UserIdentifier"`
	VersioningEnabled       bool   `json:"isVersioning"`
}

type ObjectStorageBucketUpdateRequest struct {
	Name string `json:"name"`
}

// ObjectStorageBucketListOptions specifies the optional parameters to the ObjectStoragesService.ListBuckets.
type ObjectStorageBucketListOptions struct {
	TenantID string `url:"tenantId,omitempty"`

	ListOptions
}

type ObjectStorageBucketVersioningUpdateRequest struct {
	VersioningEnabled bool `json:"isVersioning"`
}

type ObjectStorageBucketIPRestrictions struct {
	AllowedIPs []string `json:"allowedIps,omitempty"`
	Enabled    bool     `json:"enabled,omitempty"`
}

type ObjectStorageBucketIPRestrictionsUpdateRequest struct {
	AllowedIPs []string `json:"allowedIps,omitempty"`
	Enabled    bool     `json:"enabled"`
}

type objectStorageBucketRoot struct {
	ObjectStorageBucket *ObjectStorageBucket `json:"data,omitempty"`
}

type objectStorageBucketsRoot struct {
	ObjectStorageBuckets []ObjectStorageBucket `json:"data"`
	Meta                 *Meta                 `json:"meta,omitempty"`
}

type objectStorageBucketIPRestrictionsRoot struct {
	ObjectStorageBucketIPRestrictions *ObjectStorageBucketIPRestrictions `json:"data,omitempty"`
}

func (v ObjectStorageBucket) String() string { return Stringify(v) }

func (v ObjectStorageBucketIPRestrictions) String() string { return Stringify(v) }

// ListBuckets provides a list of available object storage buckets.
func (s *ObjectStoragesService) ListBuckets(ctx context.Context, opts *ObjectStorageBucketListOptions) ([]ObjectStorageBucket, *Response, error) {
	path, err := addOptions(fmt.Sprintf("%v/buckets", objectStorageBasePath), opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(objectStorageBucketsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}

	return root.ObjectStorageBuckets, resp, nil
}

// AllBuckets returns an iterator to paginate over all object storage buckets.
//
// The return iterator can be used in a for...range loop to easily process all buckets.
func (s *ObjectStoragesService) AllBuckets(ctx context.Context, opts *ListOptions) (iter.Seq2[ObjectStorageBucket, *Response], func() error) {
	return newPaginator[ObjectStorageBucket](ctx, s.client, fmt.Sprintf("%v/buckets", objectStorageBasePath), opts)
}

// GetBucket provides detailed information for object storage bucket identified by name and user id.
func (s *ObjectStoragesService) GetBucket(ctx context.Context, bucketName, objectStorageUserID string) (*ObjectStorageBucket, *Response, error) {
	if bucketName == "" {
		return nil, nil, fmt.Errorf("bucket name: %w", ErrEmptyArgument)
	}
	if objectStorageUserID == "" {
		return nil, nil, fmt.Errorf("object storage user id: %w", ErrEmptyArgument)
	}

	path := fmt.Sprintf("%v/buckets/%v/%v", objectStorageBasePath, bucketName, objectStorageUserID)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(objectStorageBucketRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if root.ObjectStorageBucket == nil {
		return nil, resp, errors.New("object storage bucket data is empty")
	}

	return root.ObjectStorageBucket, resp, nil
}

// CreateBucket makes a new object storage bucket with given payload.
func (s *ObjectStoragesService) CreateBucket(ctx context.Context, createRequest *ObjectStorageBucketCreateRequest) (*ObjectStorageBucket, *Response, error) {
	if createRequest == nil {
		return nil, nil, fmt.Errorf("payload: %w", ErrEmptyPayloadNotAllowed)
	}

	path := fmt.Sprintf("%v/buckets", objectStorageBasePath)
	req, err := s.client.NewRequest(http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(objectStorageBucketRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if root.ObjectStorageBucket == nil {
		return nil, resp, errors.New("object storage bucket data is empty")
	}

	return root.ObjectStorageBucket, resp, nil
}

// UpdateBucket changes an object storage bucket identified by id.
func (s *ObjectStoragesService) UpdateBucket(ctx context.Context, bucketID string, updateRequest *ObjectStorageBucketUpdateRequest) (*Response, error) {
	if bucketID == "" {
		return nil, fmt.Errorf("bucket id: %w", ErrEmptyArgument)
	}
	if updateRequest == nil {
		return nil, fmt.Errorf("payload: %w", ErrEmptyPayloadNotAllowed)
	}

	path := fmt.Sprintf("%v/buckets/%v", objectStorageBasePath, bucketID)
	req, err := s.client.NewRequest(http.MethodPut, path, updateRequest)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// DeleteBucket removes object storage bucket identified by name and user id.
func (s *ObjectStoragesService) DeleteBucket(ctx context.Context, bucketName, objectStorageUserID string) (*Response, error) {
	if bucketName == "" {
		return nil, fmt.Errorf("bucket name: %w", ErrEmptyArgument)
	}
	if objectStorageUserID == "" {
		return nil, fmt.Errorf("object storage user id: %w", ErrEmptyArgument)
	}

	path := fmt.Sprintf("%v/buckets/%v/%v", objectStorageBasePath, bucketName, objectStorageUserID)
	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// UpdateBucketVersioning changes versioning for an object storage bucket.
func (s *ObjectStoragesService) UpdateBucketVersioning(ctx context.Context, bucketName, objectStorageUserID string, updateRequest *ObjectStorageBucketVersioningUpdateRequest) (*Response, error) {
	if bucketName == "" {
		return nil, fmt.Errorf("bucket name: %w", ErrEmptyArgument)
	}
	if objectStorageUserID == "" {
		return nil, fmt.Errorf("object storage user id: %w", ErrEmptyArgument)
	}
	if updateRequest == nil {
		return nil, fmt.Errorf("payload: %w", ErrEmptyPayloadNotAllowed)
	}

	path := fmt.Sprintf("%v/buckets/%v/%v/versioning", objectStorageBasePath, bucketName, objectStorageUserID)
	req, err := s.client.NewRequest(http.MethodPut, path, updateRequest)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// GetBucketIPRestrictions provides IP restriction settings for an object storage bucket.
func (s *ObjectStoragesService) GetBucketIPRestrictions(ctx context.Context, bucketName, objectStorageUserID string) (*ObjectStorageBucketIPRestrictions, *Response, error) {
	if bucketName == "" {
		return nil, nil, fmt.Errorf("bucket name: %w", ErrEmptyArgument)
	}
	if objectStorageUserID == "" {
		return nil, nil, fmt.Errorf("object storage user id: %w", ErrEmptyArgument)
	}

	path := fmt.Sprintf("%v/buckets/%v/%v/ip-restrictions", objectStorageBasePath, bucketName, objectStorageUserID)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(objectStorageBucketIPRestrictionsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if root.ObjectStorageBucketIPRestrictions == nil {
		return nil, resp, errors.New("object storage bucket ip restrictions data is empty")
	}

	return root.ObjectStorageBucketIPRestrictions, resp, nil
}

// UpdateBucketIPRestrictions changes IP restriction settings for an object storage bucket.
func (s *ObjectStoragesService) UpdateBucketIPRestrictions(ctx context.Context, bucketName, objectStorageUserID string, updateRequest *ObjectStorageBucketIPRestrictionsUpdateRequest) (*Response, error) {
	if bucketName == "" {
		return nil, fmt.Errorf("bucket name: %w", ErrEmptyArgument)
	}
	if objectStorageUserID == "" {
		return nil, fmt.Errorf("object storage user id: %w", ErrEmptyArgument)
	}
	if updateRequest == nil {
		return nil, fmt.Errorf("payload: %w", ErrEmptyPayloadNotAllowed)
	}

	path := fmt.Sprintf("%v/buckets/%v/%v/ip-restrictions", objectStorageBasePath, bucketName, objectStorageUserID)
	req, err := s.client.NewRequest(http.MethodPut, path, updateRequest)
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

// ObjectStorageRegion represents a Xelon object storage region.
type ObjectStorageRegion struct {
	RegionID                 string `json:"id"`
	RegionName               string `json:"name"`
	RegionReplicationEnabled bool   `json:"hasReplication"`
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
