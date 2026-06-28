package xelon

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDomains_ListZones(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("GET /dns", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fixture := loadFixture(t, "domains_list_zones.json")
		_, _ = w.Write(fixture)
	})
	expectedZones := []DNSZone{{
		ID:        "dns-zone-1",
		Name:      "example.com",
		OwnerName: "test-tenant",
	}, {
		ID:        "dns-zone-2",
		Name:      "example.net",
		OwnerName: "test-tenant",
	}}

	actualZones, resp, err := client.Domains.ListZones(ctx, nil)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedZones, actualZones)
	assert.Equal(t, &Meta{
		Total:    2,
		LastPage: 1,
		PerPage:  10,
		Page:     1,
		From:     1,
		To:       2,
	}, resp.Meta)
}

func TestDomains_GetZone(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("GET /dns/dns-zone-1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fixture := loadFixture(t, "domains_get_zone_success.json")
		_, _ = w.Write(fixture)
	})
	expectedZone := &DNSZone{
		ID:        "dns-zone-1",
		Name:      "example.com",
		OwnerName: "test-tenant",
	}

	actualZone, resp, err := client.Domains.GetZone(ctx, "dns-zone-1")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedZone, actualZone)
}

func TestDomains_CreateZone(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("POST /dns", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		var actualRequest DNSZoneCreateRequest
		err := json.NewDecoder(r.Body).Decode(&actualRequest)
		assert.NoError(t, err)
		assert.Equal(t, DNSZoneCreateRequest{Domain: "example.com"}, actualRequest)

		fixture := loadFixture(t, "domains_create_zone_success.json")
		_, _ = w.Write(fixture)
	})
	expectedZone := &DNSZone{
		ID:        "dns-zone-1",
		Name:      "example.com",
		OwnerName: "test-tenant",
	}

	actualZone, resp, err := client.Domains.CreateZone(ctx, &DNSZoneCreateRequest{Domain: "example.com"})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedZone, actualZone)
}

func TestDomains_DeleteZone(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("DELETE /dns/dns-zone-1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		fixture := loadFixture(t, "domains_delete_zone_success.json")
		_, _ = w.Write(fixture)
	})

	resp, err := client.Domains.DeleteZone(ctx, "dns-zone-1")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestDomains_GetSOA(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("GET /dns/dns-zone-1/soa", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fixture := loadFixture(t, "domains_get_soa_success.json")
		_, _ = w.Write(fixture)
	})
	expectedSOA := &DNSSOA{
		SerialNumber: 2025011302,
		PrimaryNS:    "ns1.xelon.ch",
		AdminEmail:   "hostmaster.xelon.ch",
		Refresh:      10800,
		Retry:        3600,
		Expire:       604800,
		TTL:          3600,
	}

	actualSOA, resp, err := client.Domains.GetSOA(ctx, "dns-zone-1")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedSOA, actualSOA)
}

