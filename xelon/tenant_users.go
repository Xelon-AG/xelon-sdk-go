package xelon

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"net/http"
	"time"
)

const tenantUsersBasePath = "tenants/%s/users"

// TenantUsersService handles communication with tenant user methods of the Xelon API.
type TenantUsersService service

// TenantUser represents a user that belongs to a Xelon tenant.
type TenantUser struct {
	BusinessPhone string     `json:"businessPhone,omitempty"`
	DeletedAt     *time.Time `json:"deletedAt,omitempty"`
	Email         string     `json:"email,omitempty"`
	FirstName     string     `json:"name,omitempty"`
	ID            string     `json:"identifier,omitempty"`
	IsActive      bool       `json:"isActive"`
	JobTitle      string     `json:"jobTitle,omitempty"`
	LastName      string     `json:"surname,omitempty"`
	Phone         string     `json:"phone,omitempty"`
	TenantID      string     `json:"tenantIdentifier,omitempty"`
}

// TenantUserWithDetails represents a tenant user with roles, permissions.
type TenantUserWithDetails struct {
	TenantUser

	Permissions []TenantUserPermission `json:"permissions,omitempty"`
	Roles       []TenantUserRole       `json:"roles,omitempty"`
}

type TenantUserRole struct {
	DisplayName string `json:"friendlyName,omitempty"`
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Type        string `json:"type,omitempty"`
}

type TenantUserCreateRequest struct {
	BusinessPhone         string   `json:"businessPhone,omitempty"`
	Email                 string   `json:"email"`
	FirstName             string   `json:"name"`
	JobTitle              string   `json:"jobTitle,omitempty"`
	LastName              string   `json:"surname"`
	Password              string   `json:"password"`
	PasswordConfirmation  string   `json:"passwordConfirmation"`
	Permissions           []string `json:"permissions,omitempty"`
	Phone                 string   `json:"phone,omitempty"`
	RequirePasswordChange bool     `json:"passwordShouldBeChanged"`
	Roles                 []string `json:"roles,omitempty"`
	SendWelcomeEmail      bool     `json:"welcomeEmail"`
}

type TenantUserUpdateRequest struct {
	BusinessPhone string `json:"businessPhone,omitempty"`
	FirstName     string `json:"name"`
	JobTitle      string `json:"jobTitle,omitempty"`
	LastName      string `json:"surname"`
	Phone         string `json:"phone,omitempty"`
}

