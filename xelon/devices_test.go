package xelon

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDevices_AddDisk(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/devices/abc123/disk", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		// Verify request body
		body, _ := io.ReadAll(r.Body)
		var req DeviceAddDiskRequest
		_ = json.Unmarshal(body, &req)
		assert.Equal(t, 100, req.Size)
		assert.Equal(t, false, req.IsHDD)

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, `{"success":"Disk is being added"}`)
	})

	addRequest := &DeviceAddDiskRequest{
		Size:  100,
		IsHDD: false,
	}

	resp, err := client.Devices.AddDisk(ctx, "abc123", addRequest)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDevices_AddDisk_emptyDeviceID(t *testing.T) {
	addRequest := &DeviceAddDiskRequest{
		Size: 100,
	}

	_, err := client.Devices.AddDisk(ctx, "", addRequest)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "device id must be supplied")
}

func TestDevices_AddDisk_nilRequest(t *testing.T) {
	_, err := client.Devices.AddDisk(ctx, "abc123", nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request must be supplied")
}

func TestDevices_UpdateDisk(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/devices/abc123/disk", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)

		// Verify request body
		body, _ := io.ReadAll(r.Body)
		var req DeviceUpdateDiskRequest
		_ = json.Unmarshal(body, &req)
		assert.Equal(t, "disk-1", req.DiskID)
		assert.Equal(t, 200, req.Size)
		assert.Equal(t, true, req.ExtendPartition)
		assert.Equal(t, true, req.CreateSnapshot)

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, `{"success":"Disk is being updated"}`)
	})

	updateRequest := &DeviceUpdateDiskRequest{
		DiskID:          "disk-1",
		Size:            200,
		ExtendPartition: true,
		CreateSnapshot:  true,
	}

	resp, err := client.Devices.UpdateDisk(ctx, "abc123", updateRequest)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDevices_UpdateDisk_emptyDeviceID(t *testing.T) {
	updateRequest := &DeviceUpdateDiskRequest{
		DiskID: "disk-1",
		Size:   200,
	}

	_, err := client.Devices.UpdateDisk(ctx, "", updateRequest)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "device id must be supplied")
}

func TestDevices_UpdateDisk_nilRequest(t *testing.T) {
	_, err := client.Devices.UpdateDisk(ctx, "abc123", nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request must be supplied")
}

func TestDevices_DeleteDisk(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/devices/abc123/disk", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)

		// Verify request body
		body, _ := io.ReadAll(r.Body)
		var req DeviceDeleteDiskRequest
		_ = json.Unmarshal(body, &req)
		assert.Equal(t, "disk-1", req.DiskID)
		assert.Equal(t, "password123", req.Password)

		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, `{"success":"Disk is being deleted"}`)
	})

	deleteRequest := &DeviceDeleteDiskRequest{
		DiskID:   "disk-1",
		Password: "password123",
	}

	resp, err := client.Devices.DeleteDisk(ctx, "abc123", deleteRequest)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDevices_DeleteDisk_emptyDeviceID(t *testing.T) {
	deleteRequest := &DeviceDeleteDiskRequest{
		DiskID: "disk-1",
	}

	_, err := client.Devices.DeleteDisk(ctx, "", deleteRequest)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "device id must be supplied")
}

func TestDevices_DeleteDisk_nilRequest(t *testing.T) {
	_, err := client.Devices.DeleteDisk(ctx, "abc123", nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request must be supplied")
}
