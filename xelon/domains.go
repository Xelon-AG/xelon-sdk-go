package xelon

import (
	"context"
	"errors"
	"fmt"
	"iter"
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

// AllDNSZones returns an iterator to paginate over all DNS zones.
//
// The return iterator can be used in a for...range loop to easily process all zones.
func (s *DomainsService) AllDNSZones(ctx context.Context, opts *ListOptions) (iter.Seq2[DNSZone, *Response], func() error) {
	return newPaginator[DNSZone](ctx, s.client, dnsBasePath, opts)
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

// DeleteDNSZone removes DNS zone by id.
func (s *DomainsService) DeleteDNSZone(ctx context.Context, dnsZoneID string) (*Response, error) {
	if dnsZoneID == "" {
		return nil, errors.New("failed to delete dns zone: id must be supplied")
	}

	path := fmt.Sprintf("%v/%v", dnsBasePath, dnsZoneID)
	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// DNSRecord represetns a Xelon DNS records.
type DNSRecord struct {
	Failover int    `json:"failover,omitempty"`
	Host     string `json:"host,omitempty"`
	ID       int    `json:"id,omitempty"`
	Record   string `json:"record,omitempty"`
	Status   int    `json:"status,omitempty"`
	TTL      int    `json:"ttl,omitempty"`
	Type     string `json:"type,omitempty"`
}

type DNSRecordCreateRequest struct {
	Algorithm       string `json:"algorithm,omitempty"`   // for SSHFP type
	CAAFlag         int    `json:"caaFlag,omitempty"`     // for CAA type
	CAAValue        int    `json:"caaValue,omitempty"`    // for CAA type
	Certificate     string `json:"certificate,omitempty"` // for TLSA type
	Fingerprint     string `json:"fingerprint,omitempty"` // for SSHFP type
	FingerprintType string `json:"fpType,omitempty"`      // for SSHFP type
	Host            string `json:"host"`
	Mail            string `json:"mail,omitempty"`         // for RP type
	MatchingType    int    `json:"matchingType,omitempty"` // for TLSA type
	Port            int    `json:"port,omitempty"`         // for SRV type
	Priority        int    `json:"priority,omitempty"`     // for MX, SRV types
	Record          string `json:"record"`
	Selector        int    `json:"selector,omitempty"` // for TLSA type
	Tag             string `json:"tag,omitempty"`      // for CAA type
	Type            string `json:"type"`
	TTL             int    `json:"ttl"`
	Usage           int    `json:"usage,omitempty"`  // for TLSA type
	Weight          int    `json:"weight,omitempty"` // for SRV type
}

type DNSRecordUpdateRequest struct {
	Algorithm       string `json:"algorithm,omitempty"`   // for SSHFP type
	CAAFlag         int    `json:"caaFlag,omitempty"`     // for CAA type
	CAAValue        int    `json:"caaValue,omitempty"`    // for CAA type
	Certificate     string `json:"certificate,omitempty"` // for TLSA type
	Fingerprint     string `json:"fingerprint,omitempty"` // for SSHFP type
	FingerprintType string `json:"fpType,omitempty"`      // for SSHFP type
	Host            string `json:"host"`
	Mail            string `json:"mail,omitempty"`         // for RP type
	MatchingType    int    `json:"matchingType,omitempty"` // for TLSA type
	Port            int    `json:"port,omitempty"`         // for SRV type
	Priority        int    `json:"priority,omitempty"`     // for MX, SRV types
	Record          string `json:"record"`
	Selector        int    `json:"selector,omitempty"` // for TLSA type
	Tag             string `json:"tag,omitempty"`      // for CAA type
	Type            string `json:"type"`
	TTL             int    `json:"ttl"`
	Usage           int    `json:"usage,omitempty"`  // for TLSA type
	Weight          int    `json:"weight,omitempty"` // for SRV type
}

type dnsRecordRoot struct {
	Message string `json:"message,omitempty"`
}

type dnsRecordsRoot struct {
	DNSRecords []DNSRecord `json:"data"`
}

func (v DNSRecord) String() string { return Stringify(v) }

// ListDNSRecords provides a list of all DNS records for DNS zone.
func (s *DomainsService) ListDNSRecords(ctx context.Context, dnsZoneID string) ([]DNSRecord, *Response, error) {
	if dnsZoneID == "" {
		return nil, nil, errors.New("failed to list dns records: dns zone id must be supplied")
	}

	path := fmt.Sprintf("%v/%v/records", dnsBasePath, dnsZoneID)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(dnsRecordsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.DNSRecords, resp, nil
}

// CreateDNSRecord makes a new DNS record with given payload.
func (s *DomainsService) CreateDNSRecord(ctx context.Context, dnsZoneID string, createRequest *DNSRecordCreateRequest) (*Response, error) {
	if dnsZoneID == "" {
		return nil, errors.New("failed to create dns record: dns zone id must be supplied")
	}
	if createRequest == nil {
		return nil, errors.New("failed to create dns record: payload must be supplied")
	}

	path := fmt.Sprintf("%v/%v/records", dnsBasePath, dnsZoneID)
	req, err := s.client.NewRequest(http.MethodPost, path, createRequest)
	if err != nil {
		return nil, err
	}

	root := new(dnsRecordRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// UpdateDNSRecord changes DNS record identified by id.
func (s *DomainsService) UpdateDNSRecord(ctx context.Context, dnsZoneID, dnsRecordID string, updateRequest *DNSRecordUpdateRequest) (*Response, error) {
	if dnsZoneID == "" {
		return nil, errors.New("failed to update dns record: dns zone id must be supplied")
	}
	if dnsRecordID == "" {
		return nil, errors.New("failed to update dns record: dns record id must be supplied")
	}
	if updateRequest == nil {
		return nil, errors.New("failed to update dns record: payload must be supplied")
	}

	path := fmt.Sprintf("%v/%v/records/%v", dnsBasePath, dnsZoneID, dnsRecordID)
	req, err := s.client.NewRequest(http.MethodPut, path, updateRequest)
	if err != nil {
		return nil, err
	}

	root := new(dnsRecordRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// DeleteDNSRecord removes DNS record identified by id.
func (s *DomainsService) DeleteDNSRecord(ctx context.Context, dnsZoneID, dnsRecordID string) (*Response, error) {
	if dnsZoneID == "" {
		return nil, errors.New("failed to delete dns record: dns zone id must be supplied")
	}
	if dnsRecordID == "" {
		return nil, errors.New("failed to delete dns record: dns record id must be supplied")
	}

	path := fmt.Sprintf("%v/%v/records/%v", dnsBasePath, dnsZoneID, dnsRecordID)
	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
