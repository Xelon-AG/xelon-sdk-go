package xelon

import (
	"context"
	"encoding/json"
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
	Cloud             *Cloud                   `json:"cloud,omitempty"`
	CreatedAt         *time.Time               `json:"createdAt,omitempty"`
	ExternalIPAddress string                   `json:"externalIp,omitempty"`
	ForwardingRules   []FirewallForwardingRule `json:"forwardingRules,omitempty"`
	HealthStatus      string                   `json:"health,omitempty"`
	ID                string                   `json:"identifier,omitempty"`
	InternalIPAddress string                   `json:"internalIp,omitempty"`
	Name              string                   `json:"name,omitempty"`
	State             int                      `json:"state,omitempty"`
	Tenant            *Tenant                  `json:"tenant,omitempty"`
}

type FirewallForwardingRule struct {
	DestinationIPAddressWrapper any      `json:"destinationIp,omitempty"`
	DestinationIPAddress        string   `json:"-"`
	DestinationIPAddresses      []string `json:"-"`
	ExternalPort                int      `json:"externalPort,omitempty"`
	ID                          string   `json:"identifier,omitempty"`
	InternalPort                int      `json:"port,omitempty"`
	Protocol                    string   `json:"protocol,omitempty"`
	SourceIPAddressWrapper      any      `json:"sourceIp,omitempty"`
	SourceIPAddress             string   `json:"-"`
	SourceIPAddresses           []string `json:"-"`
	Type                        string   `json:"type,omitempty"`
}

// MarshalJSON marshals sourceIp and destinationIp in FirewallForwardingRule depending on type:
//   - inbound: `sourceIp` is slice of strings, `destinationIp` is string
//   - outbound: `sourceIp` is string, `destinationIp` is slice of strings
//
// This workaround is needed because of backend API.
func (v *FirewallForwardingRule) MarshalJSON() ([]byte, error) {
	if v.DestinationIPAddress != "" {
		v.DestinationIPAddressWrapper = v.DestinationIPAddress
	}
	if len(v.DestinationIPAddresses) > 0 {
		var ipAddresses []string
		ipAddresses = append(ipAddresses, v.DestinationIPAddresses...)
		v.DestinationIPAddressWrapper = ipAddresses
	}

	if v.SourceIPAddress != "" {
		v.SourceIPAddressWrapper = v.SourceIPAddress
	}
	if len(v.SourceIPAddresses) > 0 {
		var ipAddresses []string
		ipAddresses = append(ipAddresses, v.SourceIPAddresses...)
		v.SourceIPAddressWrapper = ipAddresses
	}

	type alias FirewallForwardingRule
	return json.Marshal(&struct {
		*alias
		DestinationIPAddressWrapper any `json:"destinationIp,omitempty"`
		SourceIPAddressWrapper      any `json:"sourceIp,omitempty"`
	}{
		alias:                       (*alias)(v),
		DestinationIPAddressWrapper: v.DestinationIPAddressWrapper,
		SourceIPAddressWrapper:      v.SourceIPAddressWrapper,
	})
}

// UnmarshalJSON parses sourceIp and destinationIp in FirewallForwardingRule depending on type:
//   - inbound: `sourceIp` is slice of strings, `destinationIp` is string
//   - outbound: `sourceIp` is string, `destinationIp` is slice of strings
//
// This workaround is needed because of backend API.
func (v *FirewallForwardingRule) UnmarshalJSON(data []byte) error {
	type alias FirewallForwardingRule
	rule := &struct{ *alias }{alias: (*alias)(v)}
	if err := json.Unmarshal(data, &rule); err != nil {
		return err
	}

	switch val := rule.DestinationIPAddressWrapper.(type) {
	case string:
		v.DestinationIPAddress = val
		v.DestinationIPAddresses = nil
	case []any:
		var ipAddresses []string
		for _, ipAddress := range val {
			ipAddresses = append(ipAddresses, ipAddress.(string))
		}
		v.DestinationIPAddress = ""
		v.DestinationIPAddresses = ipAddresses
	}

	switch val := rule.SourceIPAddressWrapper.(type) {
	case string:
		v.SourceIPAddress = val
		v.SourceIPAddresses = nil
	case []any:
		var ipAddresses []string
		for _, ipAddress := range val {
			ipAddresses = append(ipAddresses, ipAddress.(string))
		}
		v.SourceIPAddress = ""
		v.SourceIPAddresses = ipAddresses
	}

	return nil
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

type FirewallCreateForwardingRuleRequest struct {
	FirewallForwardingRule
}

type FirewallUpdateForwardingRuleRequest struct {
	FirewallForwardingRule
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

type firewallForwardingRuleRoot struct {
	ForwardingRule *FirewallForwardingRule `json:"data,omitempty"`
	Message        string                  `json:"message,omitempty"`
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

// Create makes a new firewall with given payload.
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

// Update changes firewall identified by id.
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

// Delete removes firewall identified by id.
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

// CreateForwardingRule makes a new forwarding rule.
func (s *FirewallsService) CreateForwardingRule(ctx context.Context, firewallID string, createRequest *FirewallCreateForwardingRuleRequest) (*FirewallForwardingRule, *Response, error) {
	if firewallID == "" {
		return nil, nil, errors.New("failed to create forwarding rule: firewall id must be supplied")
	}
	if createRequest == nil {
		return nil, nil, errors.New("failed to create forwarding rule: payload must be supplied")
	}

	path := fmt.Sprintf("%v/%v/rules", firewallBasePath, firewallID)
	req, err := s.client.NewRequest(http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	firewallForwardingRuleRoot := new(firewallForwardingRuleRoot)
	resp, err := s.client.Do(ctx, req, firewallForwardingRuleRoot)
	if err != nil {
		return nil, resp, err
	}

	return firewallForwardingRuleRoot.ForwardingRule, resp, nil
}

// UpdateForwardingRule changes the configuration of a forwarding rule.
func (s *FirewallsService) UpdateForwardingRule(ctx context.Context, firewallID, forwardingRuleID string, updateRequest *FirewallUpdateForwardingRuleRequest) (*FirewallForwardingRule, *Response, error) {
	if firewallID == "" {
		return nil, nil, errors.New("failed to update forwarding rule: firewall id must be supplied")
	}
	if forwardingRuleID == "" {
		return nil, nil, errors.New("failed to update forwarding rule: forwarding rule id must be supplied")
	}
	if updateRequest == nil {
		return nil, nil, errors.New("failed to update forwarding rule: payload must be supplied")
	}

	path := fmt.Sprintf("%v/%v/rules/%v", firewallBasePath, firewallID, forwardingRuleID)
	req, err := s.client.NewRequest(http.MethodPut, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	firewallForwardingRuleRoot := new(firewallForwardingRuleRoot)
	resp, err := s.client.Do(ctx, req, firewallForwardingRuleRoot)
	if err != nil {
		return nil, resp, err
	}

	return firewallForwardingRuleRoot.ForwardingRule, resp, nil
}

// DeleteForwardingRule removes a forwarding rule.
func (s *FirewallsService) DeleteForwardingRule(ctx context.Context, firewallID string, forwardingRuleID int) (*Response, error) {
	if firewallID == "" {
		return nil, errors.New("failed to delete forwarding rule: firewall id must be supplied")
	}
	if forwardingRuleID == 0 {
		return nil, errors.New("failed to delete forwarding rule: forwarding rule id must be supplied")
	}

	path := fmt.Sprintf("%v/%v/rules/%v", firewallBasePath, firewallID, forwardingRuleID)
	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
