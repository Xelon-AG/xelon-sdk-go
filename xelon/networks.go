package xelon

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

const networkBasePath = "networks"

// NetworksService handles communication with the network related methods of the Xelon REST API.
type NetworksService service

type Network struct {
	Clouds       []Cloud `json:"clouds,omitempty"`
	DNSPrimary   string  `json:"dns1,omitempty"`
	DNSSecondary string  `json:"dns2,omitempty"`
	Free         bool    `json:"isFree,omitempty"`
	Gateway      string  `json:"gateway,omitempty"`
	ID           string  `json:"identifier,omitempty"`
	Name         string  `json:"name,omitempty"`
	Network      string  `json:"network,omitempty"`
	NetworkSpeed int     `json:"networkSpeedValue,omitempty"`
	Stretched    bool    `json:"isStretched,omitempty"`
	SubnetSize   int     `json:"networkSize,omitempty"`
	Type         string  `json:"type,omitempty"`
}

type NetworkLANCreateRequest struct {
	CloudID            string `json:"cloudIdentifier"`
	CloudForStretching string `json:"cloudForStretching,omitempty"`
	DNSPrimary         string `json:"dns1"`
	DNSSecondary       string `json:"dns2,omitempty"`
	Gateway            string `json:"gateway"`
	Name               string `json:"name"`
	Network            string `json:"network"`
	NetworkSpeed       int    `json:"networkSpeedValue"`
	Stretched          bool   `json:"isStretched,omitempty"`
	SubnetSize         int    `json:"networkSize"`
	TenantID           string `json:"tenantIdentifier,omitempty"`
}

type NetworkLANUpdateRequest struct {
	DNSPrimary   string `json:"dns1"`
	DNSSecondary string `json:"dns2,omitempty"`
	Gateway      string `json:"gateway"`
	Name         string `json:"name"`
	Network      string `json:"network"`
	NetworkSpeed int    `json:"networkSpeedValue"`
}

type NetworkWANCreateRequest struct {
	CloudID            string `json:"cloudIdentifier"`
	CloudForStretching string `json:"cloudForStretching,omitempty"`
	Name               string `json:"name"`
	NetworkSpeed       int    `json:"networkSpeedValue"`
	Stretched          bool   `json:"isStretched,omitempty"`
	SubnetSize         int    `json:"networkSize"`
	TenantID           string `json:"tenantIdentifier,omitempty"`
}

// NetworkListOptions specifies the optional parameters to the NetworksService.List.
type NetworkListOptions struct {
	Sort   string `url:"sort,omitempty"`
	Search string `url:"search,omitempty"`

	ListOptions
}

type networkRoot struct {
	Network *Network `json:"data,omitempty"`
	Message string   `json:"message,omitempty"`
}

type networksRoot struct {
	Networks []Network `json:"data"`
	Meta     *Meta     `json:"meta,omitempty"`
}

func (v Network) String() string { return Stringify(v) }

// List provides a list of all networks.
func (s *NetworksService) List(ctx context.Context, opts *NetworkListOptions) ([]Network, *Response, error) {
	path, err := addOptions(networkBasePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(networksRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}

	return root.Networks, resp, nil
}

// Get provides detailed information for network identified by id.
func (s *NetworksService) Get(ctx context.Context, networkID string) (*Network, *Response, error) {
	if networkID == "" {
		return nil, nil, errors.New("failed to get network: id must be supplied")
	}

	path := fmt.Sprintf("%v/%v", networkBasePath, networkID)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	network := new(Network)
	resp, err := s.client.Do(ctx, req, network)
	if err != nil {
		return nil, resp, err
	}

	return network, resp, nil
}

// CreateLAN makes a new LAN network with given payload.
func (s *NetworksService) CreateLAN(ctx context.Context, createRequest *NetworkLANCreateRequest) (*Network, *Response, error) {
	if createRequest == nil {
		return nil, nil, errors.New("failed to create LAN network: payload must be supplied")
	}

	path := fmt.Sprintf("%v/lan", networkBasePath)
	req, err := s.client.NewRequest(http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	networkRoot := new(networkRoot)
	resp, err := s.client.Do(ctx, req, networkRoot)
	if err != nil {
		return nil, resp, err
	}

	return networkRoot.Network, resp, nil
}

// UpdateLAN changes network identified by id.
func (s *NetworksService) UpdateLAN(ctx context.Context, networkID string, updateRequest *NetworkLANUpdateRequest) (*Network, *Response, error) {
	if networkID == "" {
		return nil, nil, errors.New("failed to update LAN network: id must be supplied")
	}
	if updateRequest == nil {
		return nil, nil, errors.New("failed to update LAN network: payload must be supplied")
	}

	path := fmt.Sprintf("%v/%v/lan", networkBasePath, networkID)
	req, err := s.client.NewRequest(http.MethodPatch, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	network := new(Network)
	resp, err := s.client.Do(ctx, req, network)
	if err != nil {
		return nil, resp, err
	}

	return network, resp, nil
}

// CreateWAN makes a new WAN network with given payload.
func (s *NetworksService) CreateWAN(ctx context.Context, createRequest *NetworkWANCreateRequest) (*Network, *Response, error) {
	if createRequest == nil {
		return nil, nil, errors.New("failed to create WAN network: payload must be supplied")
	}

	path := fmt.Sprintf("%v/wan", networkBasePath)
	req, err := s.client.NewRequest(http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	networkRoot := new(networkRoot)
	resp, err := s.client.Do(ctx, req, networkRoot)
	if err != nil {
		return nil, resp, err
	}

	return networkRoot.Network, resp, nil
}

// Delete removes network identified by id.
func (s *NetworksService) Delete(ctx context.Context, networkID string) (*Response, error) {
	if networkID == "" {
		return nil, errors.New("failed to delete network: id must be supplied")
	}

	path := fmt.Sprintf("%v/%v", networkBasePath, networkID)
	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
