package xelon

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBackups_ListPlans(t *testing.T) {
	setup(t)
	defer teardown()

	mux.HandleFunc("GET /backups/tags", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "cloud-1", r.URL.Query().Get("cloudIdentifier"))
		assert.Len(t, r.URL.Query(), 1)

		fixture := loadFixture(t, "backups_list_plans.json")
		_, _ = w.Write(fixture)
	})

	plans, resp, err := client.Backups.ListPlans(ctx, "cloud-1")

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Nil(t, resp.Meta)
	assert.Equal(t, []BackupPlan{
		{ID: 17, Name: "Daily Backup"},
		{ID: 23, Name: "Weekly Backup"},
	}, plans)
}

func TestBackups_ListPlans_MissingData(t *testing.T) {
	tests := map[string]string{
		"missing data": `{"message":"Backup tags retrieved."}`,
		"null data":    `{"data":null}`,
	}

	for name, responseBody := range tests {
		t.Run(name, func(t *testing.T) {
			setup(t)
			defer teardown()

			mux.HandleFunc("GET /backups/tags", func(w http.ResponseWriter, _ *http.Request) {
				_, _ = io.WriteString(w, responseBody)
			})

			plans, resp, err := client.Backups.ListPlans(ctx, "cloud-1")

			assert.Nil(t, plans)
			assert.NotNil(t, resp)
			assert.EqualError(t, err, "failed to list backup plans: response data is empty")
		})
	}
}

func TestBackups_GetDevicePlan(t *testing.T) {
	tests := map[string]struct {
		responseBody string
		want         *BackupPlan
		wantErr      error
	}{
		"no assigned plan": {
			responseBody: `{"data":[]}`,
		},
		"one assigned plan": {
			responseBody: string(loadFixture(t, "backups_get_device_plan_success.json")),
			want:         &BackupPlan{ID: 17, Name: "Daily Backup"},
		},
		"multiple assigned plans": {
			responseBody: `{"data":[{"id":17,"name":"Daily Backup"},{"id":23,"name":"Weekly Backup"}]}`,
			wantErr:      ErrAmbiguousBackupPlanAssignment,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			setup(t)
			defer teardown()

			mux.HandleFunc("GET /devices/device-1/backups/tags", func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				_, _ = io.WriteString(w, test.responseBody)
			})

			plan, resp, err := client.Backups.GetDevicePlan(ctx, "device-1")

			if test.wantErr == nil {
				require.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, test.wantErr)
			}
			assert.NotNil(t, resp)
			assert.Equal(t, test.want, plan)
		})
	}
}

func TestBackups_GetDevicePlan_MissingData(t *testing.T) {
	tests := map[string]string{
		"missing data": `{"message":"Backup tags retrieved."}`,
		"null data":    `{"data":null}`,
	}

	for name, responseBody := range tests {
		t.Run(name, func(t *testing.T) {
			setup(t)
			defer teardown()

			mux.HandleFunc("GET /devices/device-1/backups/tags", func(w http.ResponseWriter, _ *http.Request) {
				_, _ = io.WriteString(w, responseBody)
			})

			plan, resp, err := client.Backups.GetDevicePlan(ctx, "device-1")

			assert.Nil(t, plan)
			assert.NotNil(t, resp)
			assert.EqualError(t, err, "failed to get device backup plan: response data is empty")
		})
	}
}

func TestBackups_SetDevicePlan(t *testing.T) {
	setup(t)
	defer teardown()

	mux.HandleFunc("PUT /devices/device-1/backups/tags", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)

		var payload map[string]int
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.NoError(t, err)
		assert.Equal(t, map[string]int{"backupId": 17}, payload)

		_, _ = io.WriteString(w, `{"message":"Backup has been changed."}`)
	})

	resp, err := client.Backups.SetDevicePlan(ctx, "device-1", 17)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestBackups_DisableDeviceBackup(t *testing.T) {
	setup(t)
	defer teardown()

	mux.HandleFunc("DELETE /devices/device-1/backups/tags", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)

		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.Empty(t, body)

		_, _ = io.WriteString(w, `{"message":"Backup has been deleted."}`)
	})

	resp, err := client.Backups.DisableDeviceBackup(ctx, "device-1")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestBackups_ValidationErrors(t *testing.T) {
	c := newTestClient(t)
	tests := map[string]error{
		"list plans empty cloud id": func() error {
			_, _, err := c.Backups.ListPlans(ctx, "")
			return err
		}(),
		"get device plan empty device id": func() error {
			_, _, err := c.Backups.GetDevicePlan(ctx, "")
			return err
		}(),
		"set device plan empty device id": func() error {
			_, err := c.Backups.SetDevicePlan(ctx, "", 17)
			return err
		}(),
		"set device plan empty backup plan id": func() error {
			_, err := c.Backups.SetDevicePlan(ctx, "device-1", 0)
			return err
		}(),
		"set device plan negative backup plan id": func() error {
			_, err := c.Backups.SetDevicePlan(ctx, "device-1", -1)
			return err
		}(),
		"disable device backup empty device id": func() error {
			_, err := c.Backups.DisableDeviceBackup(ctx, "")
			return err
		}(),
	}

	for name, err := range tests {
		t.Run(name, func(t *testing.T) {
			assert.ErrorIs(t, err, ErrEmptyArgument)
		})
	}
}
