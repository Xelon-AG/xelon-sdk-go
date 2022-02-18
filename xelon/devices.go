package xelon

import (
	"context"
	"fmt"
	"net/http"
)

// DevicesService handles communication with the devices related methods of the Xelon API.
type DevicesService service

// Device represents a Xelon device.
type Device struct {
	CPU            int            `json:"cpu"`
	LocalVMDetails LocalVMDetails `json:"localvmdetails,omitempty"`
	Networks       []Network      `json:"networks,omitempty"`
	PowerState     bool           `json:"powerstate"`
	RAM            int            `json:"ram"`
}

// LocalVMDetails represents a Xelon device's details.
type LocalVMDetails struct {
	CreatedAt     string `json:"created_at"`
	HVSystemID    int    `json:"hv_system_id"`
	ISOMounted    string `json:"iso_mounted,omitempty"`
	LocalVMID     string `json:"localvmid"`
	State         int    `json:"state"`
	TemplateID    int    `json:"template_id"`
	UpdatedAt     string `json:"updated_at"`
	UserID        int    `json:"user_id"`
	VMDisplayName string `json:"vmdisplayname"`
	VMHostname    string `json:"vmhostname"`
}

// Network represents a Xelon device's network information.
type Network struct {
	IPAddress  string `json:"ip,omitempty"`
	Label      string `json:"label,omitempty"`
	MacAddress string `json:"macAddress,omitempty"`
}

// DeviceRoot represents a Xelon device root object.
type DeviceRoot struct {
	Device Device `json:"device,omitempty"`
}

// Get provides detailed information for a device identified by tenant and localvmid.
func (s *DevicesService) Get(ctx context.Context, tenantID, localVMID string) (*DeviceRoot, *Response, error) {
	if tenantID == "" || localVMID == "" {
		return nil, nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("device?tenant=%v&localvmid=%v", tenantID, localVMID)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	deviceRoot := new(DeviceRoot)
	resp, err := s.client.Do(ctx, req, deviceRoot)
	if err != nil {
		return nil, resp, err
	}

	return deviceRoot, resp, nil
}
