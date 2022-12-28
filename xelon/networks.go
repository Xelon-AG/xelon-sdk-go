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

type NIC struct {
	ControllerKey int             `json:"niccontrollerkey1,omitempty"`
	Key           int             `json:"nickey1,omitempty"`
	IPs           map[string][]IP `json:"ips,omitempty"`
	Name          string          `json:"nicname,omitempty"`
	Networks      []NICNetwork    `json:"networks,omitempty"`
	Number        int             `json:"nicnumber,omitempty"`
	Unit          int             `json:"nicunit1,omitempty"`
}

type NICNetwork struct {
	Name  string `json:"name,omitempty"`
	Type  string `json:"type,omitempty"`
	Value int    `json:"value,omitempty"`
}

type IP struct {
	ID      int    `json:"value,omitempty"`
	Address string `json:"text,omitempty"`
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

func (s *NetworksService) Get(ctx context.Context) {}
