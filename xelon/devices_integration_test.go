package xelon

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDevices_AddDisk_Integration tests adding a disk to a real device.
// Set XELON_INTEGRATION_TEST=1 and XELON_TEST_DEVICE_ID=<device-id> to run this test.
func TestDevices_AddDisk_Integration(t *testing.T) {
	if os.Getenv("XELON_INTEGRATION_TEST") != "1" {
		t.Skip("Skipping integration test. Set XELON_INTEGRATION_TEST=1 to run.")
	}

	deviceID := os.Getenv("XELON_TEST_DEVICE_ID")
	if deviceID == "" {
		t.Fatal("XELON_TEST_DEVICE_ID must be set for integration tests")
	}

	// Setup client from environment
	baseURL := os.Getenv("XELON_BASE_URL")
	token := os.Getenv("XELON_TOKEN")
	clientID := os.Getenv("XELON_CLIENT_ID")

	require.NotEmpty(t, baseURL, "XELON_BASE_URL must be set")
	require.NotEmpty(t, token, "XELON_TOKEN must be set")
	require.NotEmpty(t, clientID, "XELON_CLIENT_ID must be set")

	client := NewClient(token,
		WithBaseURL(baseURL),
		WithClientID(clientID),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Test AddDisk
	t.Run("AddDisk", func(t *testing.T) {
		addRequest := &DeviceAddDiskRequest{
			Size:  10, // Small disk for testing
			IsHDD: false,
		}

		t.Logf("Adding 10GB disk to device %s", deviceID)
		resp, err := client.Devices.AddDisk(ctx, deviceID, addRequest)

		assert.NoError(t, err)
		if err != nil {
			t.Logf("Error response: %+v", resp)
		}
		assert.NotNil(t, resp)

		if resp != nil {
			t.Logf("Response status: %d", resp.StatusCode)
			assert.True(t, resp.StatusCode >= 200 && resp.StatusCode < 300, "Expected 2xx status code")
		}
	})
}

// TestDevices_UpdateDisk_Integration tests updating a disk on a real device.
// Set XELON_INTEGRATION_TEST=1, XELON_TEST_DEVICE_ID=<device-id>, and XELON_TEST_DISK_ID=<disk-id> to run this test.
func TestDevices_UpdateDisk_Integration(t *testing.T) {
	if os.Getenv("XELON_INTEGRATION_TEST") != "1" {
		t.Skip("Skipping integration test. Set XELON_INTEGRATION_TEST=1 to run.")
	}

	deviceID := os.Getenv("XELON_TEST_DEVICE_ID")
	diskID := os.Getenv("XELON_TEST_DISK_ID")

	if deviceID == "" {
		t.Fatal("XELON_TEST_DEVICE_ID must be set for integration tests")
	}
	if diskID == "" {
		t.Fatal("XELON_TEST_DISK_ID must be set for integration tests")
	}

	// Setup client from environment
	baseURL := os.Getenv("XELON_BASE_URL")
	token := os.Getenv("XELON_TOKEN")
	clientID := os.Getenv("XELON_CLIENT_ID")

	require.NotEmpty(t, baseURL, "XELON_BASE_URL must be set")
	require.NotEmpty(t, token, "XELON_TOKEN must be set")
	require.NotEmpty(t, clientID, "XELON_CLIENT_ID must be set")

	client := NewClient(token,
		WithBaseURL(baseURL),
		WithClientID(clientID),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Test UpdateDisk
	t.Run("UpdateDisk", func(t *testing.T) {
		updateRequest := &DeviceUpdateDiskRequest{
			DiskID:          diskID,
			Size:            20, // Increase to 20GB
			ExtendPartition: false,
			CreateSnapshot:  false,
		}

		t.Logf("Updating disk %s on device %s to 20GB", diskID, deviceID)
		resp, err := client.Devices.UpdateDisk(ctx, deviceID, updateRequest)

		assert.NoError(t, err)
		if err != nil {
			t.Logf("Error response: %+v", resp)
		}
		assert.NotNil(t, resp)

		if resp != nil {
			t.Logf("Response status: %d", resp.StatusCode)
			assert.True(t, resp.StatusCode >= 200 && resp.StatusCode < 300, "Expected 2xx status code")
		}
	})
}
