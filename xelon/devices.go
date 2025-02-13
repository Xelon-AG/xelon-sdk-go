package xelon

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

const deviceBasePath = "devices"

// DevicesService handles communication with the devices related methods of the Xelon REST API.
type DevicesService service

// Device represents a Xelon device (virtual machine).
type Device struct {
	CPUCores          int    `json:"cpu,omitempty"`
	DisplayName       string `json:"displayName,omitempty"`
	HostName          string `json:"hostName,omitempty"`
	ID                string `json:"identifier,omitempty"`
	MonitoringEnabled bool   `json:"monitoring,omitempty"`
	PoweredOn         bool   `json:"isPoweredOn,omitempty"`
	RAM               int    `json:"ram,omitempty"`
	State             int    `json:"state,omitempty"`
	TemplateID        string `json:"templateId,omitempty"`
	TenantID          string `json:"tenantIdentifier,omitempty"`
}

type DeviceCreateRequest struct {
	BackupJobID          int                   `json:"backJobId,omitempty"`
	CPUCores             int                   `json:"cpu"`
	DiskSize             int                   `json:"diskSize"`
	DisplayName          string                `json:"displayName"`
	HostName             string                `json:"hostName"`
	EnableMonitoring     bool                  `json:"isMonitoring"`
	Networks             []DeviceCreateNetwork `json:"networks,omitempty"`
	Password             string                `json:"password"`
	PasswordConfirmation string                `json:"passwordConfirmation"`
	RAM                  int                   `json:"ram"`
	ScriptID             string                `json:"scriptId,omitempty"`
	SendEmail            bool                  `json:"sendEmail,omitempty"`
	SSHKeyID             string                `json:"sshKeyId,omitempty"`
	SwapDiskSize         int                   `json:"swapDiskSize"`
	TemplateID           string                `json:"templateId"`
	TenantID             string                `json:"tenantIdentifier"`
}

type DeviceCreateNetwork struct {
	ConnectOnPowerOn bool   `json:"connectOnPowerOn"`
	IPAddress        string `json:"ip,omitempty"`
	IPAddressID      string `json:"ipId,omitempty"`
	NetworkID        string `json:"networkId"`
}

// DeviceListOptions specifies the optional parameters to the DevicesService.List.
type DeviceListOptions struct {
	Sort   string `url:"sort,omitempty"`
	Search string `url:"search,omitempty"`

	ListOptions
}

type deviceRoot struct {
	Device  *Device `json:"data,omitempty"`
	Message string  `json:"message,omitempty"`
}

type devicesRoot struct {
	Devices []Device `json:"data"`
	Meta    *Meta    `json:"meta,omitempty"`
}

func (v Device) String() string { return Stringify(v) }

// List provides a list of all devices.
func (s *DevicesService) List(ctx context.Context, opts *DeviceListOptions) ([]Device, *Response, error) {
	path, err := addOptions(deviceBasePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(devicesRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}

	return root.Devices, resp, nil
}

// Get provides detailed information for device identified by id.
func (s *DevicesService) Get(ctx context.Context, deviceID string) (*Device, *Response, error) {
	if deviceID == "" {
		return nil, nil, errors.New("failed to get device: id must be supplied")
	}

	path := fmt.Sprintf("%v/%v", deviceBasePath, deviceID)
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	device := new(Device)
	resp, err := s.client.Do(ctx, req, device)
	if err != nil {
		return nil, resp, err
	}

	return device, resp, err
}

// Create makes a device with given payload.
func (s *DevicesService) Create(ctx context.Context, createRequest *DeviceCreateRequest) (*Device, *Response, error) {
	if createRequest == nil {
		return nil, nil, errors.New("failed to create device: payload must be supplied")
	}

	req, err := s.client.NewRequest(http.MethodPost, deviceBasePath, createRequest)
	if err != nil {
		return nil, nil, err
	}

	deviceRoot := new(deviceRoot)
	resp, err := s.client.Do(ctx, req, deviceRoot)
	if err != nil {
		return nil, resp, err
	}

	return deviceRoot.Device, resp, nil
}

// Delete removes device identified by id.
func (s *DevicesService) Delete(ctx context.Context, deviceID string) (*Response, error) {
	if deviceID == "" {
		return nil, errors.New("failed to delete device: id must be supplied")
	}

	path := fmt.Sprintf("%v/%v", deviceBasePath, deviceID)
	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Start sends 'start' action and starts device identified by id.
func (s *DevicesService) Start(ctx context.Context, deviceID string) (*Response, error) {
	if deviceID == "" {
		return nil, errors.New("failed to start device: id must be supplied")
	}

	path := fmt.Sprintf("%v/%v/start", deviceBasePath, deviceID)
	req, err := s.client.NewRequest(http.MethodPost, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// Stop sends an ACPI shutdown to device identified by id.
func (s *DevicesService) Stop(ctx context.Context, deviceID string) (*Response, error) {
	if deviceID == "" {
		return nil, errors.New("failed to stop device: id must be supplied")
	}

	path := fmt.Sprintf("%v/%v/stop", deviceBasePath, deviceID)
	req, err := s.client.NewRequest(http.MethodPost, path, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