type TenantUserPasswordUpdateRequest struct {
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

// TenantUserListOptions specifies the optional parameters to the TenantUsersService.List.
type TenantUserListOptions struct {
	Search string `url:"search,omitempty"`
	Sort   string `url:"sort,omitempty"`

	ListOptions
}

type tenantUserRoot struct {
	TenantUser *TenantUser `json:"data,omitempty"`
}

type tenantUserWithDetailsRoot struct {
	TenantUser *TenantUserWithDetails `json:"data,omitempty"`
}

type tenantUsersRoot struct {
	TenantUsers []TenantUser `json:"data"`
	Meta        *Meta        `json:"meta,omitempty"`
}

func (v TenantUser) String() string { return Stringify(v) }

func (v TenantUserWithDetails) String() string { return Stringify(v) }

func (v TenantUserRole) String() string { return Stringify(v) }

// List lists users that belong to a tenant.
func (s *TenantUsersService) List(ctx context.Context, tenantID string, opts *TenantUserListOptions) ([]TenantUser, *Response, error) {
	if tenantID == "" {
		return nil, nil, fmt.Errorf("tenant id: %w", ErrEmptyArgument)
	}

	path, err := addOptions(fmt.Sprintf(tenantUsersBasePath, tenantID), opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(tenantUsersRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}

	return root.TenantUsers, resp, nil
}

// All returns an iterator to paginate over all users that belong to a tenant.
//
// The return iterator can be used in a for...range loop to easily process all users.
func (s *TenantUsersService) All(ctx context.Context, tenantID string, opts *ListOptions) (iter.Seq2[TenantUser, *Response], func() error) {
	if tenantID == "" {
		err := fmt.Errorf("tenant id: %w", ErrEmptyArgument)
		return func(yield func(TenantUser, *Response) bool) {}, func() error { return err }
	}

	return newPaginator[TenantUser](ctx, s.client, fmt.Sprintf(tenantUsersBasePath, tenantID), opts)
}

// Get gets a tenant user by id, including detailed roles, permissions, and active state.
func (s *TenantUsersService) Get(ctx context.Context, tenantID, userID string) (*TenantUserWithDetails, *Response, error) {
	if tenantID == "" {
		return nil, nil, fmt.Errorf("tenant id: %w", ErrEmptyArgument)
	}
	if userID == "" {
		return nil, nil, fmt.Errorf("user id: %w", ErrEmptyArgument)
	}

	path := fmt.Sprintf(tenantUsersBasePath+"/%s", tenantID, userID)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(tenantUserWithDetailsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if root.TenantUser == nil {
		return nil, resp, errors.New("tenant user data is empty")
	}

	return root.TenantUser, resp, nil
}

// Create creates a tenant user.
func (s *TenantUsersService) Create(ctx context.Context, tenantID string, createRequest *TenantUserCreateRequest) (*TenantUser, *Response, error) {
	if tenantID == "" {
		return nil, nil, fmt.Errorf("tenant id: %w", ErrEmptyArgument)
	}
	if createRequest == nil {
		return nil, nil, fmt.Errorf("payload: %w", ErrEmptyPayloadNotAllowed)
	}

	req, err := s.client.NewRequest(http.MethodPost, fmt.Sprintf(tenantUsersBasePath, tenantID), createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(tenantUserRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if root.TenantUser == nil {
		return nil, resp, errors.New("tenant user data is empty")
	}

	return root.TenantUser, resp, nil
}

// Update updates a tenant user by id.
func (s *TenantUsersService) Update(ctx context.Context, tenantID, userID string, updateRequest *TenantUserUpdateRequest) (*TenantUser, *Response, error) {
	if tenantID == "" {
		return nil, nil, fmt.Errorf("tenant id: %w", ErrEmptyArgument)
	}
	if userID == "" {
		return nil, nil, fmt.Errorf("user id: %w", ErrEmptyArgument)
	}
	if updateRequest == nil {
		return nil, nil, fmt.Errorf("payload: %w", ErrEmptyPayloadNotAllowed)
	}

	path := fmt.Sprintf(tenantUsersBasePath+"/%s", tenantID, userID)
	req, err := s.client.NewRequest(http.MethodPut, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(tenantUserRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if root.TenantUser == nil {
		return nil, resp, errors.New("tenant user data is empty")
	}

	return root.TenantUser, resp, nil
}

// Delete soft-deletes a tenant user by id.
func (s *TenantUsersService) Delete(ctx context.Context, tenantID, userID string) (*Response, error) {
	if tenantID == "" {
		return nil, fmt.Errorf("tenant id: %w", ErrEmptyArgument)
	}
	if userID == "" {
		return nil, fmt.Errorf("user id: %w", ErrEmptyArgument)
	}

	path := fmt.Sprintf(tenantUsersBasePath+"/%s", tenantID, userID)
	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Restore restores a previously deleted tenant user.
func (s *TenantUsersService) Restore(ctx context.Context, tenantID, userID string) (*Response, error) {
	if tenantID == "" {
		return nil, fmt.Errorf("tenant id: %w", ErrEmptyArgument)
	}
	if userID == "" {
		return nil, fmt.Errorf("user id: %w", ErrEmptyArgument)
	}

	restoreRequest := struct {
		UserID string `json:"userIdentifier"`
	}{
		UserID: userID,
	}
	path := fmt.Sprintf(tenantUsersBasePath+"/restore", tenantID)
	req, err := s.client.NewRequest(http.MethodPost, path, restoreRequest)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// UpdatePassword updates a tenant user's password.
func (s *TenantUsersService) UpdatePassword(ctx context.Context, tenantID, userID string, updateRequest *TenantUserPasswordUpdateRequest) (*Response, error) {
	if tenantID == "" {
		return nil, fmt.Errorf("tenant id: %w", ErrEmptyArgument)
	}
	if userID == "" {
		return nil, fmt.Errorf("user id: %w", ErrEmptyArgument)
	}
	if updateRequest == nil {
		return nil, fmt.Errorf("payload: %w", ErrEmptyPayloadNotAllowed)
	}

	path := fmt.Sprintf(tenantUsersBasePath+"/%s/password", tenantID, userID)
	req, err := s.client.NewRequest(http.MethodPost, path, updateRequest)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

type TenantUserPermission struct {
	DisplayName string `json:"friendlyName,omitempty"`
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Type        string `json:"type,omitempty"`
}

type TenantUserPermissionsUpdateRequest struct {
	ChildTenants []string `json:"childTenants,omitempty"`
	Permissions  []string `json:"permissions"`
	Roles        []string `json:"roles"`
}

func (v TenantUserPermission) String() string { return Stringify(v) }

// ListAvailablePermissions lists available permissions that can be assigned to users within a tenant.
func (s *TenantUsersService) ListAvailablePermissions(ctx context.Context, tenantID string) ([]TenantUserPermission, *Response, error) {
	if tenantID == "" {
		return nil, nil, fmt.Errorf("tenant id: %w", ErrEmptyArgument)
	}

	path := fmt.Sprintf(tenantUsersBasePath+"/permissions", tenantID)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var permissions []TenantUserPermission
	resp, err := s.client.Do(ctx, req, &permissions)
	if err != nil {
		return nil, resp, err
	}

	return permissions, resp, nil
}

// UpdatePermissions updates roles and permissions for a tenant user.
func (s *TenantUsersService) UpdatePermissions(ctx context.Context, tenantID, userID string, updateRequest *TenantUserPermissionsUpdateRequest) (*Response, error) {
	if tenantID == "" {
		return nil, fmt.Errorf("tenant id: %w", ErrEmptyArgument)
	}
	if userID == "" {
		return nil, fmt.Errorf("user id: %w", ErrEmptyArgument)
	}
	if updateRequest == nil {
		return nil, fmt.Errorf("payload: %w", ErrEmptyPayloadNotAllowed)
	}

	path := fmt.Sprintf(tenantUsersBasePath+"/%s/permissions", tenantID, userID)
	req, err := s.client.NewRequest(http.MethodPost, path, updateRequest)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
