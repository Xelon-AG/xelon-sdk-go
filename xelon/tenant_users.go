package xelon

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

const tenantUsersBasePath = "tenants/%s/users"

// TenantUsersService handles communication with tenant user methods of the Xelon API.
type TenantUsersService service

// TenantUser represents a user that belongs to a Xelon tenant.
type TenantUser struct {
	Email       string                 `json:"email,omitempty"`
	ID          string                 `json:"identifier,omitempty"`
	JobTitle    string                 `json:"jobTitle,omitempty"`
	Name        string                 `json:"name,omitempty"`
	Permissions []TenantUserPermission `json:"permissions,omitempty"`
	Phone       string                 `json:"phone,omitempty"`
	Roles       []TenantUserRole       `json:"roles,omitempty"`
	Surname     string                 `json:"surname,omitempty"`
	TenantID    string                 `json:"tenantIdentifier,omitempty"`
}

type TenantUserRole struct {
	DisplayName string `json:"friendlyName,omitempty"`
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Type        string `json:"type,omitempty"`
}

type TenantUserCreateRequest struct {
	BusinessPhone         string   `json:"business_phone,omitempty"`
	Email                 string   `json:"email"`
	JobTitle              string   `json:"job_title,omitempty"`
	Name                  string   `json:"name"`
	Password              string   `json:"password"`
	PasswordConfirmation  string   `json:"password_confirmation"`
	Permissions           []string `json:"permissions,omitempty"`
	Phone                 string   `json:"phone,omitempty"`
	RequirePasswordChange bool     `json:"passwordShouldBeChanged"`
	Roles                 []string `json:"roles,omitempty"`
	SendWelcomeEmail      bool     `json:"welcomeEmail"`
	Surname               string   `json:"surname"`
}

type TenantUserUpdateRequest struct {
	BusinessPhone string `json:"business_phone,omitempty"`
	JobTitle      string `json:"job_title,omitempty"`
	Name          string `json:"name"`
	Phone         string `json:"phone,omitempty"`
	Surname       string `json:"surname"`
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

type tenantUsersRoot struct {
	TenantUsers []TenantUser `json:"data"`
	Meta        *Meta        `json:"meta,omitempty"`
}

func (v TenantUser) String() string { return Stringify(v) }

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

// Get gets a tenant user by id.
func (s *TenantUsersService) Get(ctx context.Context, tenantID, userID string) (*TenantUser, *Response, error) {
	if tenantID == "" {
		return nil, nil, fmt.Errorf("tenant id: %w", ErrEmptyArgument)
	}
	if userID == "" {
		return nil, nil, fmt.Errorf("user id: %w", ErrEmptyArgument)
	}

	path := fmt.Sprintf("%v/%v", fmt.Sprintf(tenantUsersBasePath, tenantID), userID)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
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

	path := fmt.Sprintf("%v/%v", fmt.Sprintf(tenantUsersBasePath, tenantID), userID)
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

// Delete deletes a tenant user by id.
func (s *TenantUsersService) Delete(ctx context.Context, tenantID, userID string) (*Response, error) {
	if tenantID == "" {
		return nil, fmt.Errorf("tenant id: %w", ErrEmptyArgument)
	}
	if userID == "" {
		return nil, fmt.Errorf("user id: %w", ErrEmptyArgument)
	}

	path := fmt.Sprintf("%v/%v", fmt.Sprintf(tenantUsersBasePath, tenantID), userID)
	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
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

	path := fmt.Sprintf("%v/permissions", fmt.Sprintf(tenantUsersBasePath, tenantID))
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

	path := fmt.Sprintf("%v/%v/permissions", fmt.Sprintf(tenantUsersBasePath, tenantID), userID)
	req, err := s.client.NewRequest(http.MethodPost, path, updateRequest)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
