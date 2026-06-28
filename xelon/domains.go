package xelon

import (
	"context"
	"fmt"
	"iter"
	"net/http"
)

const dnsBasePath = "dns"

// DomainsService handles communication with DNS methods of the Xelon API.
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
	Search string `url:"search,omitempty"`
	Sort   string `url:"sort,omitempty"`

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

// ListZones lists DNS zones.
func (s *DomainsService) ListZones(ctx context.Context, opts *DNSZoneListOptions) ([]DNSZone, *Response, error) {
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

// AllZones returns an iterator over all DNS zones.
//
// The returned iterator can be used in a for...range loop.
func (s *DomainsService) AllZones(ctx context.Context, opts *ListOptions) (iter.Seq2[DNSZone, *Response], func() error) {
	return newPaginator[DNSZone](ctx, s.client, dnsBasePath, opts)
}

// GetZone gets a DNS zone by id.
func (s *DomainsService) GetZone(ctx context.Context, dnsZoneID string) (*DNSZone, *Response, error) {
	if dnsZoneID == "" {
		return nil, nil, fmt.Errorf("zone id: %w", ErrEmptyArgument)
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

// CreateZone creates a DNS zone.
func (s *DomainsService) CreateZone(ctx context.Context, createRequest *DNSZoneCreateRequest) (*DNSZone, *Response, error) {
	if createRequest == nil {
		return nil, nil, fmt.Errorf("payload: %w", ErrEmptyPayloadNotAllowed)
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

// DeleteZone deletes a DNS zone by id.
func (s *DomainsService) DeleteZone(ctx context.Context, dnsZoneID string) (*Response, error) {
	if dnsZoneID == "" {
		return nil, fmt.Errorf("zone id: %w", ErrEmptyArgument)
	}

	path := fmt.Sprintf("%v/%v", dnsBasePath, dnsZoneID)
	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// DNSRecordType represents a supported DNS record type.
type DNSRecordType string

const (
	DNSRecordTypeA     DNSRecordType = "A"
	DNSRecordTypeAAAA  DNSRecordType = "AAAA"
	DNSRecordTypeALIAS DNSRecordType = "ALIAS"
	DNSRecordTypeCAA   DNSRecordType = "CAA"
	DNSRecordTypeCNAME DNSRecordType = "CNAME"
	DNSRecordTypeMX    DNSRecordType = "MX"
	DNSRecordTypeNS    DNSRecordType = "NS"
	DNSRecordTypePTR   DNSRecordType = "PTR"
	DNSRecordTypeRP    DNSRecordType = "RP"
	DNSRecordTypeSRV   DNSRecordType = "SRV"
	DNSRecordTypeSSHFP DNSRecordType = "SSHFP"
	DNSRecordTypeTLSA  DNSRecordType = "TLSA"
	DNSRecordTypeTXT   DNSRecordType = "TXT"
)

// DNSRecord represents a Xelon DNS record.
type DNSRecord struct {
	Failover int           `json:"failover,omitempty"`
	Host     string        `json:"host,omitempty"`
	ID       int           `json:"id,omitempty"`
	Record   string        `json:"record,omitempty"`
	Status   int           `json:"status,omitempty"`
	TTL      int           `json:"ttl,omitempty"`
	Type     DNSRecordType `json:"type,omitempty"`
}

type DNSRecordCreateRequest struct {
	Algorithm       int           `json:"algorithm,omitempty"`    // for SSHFP type
	CAAFlag         int           `json:"caaFlag,omitempty"`      // for CAA type
	CAAType         string        `json:"caaType,omitempty"`      // for CAA type
	CAAValue        string        `json:"caaValue,omitempty"`     // for CAA type
	Certificate     string        `json:"certificate,omitempty"`  // for TLSA type
	Fingerprint     string        `json:"fingerprint,omitempty"`  // for SSHFP type
	FingerprintType string        `json:"fpType,omitempty"`       // for SSHFP type
	Host            string        `json:"host"`                   // required
	Mail            string        `json:"mail,omitempty"`         // for RP type
	MatchingType    int           `json:"matchingType,omitempty"` // for TLSA type
	Port            int           `json:"port,omitempty"`         // for SRV type
	Priority        int           `json:"priority,omitempty"`     // for MX, SRV types
	Record          string        `json:"record"`
	Selector        int           `json:"selector,omitempty"` // for TLSA type
	Tag             string        `json:"tag,omitempty"`      // for CAA type
	TTL             int           `json:"ttl"`                // required
	Type            DNSRecordType `json:"type"`               // required
	Usage           int           `json:"usage,omitempty"`    // for TLSA type
	Weight          int           `json:"weight,omitempty"`   // for SRV type
}

type DNSRecordUpdateRequest struct {
	Algorithm       int           `json:"algorithm,omitempty"`    // for SSHFP type
	CAAFlag         int           `json:"caaFlag,omitempty"`      // for CAA type
	CAAType         string        `json:"caaType,omitempty"`      // for CAA type
	CAAValue        string        `json:"caaValue,omitempty"`     // for CAA type
	Certificate     string        `json:"certificate,omitempty"`  // for TLSA type
	Fingerprint     string        `json:"fingerprint,omitempty"`  // for SSHFP type
	FingerprintType string        `json:"fpType,omitempty"`       // for SSHFP type
	Host            string        `json:"host"`                   // required
	Mail            string        `json:"mail,omitempty"`         // for RP type
	MatchingType    int           `json:"matchingType,omitempty"` // for TLSA type
	Port            int           `json:"port,omitempty"`         // for SRV type
	Priority        int           `json:"priority,omitempty"`     // for MX, SRV types
	Record          string        `json:"record"`
	Selector        int           `json:"selector,omitempty"` // for TLSA type
	Tag             string        `json:"tag,omitempty"`      // for CAA type
	TTL             int           `json:"ttl"`                // required
	Type            DNSRecordType `json:"type"`               // required
	Usage           int           `json:"usage,omitempty"`    // for TLSA type
	Weight          int           `json:"weight,omitempty"`   // for SRV type
}

type dnsRecordRoot struct {
	Message string `json:"message,omitempty"`
}

type dnsRecordsRoot struct {
	DNSRecords []DNSRecord `json:"data"`
}

func (v DNSRecord) String() string { return Stringify(v) }

// ListRecords lists records for a DNS zone.
func (s *DomainsService) ListRecords(ctx context.Context, dnsZoneID string) ([]DNSRecord, *Response, error) {
	if dnsZoneID == "" {
		return nil, nil, fmt.Errorf("zone id: %w", ErrEmptyArgument)
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

// CreateRecord creates a DNS record.
func (s *DomainsService) CreateRecord(ctx context.Context, dnsZoneID string, createRequest *DNSRecordCreateRequest) (*Response, error) {
	if dnsZoneID == "" {
		return nil, fmt.Errorf("zone id: %w", ErrEmptyArgument)
	}
	if createRequest == nil {
		return nil, fmt.Errorf("payload: %w", ErrEmptyPayloadNotAllowed)
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

// UpdateRecord updates a DNS record by id.
func (s *DomainsService) UpdateRecord(ctx context.Context, dnsZoneID string, dnsRecordID int, updateRequest *DNSRecordUpdateRequest) (*Response, error) {
	if dnsZoneID == "" {
		return nil, fmt.Errorf("zone id: %w", ErrEmptyArgument)
	}
	if dnsRecordID <= 0 {
		return nil, fmt.Errorf("record id: %w", ErrEmptyArgument)
	}
	if updateRequest == nil {
		return nil, fmt.Errorf("payload: %w", ErrEmptyPayloadNotAllowed)
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

// DeleteRecord deletes a DNS record by id.
func (s *DomainsService) DeleteRecord(ctx context.Context, dnsZoneID string, dnsRecordID int) (*Response, error) {
	if dnsZoneID == "" {
		return nil, fmt.Errorf("zone id: %w", ErrEmptyArgument)
	}
	if dnsRecordID <= 0 {
		return nil, fmt.Errorf("record id: %w", ErrEmptyArgument)
	}

	path := fmt.Sprintf("%v/%v/records/%v", dnsBasePath, dnsZoneID, dnsRecordID)
	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// DNSSOA represents a Xelon DNS SOA record.
type DNSSOA struct {
	AdminEmail   string `json:"adminMail,omitempty"`
	Expire       int    `json:"expire,omitempty"`
	PrimaryNS    string `json:"primaryNs,omitempty"`
	Refresh      int    `json:"refresh,omitempty"`
	Retry        int    `json:"retry,omitempty"`
	SerialNumber int    `json:"serialNumber,omitempty"`
	TTL          int    `json:"ttl,omitempty"`
}

type DNSSOAUpdateRequest struct {
	AdminEmail string `json:"adminMail"`
	Expire     int    `json:"expire"`
	PrimaryNS  string `json:"primaryNs"`
	Refresh    int    `json:"refresh"`
	Retry      int    `json:"retry"`
	TTL        int    `json:"ttl"`
}

type dnsSOARoot struct {
	Message string `json:"message,omitempty"`
}

func (v DNSSOA) String() string { return Stringify(v) }

// GetSOA gets the SOA record for a DNS zone.
func (s *DomainsService) GetSOA(ctx context.Context, dnsZoneID string) (*DNSSOA, *Response, error) {
	if dnsZoneID == "" {
		return nil, nil, fmt.Errorf("zone id: %w", ErrEmptyArgument)
	}

	path := fmt.Sprintf("%v/%v/soa", dnsBasePath, dnsZoneID)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	dnsSOA := new(DNSSOA)
	resp, err := s.client.Do(ctx, req, dnsSOA)
	if err != nil {
		return nil, resp, err
	}

	return dnsSOA, resp, nil
}

// UpdateSOA updates the SOA record for a DNS zone.
func (s *DomainsService) UpdateSOA(ctx context.Context, dnsZoneID string, updateRequest *DNSSOAUpdateRequest) (*Response, error) {
	if dnsZoneID == "" {
		return nil, fmt.Errorf("zone id: %w", ErrEmptyArgument)
	}
	if updateRequest == nil {
		return nil, fmt.Errorf("payload: %w", ErrEmptyPayloadNotAllowed)
	}

	path := fmt.Sprintf("%v/%v/soa", dnsBasePath, dnsZoneID)
	req, err := s.client.NewRequest(http.MethodPut, path, updateRequest)
	if err != nil {
		return nil, err
	}

	root := new(dnsSOARoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
