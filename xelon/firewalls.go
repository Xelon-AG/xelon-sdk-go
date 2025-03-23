package xelon

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const firewallBasePath = "firewalls"

// FirewallsService handles communication with the firewalls related methods of the Xelon REST API.
type FirewallsService service

// Firewall represents a Xelon firewall.
type Firewall struct {
	Cloud             *Cloud     `json:"cloud,omitempty"`
	CreatedAt         *time.Time `json:"createdAt,omitempty"`
	ExternalIPAddress string     `json:"externalIp,omitempty"`
	HealthStatus      string     `json:"health,omitempty"`
	ID                string     `json:"identifier,omitempty"`
	InternalIPAddress string     `json:"internalIp,omitempty"`
	Name              string     `json:"name,omitempty"`
	State             int        `json:"state,omitempty"`
	Tenant            *Tenant    `json:"tenant,omitempty"`
}

type FirewallCreateRequest struct {
	CloudID           string `json:"cloudIdentifier"`
	InternalIPAddress string `json:"internalIp,omitempty"`
	InternalNetworkID string `json:"internalNetworkIdentifier"`
	Name              string `json:"name"`
	TenantID          string `json:"tenantIdentifier"`
}

type FirewallUpdateRequest struct {
	Name string `json:"name"`
}

// FirewallListOptions specifies the optional parameters to the FirewallsService.List.
type FirewallListOptions struct {
	Sort   string `url:"sort,omitempty"`
	Search string `url:"search,omitempty"`

	ListOptions
}

type firewallRoot struct {
	Firewall *Firewall `json:"data,omitempty"`
	Message  string    `json:"message,omitempty"`
}

type firewallsRoot struct {
	Firewalls []Firewall `json:"data"`
	Meta      *Meta      `json:"meta,omitempty"`
}

func (v Firewall) String() string { return Stringify(v) }

// List provides a list of all firewalls.
func (s *FirewallsService) List(ctx context.Context, opts *FirewallListOptions) ([]Firewall, *Response, error) {
	path, err := addOptions(firewallBasePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(firewallsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}

	return root.Firewalls, resp, nil
}

// Get provides detailed information for firewall identified by id.
func (s *FirewallsService) Get(ctx context.Context, firewallID string) (*Firewall, *Response, error) {
	if firewallID == "" {
		return nil, nil, errors.New("failed to get firewall: id must be supplied")
	}

	path := fmt.Sprintf("%v/%v", firewallBasePath, firewallID)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	firewall := new(Firewall)
	resp, err := s.client.Do(ctx, req, firewall)
	if err != nil {
		return nil, resp, err
	}

	return firewall, resp, err
}

func (s *FirewallsService) Create(ctx context.Context, createRequest *FirewallCreateRequest) (*Firewall, *Response, error) {
	if createRequest == nil {
		return nil, nil, errors.New("failed to create firewall: payload must be supplied")
	}

	req, err := s.client.NewRequest(http.MethodPost, firewallBasePath, createRequest)
	if err != nil {
		return nil, nil, err
	}

	firewallRoot := new(firewallRoot)
	resp, err := s.client.Do(ctx, req, firewallRoot)
	if err != nil {
		return nil, resp, err
	}

	return firewallRoot.Firewall, resp, nil
}

func (s *FirewallsService) Update(ctx context.Context, firewallID string, updateRequest *FirewallUpdateRequest) (*Firewall, *Response, error) {
	if firewallID == "" {
		return nil, nil, errors.New("failed to update firewall: id must be supplied")
	}
	if updateRequest == nil {
		return nil, nil, errors.New("failed to update firewall: payload must be supplied")
	}

	path := fmt.Sprintf("%v/%v", firewallBasePath, firewallID)
	req, err := s.client.NewRequest(http.MethodPut, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	firewallRoot := new(firewallRoot)
	resp, err := s.client.Do(ctx, req, firewallRoot)
	if err != nil {
		return nil, resp, err
	}

	return firewallRoot.Firewall, resp, nil
}

func (s *FirewallsService) Delete(ctx context.Context, firewallID string) (*Response, error) {
	if firewallID == "" {
		return nil, errors.New("failed to delete firewall: id must be supplied")
	}

	path := fmt.Sprintf("%v/%v", firewallBasePath, firewallID)
	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
