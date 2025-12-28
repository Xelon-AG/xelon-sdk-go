package xelon

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

const dnsBasePath = "dns"

// DomainsService handles communication with the DNS related methods of the Xelon REST API.
type DomainsService service

// DNSZone represents a Xelon DNS zone.
type DNSZone struct {
	ID        string `json:"identifier,omitempty"`
	Name      string `json:"name,omitempty"`
	OwnerName string `json:"ownerName,omitempty"`
}

type DNSZoneCreateRequest struct {
	Domain string `json:"domain"`
}

type DNSZoneListOptions struct {
	Sort   string `url:"sort,omitempty"`
	Search string `url:"search,omitempty"`

	ListOptions
}

type dnsZoneRoot struct {
	DNSZone *DNSZone `json:"data,omitempty"`
	Meta    *Meta    `json:"meta,omitempty"`
}

type dnsZonesRoot struct {
	DNSZones []DNSZone `json:"data"`
	Meta     *Meta     `json:"meta,omitempty"`
}

func (v DNSZone) String() string { return Stringify(v) }

// ListDNSZones provides a list of all DNS zones.
func (s *DomainsService) ListDNSZones(ctx context.Context, opts *DNSZoneListOptions) ([]DNSZone, *Response, error) {
	path, err := addOptions(dnsBasePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(dnsZonesRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}

	return root.DNSZones, resp, nil
}

// GetDNSZone provides detailed information for DNS zone identified by id.
func (s *DomainsService) GetDNSZone(ctx context.Context, dnsZoneID string) (*DNSZone, *Response, error) {
	if dnsZoneID == "" {
		return nil, nil, errors.New("failed to get dns zone: id must be supplied")
	}

	path := fmt.Sprintf("%v/%v", dnsBasePath, dnsZoneID)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	dnsZone := new(DNSZone)
	resp, err := s.client.Do(ctx, req, dnsZone)
	if err != nil {
		return nil, resp, err
	}

	return dnsZone, resp, err
}

// CreateDNSZone makes a new DNS with given payload.
func (s *DomainsService) CreateDNSZone(ctx context.Context, createRequest *DNSZoneCreateRequest) (*DNSZone, *Response, error) {
	if createRequest == nil {
		return nil, nil, errors.New("failed to create dns zone: payload must be supplied")
	}

	req, err := s.client.NewRequest(http.MethodPost, dnsBasePath, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(dnsZoneRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.DNSZone, resp, nil
}