func TestDomains_UpdateSOA(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("PUT /dns/dns-zone-1/soa", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)

		var actualRequest DNSSOAUpdateRequest
		err := json.NewDecoder(r.Body).Decode(&actualRequest)
		assert.NoError(t, err)
		assert.Equal(t, DNSSOAUpdateRequest{
			PrimaryNS:  "ns1.xelon.ch",
			AdminEmail: "hostmaster.xelon.ch",
			Refresh:    10800,
			Retry:      3600,
			Expire:     604800,
			TTL:        3600,
		}, actualRequest)

		fixture := loadFixture(t, "domains_update_soa_success.json")
		_, _ = w.Write(fixture)
	})

	resp, err := client.Domains.UpdateSOA(ctx, "dns-zone-1", &DNSSOAUpdateRequest{
		PrimaryNS:  "ns1.xelon.ch",
		AdminEmail: "hostmaster.xelon.ch",
		Refresh:    10800,
		Retry:      3600,
		Expire:     604800,
		TTL:        3600,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestDomains_ListRecords(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("GET /dns/dns-zone-1/records", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fixture := loadFixture(t, "domains_list_records.json")
		_, _ = w.Write(fixture)
	})
	expectedRecords := []DNSRecord{{
		ID:       244042432,
		Type:     DNSRecordTypeA,
		Host:     "www",
		Record:   "192.0.2.10",
		Failover: 0,
		TTL:      3600,
		Status:   1,
	}, {
		ID:       244042433,
		Type:     DNSRecordTypeCNAME,
		Host:     "app",
		Record:   "www.example.com",
		Failover: 0,
		TTL:      3600,
		Status:   1,
	}}

	actualRecords, resp, err := client.Domains.ListRecords(ctx, "dns-zone-1")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Nil(t, resp.Meta)
	assert.Equal(t, expectedRecords, actualRecords)
}

func TestDomains_CreateRecord(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("POST /dns/dns-zone-1/records", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		var actualRequest DNSRecordCreateRequest
		err := json.NewDecoder(r.Body).Decode(&actualRequest)
		assert.NoError(t, err)
		assert.Equal(t, DNSRecordCreateRequest{
			Type:   DNSRecordTypeA,
			Host:   "www",
			Record: "192.0.2.10",
			TTL:    3600,
		}, actualRequest)

		fixture := loadFixture(t, "domains_create_record_success.json")
		_, _ = w.Write(fixture)
	})

	resp, err := client.Domains.CreateRecord(ctx, "dns-zone-1", &DNSRecordCreateRequest{
		Type:   DNSRecordTypeA,
		Host:   "www",
		Record: "192.0.2.10",
		TTL:    3600,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestDomains_UpdateRecord(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("PUT /dns/dns-zone-1/records/244042432", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)

		var actualRequest DNSRecordUpdateRequest
		err := json.NewDecoder(r.Body).Decode(&actualRequest)
		assert.NoError(t, err)
		assert.Equal(t, DNSRecordUpdateRequest{
			Type:   DNSRecordTypeCNAME,
			Host:   "app",
			Record: "www.example.com",
			TTL:    3600,
		}, actualRequest)

		fixture := loadFixture(t, "domains_update_record_success.json")
		_, _ = w.Write(fixture)
	})

	resp, err := client.Domains.UpdateRecord(ctx, "dns-zone-1", 244042432, &DNSRecordUpdateRequest{
		Type:   DNSRecordTypeCNAME,
		Host:   "app",
		Record: "www.example.com",
		TTL:    3600,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestDomains_DeleteRecord(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("DELETE /dns/dns-zone-1/records/244042432", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		fixture := loadFixture(t, "domains_delete_record_success.json")
		_, _ = w.Write(fixture)
	})

	resp, err := client.Domains.DeleteRecord(ctx, "dns-zone-1", 244042432)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestDomains_ValidationErrors(t *testing.T) {
	setup()
	defer teardown()

	tests := map[string]struct {
		err    error
		target error
	}{
		"get zone empty id": {
			err: func() error {
				_, _, err := client.Domains.GetZone(ctx, "")
				return err
			}(),
			target: ErrEmptyArgument,
		},
		"create zone nil payload": {
			err: func() error {
				_, _, err := client.Domains.CreateZone(ctx, nil)
				return err
			}(),
			target: ErrEmptyPayloadNotAllowed,
		},
		"delete zone empty id": {
			err: func() error {
				_, err := client.Domains.DeleteZone(ctx, "")
				return err
			}(),
			target: ErrEmptyArgument,
		},
		"get soa empty zone id": {
			err: func() error {
				_, _, err := client.Domains.GetSOA(ctx, "")
				return err
			}(),
			target: ErrEmptyArgument,
		},
		"update soa empty zone id": {
			err: func() error {
				_, err := client.Domains.UpdateSOA(ctx, "", &DNSSOAUpdateRequest{})
				return err
			}(),
			target: ErrEmptyArgument,
		},
		"update soa nil payload": {
			err: func() error {
				_, err := client.Domains.UpdateSOA(ctx, "dns-zone-1", nil)
				return err
			}(),
			target: ErrEmptyPayloadNotAllowed,
		},
		"list records empty zone id": {
			err: func() error {
				_, _, err := client.Domains.ListRecords(ctx, "")
				return err
			}(),
			target: ErrEmptyArgument,
		},
		"create record empty zone id": {
			err: func() error {
				_, err := client.Domains.CreateRecord(ctx, "", &DNSRecordCreateRequest{})
				return err
			}(),
			target: ErrEmptyArgument,
		},
		"create record nil payload": {
			err: func() error {
				_, err := client.Domains.CreateRecord(ctx, "dns-zone-1", nil)
				return err
			}(),
			target: ErrEmptyPayloadNotAllowed,
		},
		"update record empty zone id": {
			err: func() error {
				_, err := client.Domains.UpdateRecord(ctx, "", 244042432, &DNSRecordUpdateRequest{})
				return err
			}(),
			target: ErrEmptyArgument,
		},
		"update record empty record id": {
			err: func() error {
				_, err := client.Domains.UpdateRecord(ctx, "dns-zone-1", 0, &DNSRecordUpdateRequest{})
				return err
			}(),
			target: ErrEmptyArgument,
		},
		"update record negative record id": {
			err: func() error {
				_, err := client.Domains.UpdateRecord(ctx, "dns-zone-1", -1, &DNSRecordUpdateRequest{})
				return err
			}(),
			target: ErrEmptyArgument,
		},
		"update record nil payload": {
			err: func() error {
				_, err := client.Domains.UpdateRecord(ctx, "dns-zone-1", 244042432, nil)
				return err
			}(),
			target: ErrEmptyPayloadNotAllowed,
		},
		"delete record empty zone id": {
			err: func() error {
				_, err := client.Domains.DeleteRecord(ctx, "", 244042432)
				return err
			}(),
			target: ErrEmptyArgument,
		},
		"delete record empty record id": {
			err: func() error {
				_, err := client.Domains.DeleteRecord(ctx, "dns-zone-1", 0)
				return err
			}(),
			target: ErrEmptyArgument,
		},
		"delete record negative record id": {
			err: func() error {
				_, err := client.Domains.DeleteRecord(ctx, "dns-zone-1", -1)
				return err
			}(),
			target: ErrEmptyArgument,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Error(t, test.err)
			assert.True(t, errors.Is(test.err, test.target))
		})
	}
}
