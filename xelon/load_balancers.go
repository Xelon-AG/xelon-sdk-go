package xelon

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const loadBalancerBasePath = "load-balancers"

// LoadBalancersService handles communication with the load balancers related methods of the Xelon REST API.
type LoadBalancersService service

// LoadBalancer represents a Xelon load balancer.
type LoadBalancer struct {
	AssignedDevices   []LoadBalancerAssignedDevice `json:"assignedDevices,omitempty"`
	Cloud             *Cloud                       `json:"cloud,omitempty"`
	CreatedAt         *time.Time                   `json:"createdAt,omitempty"`
	ExternalIPAddress string                       `json:"externalIp,omitempty"`
	ID                string                       `json:"identifier,omitempty"`
	InternalIPAddress string                       `json:"internalIp,omitempty"`
	HealthStatus      string                       `json:"health,omitempty"`
	Name              string                       `json:"name,omitempty"`
	State             int                          `json:"state,omitempty"`
	Tenant            *Tenant                      `json:"tenant,omitempty"`
}

type LoadBalancerAssignedDevice struct {
	ID   string `json:"identifier,omitempty"`
	Name string `json:"name,omitempty"`
}

type LoadBalancerCreateRequest struct {
	CloudID           string `json:"cloudIdentifier"`
	InternalIPAddress string `json:"internalIp,omitempty"`
	InternalNetworkID string `json:"internalNetworkIdentifier"`
	Type              string `json:"loadBalancingType"` // layer4 or layer7
	Name              string `json:"name"`
	TenantID          string `json:"tenantIdentifier"`
}

type LoadBalancerUpdateRequest struct {
	Name string `json:"name"`
}

type LoadBalancerUpdateAssignedDevicesRequest struct {
	DeviceIDs []string `json:"deviceIdentifiers"`
}

// LoadBalancerListOptions specifies the optional parameters to the LoadBalancersService.List.
type LoadBalancerListOptions struct {
	Sort   string `url:"sort,omitempty"`
	Search string `url:"search,omitempty"`

	ListOptions
}

type loadBalancerRoot struct {
	LoadBalancer *LoadBalancer `json:"data,omitempty"`
	Message      string        `json:"message,omitempty"`
}

type loadBalancerAssignedDevicesRoot struct {
	AssignedDevices []LoadBalancerAssignedDevice `json:"data,omitempty"`
}

type loadBalancersRoot struct {
	LoadBalancers []LoadBalancer `json:"data"`
	Meta          *Meta          `json:"meta,omitempty"`
}

func (v LoadBalancer) String() string { return Stringify(v) }

// List provides a list of all load balancers.
func (s *LoadBalancersService) List(ctx context.Context, opts *LoadBalancerListOptions) ([]LoadBalancer, *Response, error) {
	path, err := addOptions(loadBalancerBasePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(loadBalancersRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}

	return root.LoadBalancers, resp, nil
}

// Get provides detailed information for load balancer identified by id.
func (s *LoadBalancersService) Get(ctx context.Context, loadBalancerID string) (*LoadBalancer, *Response, error) {
	if loadBalancerID == "" {
		return nil, nil, errors.New("failed to get load balancer: id must be supplied")
	}

	path := fmt.Sprintf("%v/%v", loadBalancerBasePath, loadBalancerID)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	loadBalancer := new(LoadBalancer)
	resp, err := s.client.Do(ctx, req, loadBalancer)
	if err != nil {
		return nil, resp, err
	}

	return loadBalancer, resp, err
}

// Create makes a load balancer with given payload.
func (s *LoadBalancersService) Create(ctx context.Context, createRequest *LoadBalancerCreateRequest) (*LoadBalancer, *Response, error) {
	if createRequest == nil {
		return nil, nil, errors.New("failed to create load balancer: payload must be supplied")
	}

	req, err := s.client.NewRequest(http.MethodPost, loadBalancerBasePath, createRequest)
	if err != nil {
		return nil, nil, err
	}

	loadBalancerRoot := new(loadBalancerRoot)
	resp, err := s.client.Do(ctx, req, loadBalancerRoot)
	if err != nil {
		return nil, resp, err
	}

	return loadBalancerRoot.LoadBalancer, resp, nil
}

// Update changes load balancer identified by id.
func (s *LoadBalancersService) Update(ctx context.Context, loadBalancerID string, updateRequest *LoadBalancerUpdateRequest) (*LoadBalancer, *Response, error) {
	if loadBalancerID == "" {
		return nil, nil, errors.New("failed to update load balancer: id must be supplied")
	}
	if updateRequest == nil {
		return nil, nil, errors.New("failed to update load balancer: payload must be supplied")
	}

	path := fmt.Sprintf("%v/%v", loadBalancerBasePath, loadBalancerID)
	req, err := s.client.NewRequest(http.MethodPut, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	loadBalancerRoot := new(loadBalancerRoot)
	resp, err := s.client.Do(ctx, req, loadBalancerRoot)
	if err != nil {
		return nil, resp, err
	}

	return loadBalancerRoot.LoadBalancer, resp, nil
}

// Delete removes load balancer identified by id.
func (s *LoadBalancersService) Delete(ctx context.Context, loadBalancerID string) (*Response, error) {
	if loadBalancerID == "" {
		return nil, errors.New("failed to delete load balancer: id must be supplied")
	}

	path := fmt.Sprintf("%v/%v", loadBalancerBasePath, loadBalancerID)
	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// ListAssignedDevices provides information about assigned device.
func (s *LoadBalancersService) ListAssignedDevices(ctx context.Context, loadBalancerID, networkID string) ([]LoadBalancerAssignedDevice, *Response, error) {
	if loadBalancerID == "" {
		return nil, nil, errors.New("failed to list assigned devices for load balancer: id must be supplied")
	}
	if networkID == "" {
		return nil, nil, errors.New("failed to list assigned devices for load balancer: network id must be supplied")
	}

	path := fmt.Sprintf("%v/%v/assignable-devices/%v", loadBalancerBasePath, loadBalancerID, networkID)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(loadBalancerAssignedDevicesRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.AssignedDevices, resp, nil
}

// UpdateAssignedDevices changes assigned devices.
func (s *LoadBalancersService) UpdateAssignedDevices(ctx context.Context, loadBalancerID string, updateRequest *LoadBalancerUpdateAssignedDevicesRequest) (*Response, error) {
	if loadBalancerID == "" {
		return nil, errors.New("failed to update assigned devices for load balancer: id must be supplied")
	}
	if updateRequest == nil {
		return nil, errors.New("failed to update assigned devices for load balancer: payload must be supplied")
	}

	path := fmt.Sprintf("%v/%v/assigned-devices", loadBalancerBasePath, loadBalancerID)
	req, err := s.client.NewRequest(http.MethodPut, path, updateRequest)
	if err != nil {
		return nil, err
	}

	root := new(loadBalancerRoot)
	return s.client.Do(ctx, req, root)
}
