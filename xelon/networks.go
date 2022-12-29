package xelon

import (
	"context"
	"fmt"
	"net/http"
)

const networkBasePath = "networks"

// NetworksService handles communication with the network related methods of the Xelon API.
type NetworksService service

type Network struct {
	ID      int    `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Network string `json:"network,omitempty"`
	Subnet  string `json:"subnet,omitempty"`
	Type    string `json:"type,omitempty"`
}

type NetworkDetails struct {
	Broadcast      string `json:"broadcast,omitempty"`
	DefaultGateway string `json:"defgw,omitempty"`
	DNSPrimary     string `json:"dns1,omitempty"`
	DNSSecondary   string `json:"dns2,omitempty"`
	ID             int    `json:"id,omitempty"`
	HVSystemID     int    `json:"hv_system_id,omitempty"`
	Name           string `json:"displayname,omitempty"`
	Netmask        string `json:"netmask,omitempty"`
	Network        string `json:"network,omitempty"`
	NetworkID      int    `json:"networkid,omitempty"`
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

// NetworkUpdateRequest represents a request to update a network.
//
// Following fields are mandatory:
//   - NetworkDetails.DefaultGateway
//   - NetworkDetails.DNSPrimary
//   - NetworkDetails.DNSSecondary
//   - NetworkDetails.Name
//   - NetworkDetails.Network
type NetworkUpdateRequest struct {
	*NetworkDetails
}

type NetworkAddIPRequest struct {
	IPAddress string `json:"ipaddress"`
}

type NetworkInfo struct {
	Details *NetworkDetails `json:"localnetid,omitempty"`
	IPs     []IP            `json:"getiplist,omitempty"`
}

type networkRoot struct {
	Networks []Network `json:"networks,omitempty"`
}

// List provides a list of all networks.
func (s *NetworksService) List(ctx context.Context, tenantID string) ([]Network, *Response, error) {
	if tenantID == "" {
		return nil, nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v", tenantID, networkBasePath)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	networkRoot := new(networkRoot)
	resp, err := s.client.Do(ctx, req, networkRoot)
	if err != nil {
		return nil, resp, err
	}

	return networkRoot.Networks, resp, nil
}

// Get provides information about a network identified by local id.
func (s *NetworksService) Get(ctx context.Context, tenantID string, networkID int) (*NetworkInfo, *Response, error) {
	if tenantID == "" {
		return nil, nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v/%v/show", tenantID, networkBasePath, networkID)

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
func (s *NetworksService) CreateLAN(ctx context.Context, tenantID string, createRequest *NetworkLANCreateRequest) (*APIResponse, *Response, error) {
	if tenantID == "" {
		return nil, nil, ErrEmptyArgument
	}
	if createRequest == nil {
		return nil, nil, ErrEmptyPayloadNotAllowed
	}

	path := fmt.Sprintf("%v/%v/addlan", tenantID, networkBasePath)

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
func (s *NetworksService) Update(ctx context.Context, networkID int, updateRequest *NetworkUpdateRequest) (*APIResponse, *Response, error) {
	if updateRequest == nil {
		return nil, nil, ErrEmptyPayloadNotAllowed
	}

	path := fmt.Sprintf("%v/%v/update", networkBasePath, networkID)

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
func (s *NetworksService) Delete(ctx context.Context, networkID int) (*Response, error) {
	path := fmt.Sprintf("%v/%v/destroy", networkBasePath, networkID)

	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// AddIPAddress adds a new IP address to the specific network.
func (s *NetworksService) AddIPAddress(ctx context.Context, networkID int, addIPRequest *NetworkAddIPRequest) (*APIResponse, *Response, error) {
	if addIPRequest == nil {
		return nil, nil, ErrEmptyPayloadNotAllowed
	}

	path := fmt.Sprintf("%v/%v/addIp", networkBasePath, networkID)

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
func (s *NetworksService) DeleteIPAddress(ctx context.Context, networkID, ipAddressID int) (*Response, error) {
	path := fmt.Sprintf("%v/%v/deleteIp?ipid=%v", networkBasePath, networkID, ipAddressID)

	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
