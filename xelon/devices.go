package xelon

import (
	"context"
	"fmt"
	"net/http"
)

const devicesBasePath = "vmlist"

// DevicesService handles communication with the devices related methods of the Xelon API.
type DevicesService service

// Device represents a Xelon device.
type Device struct {
	CPU            int                   `json:"cpu"`
	LocalVMDetails *DeviceLocalVMDetails `json:"localvmdetails,omitempty"`
	Networks       []DeviceNetwork       `json:"networks,omitempty"`
	PowerState     bool                  `json:"powerstate"`
	RAM            int                   `json:"ram"`
}

// DeviceLocalVMDetails represents a Xelon device's details.
type DeviceLocalVMDetails struct {
	CreatedAt     string `json:"created_at,omitempty"`
	CPU           int    `json:"cpu,omitempty"`
	HVSystemID    int    `json:"hv_system_id,omitempty"`
	ISOMounted    string `json:"iso_mounted,omitempty"`
	LocalVMID     string `json:"localvmid,omitempty"`
	Memory        int    `json:"memory,omitempty"`
	State         int    `json:"state,omitempty"`
	TemplateID    int    `json:"template_id,omitempty"`
	UpdatedAt     string `json:"updated_at,omitempty"`
	UserID        int    `json:"user_id,omitempty"`
	VMDisplayName string `json:"vmdisplayname,omitempty"`
	VMHostname    string `json:"vmhostname,omitempty"`
}

// DeviceNetwork represents a Xelon device's network information.
type DeviceNetwork struct {
	IPAddress  string `json:"ip,omitempty"`
	Label      string `json:"label,omitempty"`
	MacAddress string `json:"macAddress,omitempty"`
}

type ToolsStatus struct {
	RunningStatus string `json:"runningStatus,omitempty"`
	Version       string `json:"version,omitempty"`
	ToolsStatus   bool   `json:"toolsStatus,omitempty"`
}

type DeviceCreateRequest struct {
	CloudID              int    `json:"cloudId"`
	CPUCores             int    `json:"cpucores"`
	DiskSize             int    `json:"disksize"`
	DisplayName          string `json:"displayname"`
	Hostname             string `json:"hostname"`
	IPAddressID          int    `json:"ipaddr1"`
	Memory               int    `json:"memory"`
	NetworkID            int    `json:"networkid1"`
	NICControllerKey     int    `json:"niccontrollerkey1"`
	NICKey               int    `json:"nickey1"`
	NICNumber            int    `json:"nicnumber"`
	NICUnit              int    `json:"nicunit1"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
	SwapDiskSize         int    `json:"swapdisksize"`
	TemplateID           int    `json:"template"`
	TenantID             string `json:"tenant_identifier"`
}

type DeviceCreationInfo struct {
	Template *Template `json:"template,omitempty"`
	NICs     []NIC     `json:"nics,omitempty"`
}

type DeviceCreateResponse struct {
	LocalVMDetails *DeviceLocalVMDetails `json:"device,omitempty"`
	IPs            []string              `json:"ips,omitempty"`
}

type DeviceListOptions struct {
	ListOptions
}

type deviceListRoot struct {
	Devices []DeviceLocalVMDetails `json:"data,omitempty"`
	Meta
}

// DeviceRoot represents a Xelon device root object.
type DeviceRoot struct {
	Device      *Device      `json:"device,omitempty"`
	ToolsStatus *ToolsStatus `json:"toolsStatus,omitempty"`
}

// List provides a list of all devices.
func (s *DevicesService) List(ctx context.Context, tenantID string, opts *DeviceListOptions) ([]DeviceLocalVMDetails, *Response, error) {
	if tenantID == "" {
		return nil, nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/devices", tenantID)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(deviceListRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	resp.Meta = &Meta{
		Page:    root.Page,
		PerPage: root.PerPage,
		Total:   root.Total,
	}

	return root.Devices, resp, nil
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

// Create makes a new device with given payload.
func (s *DevicesService) Create(ctx context.Context, createRequest *DeviceCreateRequest) (*DeviceCreateResponse, *Response, error) {
	if createRequest == nil {
		return nil, nil, ErrEmptyPayloadNotAllowed
	}

	path := fmt.Sprintf("%v/create", devicesBasePath)

	req, err := s.client.NewRequest(http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	deviceCreateResponse := new(DeviceCreateResponse)
	resp, err := s.client.Do(ctx, req, deviceCreateResponse)
	if err != nil {
		return nil, resp, err
	}

	return deviceCreateResponse, resp, nil
}

// Delete removes a device identified by localvmid.
func (s *DevicesService) Delete(ctx context.Context, localVMID string) (*Response, error) {
	if localVMID == "" {
		return nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v", devicesBasePath, localVMID)

	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Start starts a specific device identified by localvmid.
func (s *DevicesService) Start(ctx context.Context, localVMID string) (*Response, error) {
	if localVMID == "" {
		return nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v/startserver", devicesBasePath, localVMID)

	req, err := s.client.NewRequest(http.MethodPost, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Stop stops a specific device identified by localvmid.
func (s *DevicesService) Stop(ctx context.Context, localVMID string) (*Response, error) {
	if localVMID == "" {
		return nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v/stopserver", devicesBasePath, localVMID)

	req, err := s.client.NewRequest(http.MethodPost, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// GetDeviceCreationInfo retrieves a list of available templates, NICs,
// and scripts when creating a new device.
func (s *DevicesService) GetDeviceCreationInfo(ctx context.Context, tenantID, deviceType string, templateID int) (*DeviceCreationInfo, *Response, error) {
	if tenantID == "" {
		return nil, nil, ErrEmptyArgument
	}
	if deviceType == "" {
		return nil, nil, ErrEmptyArgument
	}

	path := fmt.Sprintf("%v/%v/create/Server/%v/%v", tenantID, devicesBasePath, deviceType, templateID)

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	creationInfo := new(DeviceCreationInfo)
	resp, err := s.client.Do(ctx, req, creationInfo)
	if err != nil {
		return nil, resp, err
	}

	return creationInfo, resp, nil
}
