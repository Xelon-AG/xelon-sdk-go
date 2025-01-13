package xelon

import (
	"context"
	"fmt"
	"net/http"
)

const networkBasePathV1 = "networks"

// NetworksServiceV1 handles communication with the network related methods of the Xelon API.
// Deprecated.
type NetworksServiceV1 service

type NetworkV1 struct {
	ID      int    `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Network string `json:"network,omitempty"`
	Subnet  string `json:"subnet,omitempty"`
	Type    string `json:"type,omitempty"`
}

type NetworkV1Details struct {
	Broadcast      string `json:"broadcast,omitempty"`
	DefaultGateway string `json:"defgw,omitempty"`
	DNSPrimary     string `json:"dns1,omitempty"`
	DNSSecondary   string `json:"dns2,omitempty"`
	ID             int    `json:"id,omitempty"`
	Name           string `json:"displayname,omitempty"`
	Netmask        string `json:"netmask,omitempty"`
	Network        string `json:"network,omitempty"`
	Subnet         string `json:"subnet,omitempty"`
	Type           string `json:"type,omitempty"`
}

type IP struct {
	ID        int    `json:"id,omitempty"`
	IP        string `json:"ip,omitempty"`
	IPType    int    `json:"iptype,omitempty"`
	LocalVMID string `json:"localvmid,omitempty"`
	NetworkID int    `json:"networkid,omitempty"`
	NICUnit   int    `json:"nicunit,omitempty"`
	Type      string `json:"type,omitempty"`
}

type NIC struct {
	ControllerKey int                `json:"niccontrollerkey1,omitempty"`
	Key           int                `json:"nickey1,omitempty"`
	IPs           map[string][]NICIP `json:"ips,omitempty"`
	Name          string             `json:"nicname,omitempty"`
	Networks      []NICNetwork       `json:"networks,omitempty"`
	Number        int                `json:"nicnumber,omitempty"`
	Unit          int                `json:"nicunit1,omitempty"`
}

type NICNetwork struct {
	Name  string `json:"name,omitempty"`
	Type  string `json:"type,omitempty"`
	Value int    `json:"value,omitempty"`
}

type NICIP struct {
	ID      int    `json:"value,omitempty"`
	Address string `json:"text,omitempty"`
}

// NetworkLANCreateRequest represents a request to create a LAN network.
//
// IPAddresses is optional, can either be single "10.0.0.1" or multiple
// "10.0.0.1-10.0.0.254" entries.
type NetworkLANCreateRequest struct {
	CloudID      int    `json:"cloudId"`
	DisplayName  string `json:"displayname"`
	DNSPrimary   string `json:"dns1"`
	DNSSecondary string `json:"dns2"`
	Gateway      string `json:"gateway"`
	Network      string `json:"network"`
	IPAddresses  string `json:"ip_addresses,omitempty"`
}

// NetworkV1UpdateRequest represents a request to update a network.
//
// Following fields are mandatory:
//   - NetworkDetails.DefaultGateway
//   - NetworkDetails.DNSPrimary
//   - NetworkDetails.DNSSecondary
//   - NetworkDetails.Name
//   - NetworkDetails.Network
type NetworkV1UpdateRequest struct {
	*NetworkV1Details
}

type NetworkAddIPRequest struct {
	IPAddress string `json:"ipaddress"`
}

type NetworkInfo struct {
	CloudID int               `json:"hvSystemId,omitempty"`
	Details *NetworkV1Details `json:"localnetid,omitempty"`
	IPs     []IP              `json:"getiplist,omitempty"`
}

type networkRootV1 struct {
	Networks []NetworkV1 `json:"networks,omitempty"`
}

// List provides a list of all networks.
func (s *NetworksServiceV1) List(ctx context.Context, tenantID string) ([]NetworkV1, *Response, error) {
	if tenantID == "" {
		return nil, nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v", tenantID, networkBasePathV1)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	networkRoot := new(networkRootV1)
	resp, err := s.client.Do(ctx, req, networkRoot)
	if err != nil {
		return nil, resp, err
	}

	return networkRoot.Networks, resp, nil
}

// Get provides information about a network identified by local id.
func (s *NetworksServiceV1) Get(ctx context.Context, tenantID string, networkID int) (*NetworkInfo, *Response, error) {
	if tenantID == "" {
		return nil, nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v/%v/show", tenantID, networkBasePathV1, networkID)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	networkInfo := new(NetworkInfo)
	resp, err := s.client.Do(ctx, req, networkInfo)
	if err != nil {
		return nil, resp, err
	}

	return networkInfo, resp, nil
}

// CreateLAN makes a new LAN network with given payload.
func (s *NetworksServiceV1) CreateLAN(ctx context.Context, tenantID string, createRequest *NetworkLANCreateRequest) (*APIResponse, *Response, error) {
	if tenantID == "" {
		return nil, nil, ErrEmptyArgument
	}
	if createRequest == nil {
		return nil, nil, ErrEmptyPayloadNotAllowed
	}

	path := fmt.Sprintf("%v/%v/addlan", tenantID, networkBasePathV1)

	req, err := s.client.NewRequest(http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	apiResponse := new(APIResponse)
	resp, err := s.client.Do(ctx, req, apiResponse)
	if err != nil {
		return nil, resp, err
	}

	return apiResponse, resp, nil
}

// Update changes the configuration of a network.
func (s *NetworksServiceV1) Update(ctx context.Context, networkID int, updateRequest *NetworkV1UpdateRequest) (*APIResponse, *Response, error) {
	if updateRequest == nil {
		return nil, nil, ErrEmptyPayloadNotAllowed
	}

	path := fmt.Sprintf("%v/%v/update", networkBasePathV1, networkID)

	req, err := s.client.NewRequest(http.MethodPut, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	apiResponse := new(APIResponse)
	resp, err := s.client.Do(ctx, req, apiResponse)
	if err != nil {
		return nil, resp, err
	}

	return apiResponse, resp, nil
}

// Delete removes a network identified by id.
func (s *NetworksServiceV1) Delete(ctx context.Context, networkID int) (*Response, error) {
	path := fmt.Sprintf("%v/%v/destroy", networkBasePathV1, networkID)

	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// AddIPAddress adds a new IP address to the specific network.
func (s *NetworksServiceV1) AddIPAddress(ctx context.Context, networkID int, addIPRequest *NetworkAddIPRequest) (*APIResponse, *Response, error) {
	if addIPRequest == nil {
		return nil, nil, ErrEmptyPayloadNotAllowed
	}

	path := fmt.Sprintf("%v/%v/addIp", networkBasePathV1, networkID)

	req, err := s.client.NewRequest(http.MethodPost, path, addIPRequest)
	if err != nil {
		return nil, nil, err
	}

	apiResponse := new(APIResponse)
	resp, err := s.client.Do(ctx, req, apiResponse)
	if err != nil {
		return nil, resp, err
	}

	return apiResponse, nil, nil
}

// DeleteIPAddress removes an IP from the specific network.
func (s *NetworksServiceV1) DeleteIPAddress(ctx context.Context, networkID, ipAddressID int) (*Response, error) {
	path := fmt.Sprintf("%v/%v/deleteIp?ipid=%v", networkBasePathV1, networkID, ipAddressID)

	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
