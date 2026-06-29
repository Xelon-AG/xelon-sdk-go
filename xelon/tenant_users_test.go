package xelon

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTenantUsers_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("GET /tenants/tenant-1/users", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "john", r.URL.Query().Get("search"))
		fixture := loadFixture(t, "tenantusers_list_users.json")
		_, _ = w.Write(fixture)
	})
	expectedUsers := []TenantUser{{
		Email:    "john.doe@example.com",
		ID:       "user-1",
		JobTitle: "sysadmin_devops",
		Name:     "John",
		Surname:  "Doe",
		TenantID: "tenant-1",
	}, {
		Email:    "jane.doe@example.com",
		ID:       "user-2",
		JobTitle: "ceo",
		Name:     "Jane",
		Surname:  "Doe",
		TenantID: "tenant-1",
	}}

	actualUsers, resp, err := client.TenantUsers.List(ctx, "tenant-1", &TenantUserListOptions{Search: "john"})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedUsers, actualUsers)
	assert.Equal(t, &Meta{
		Total:    3,
		LastPage: 2,
		PerPage:  2,
		Page:     1,
		From:     1,
		To:       2,
	}, resp.Meta)
}

func TestTenantUsers_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("GET /tenants/tenant-1/users/user-1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fixture := loadFixture(t, "tenantusers_get_user_success.json")
		_, _ = w.Write(fixture)
	})
	expectedUser := &TenantUser{
		Email:    "john.doe@example.com",
		ID:       "user-1",
		JobTitle: "developer",
		Name:     "John",
		Permissions: []TenantUserPermission{{
			DisplayName: "Allow view virtual machines",
			ID:          75,
			Name:        "allow_view_virtual_machines",
			Type:        "virtual_machine",
		}},
		Roles: []TenantUserRole{{
			DisplayName: "Organization Admin",
			ID:          1,
			Name:        "hq_organization_admin",
			Type:        "organization",
		}},
		Surname:  "Doe",
		TenantID: "tenant-1",
	}

	actualUser, resp, err := client.TenantUsers.Get(ctx, "tenant-1", "user-1")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedUser, actualUser)
}

func TestTenantUsers_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("POST /tenants/tenant-1/users", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.JSONEq(t, `{
			"name": "John",
			"surname": "Doe",
			"email": "john.doe@example.com",
			"password": "SecurePass123!",
			"password_confirmation": "SecurePass123!",
			"job_title": "developer",
			"passwordShouldBeChanged": false,
			"welcomeEmail": true,
			"roles": ["hq_organization_admin"],
			"permissions": ["allow_view_virtual_machines"]
		}`, string(body))

		var actualRequest TenantUserCreateRequest
		err = json.Unmarshal(body, &actualRequest)
		assert.NoError(t, err)
		assert.Equal(t, TenantUserCreateRequest{
			Email:                 "john.doe@example.com",
			JobTitle:              "developer",
			Name:                  "John",
			Password:              "SecurePass123!",
			PasswordConfirmation:  "SecurePass123!",
			Permissions:           []string{"allow_view_virtual_machines"},
			RequirePasswordChange: false,
			Roles:                 []string{"hq_organization_admin"},
			SendWelcomeEmail:      true,
			Surname:               "Doe",
		}, actualRequest)

		fixture := loadFixture(t, "tenantusers_create_user_success.json")
		_, _ = w.Write(fixture)
	})
	expectedUser := &TenantUser{
		Email:    "john.doe@example.com",
		ID:       "user-1",
		JobTitle: "developer",
		Name:     "John",
		Surname:  "Doe",
		TenantID: "tenant-1",
	}

	actualUser, resp, err := client.TenantUsers.Create(ctx, "tenant-1", &TenantUserCreateRequest{
		Email:                 "john.doe@example.com",
		JobTitle:              "developer",
		Name:                  "John",
		Password:              "SecurePass123!",
		PasswordConfirmation:  "SecurePass123!",
		Permissions:           []string{"allow_view_virtual_machines"},
		RequirePasswordChange: false,
		Roles:                 []string{"hq_organization_admin"},
		SendWelcomeEmail:      true,
		Surname:               "Doe",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedUser, actualUser)
}

func TestTenantUsers_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("PUT /tenants/tenant-1/users/user-1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)

		var actualRequest TenantUserUpdateRequest
		err := json.NewDecoder(r.Body).Decode(&actualRequest)
		assert.NoError(t, err)
		assert.Equal(t, TenantUserUpdateRequest{
			JobTitle: "developer",
			Name:     "Jane",
			Surname:  "Doe",
		}, actualRequest)

		fixture := loadFixture(t, "tenantusers_update_user_success.json")
		_, _ = w.Write(fixture)
	})
	expectedUser := &TenantUser{
		Email:    "jane.doe@example.com",
		ID:       "user-1",
		JobTitle: "developer",
		Name:     "Jane",
		Surname:  "Doe",
		TenantID: "tenant-1",
	}

	actualUser, resp, err := client.TenantUsers.Update(ctx, "tenant-1", "user-1", &TenantUserUpdateRequest{
		JobTitle: "developer",
		Name:     "Jane",
		Surname:  "Doe",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedUser, actualUser)
}

func TestTenantUsers_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("DELETE /tenants/tenant-1/users/user-1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		fixture := loadFixture(t, "tenantusers_delete_user_success.json")
		_, _ = w.Write(fixture)
	})

	resp, err := client.TenantUsers.Delete(ctx, "tenant-1", "user-1")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestTenantUsers_Restore(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("POST /tenants/tenant-1/users/restore", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.JSONEq(t, `{
			"userIdentifier": "user-1"
		}`, string(body))

		fixture := loadFixture(t, "tenantusers_restore_user_success.json")
		_, _ = w.Write(fixture)
	})

	resp, err := client.TenantUsers.Restore(ctx, "tenant-1", "user-1")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestTenantUsers_UpdatePassword(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("POST /tenants/tenant-1/users/user-1/password", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.JSONEq(t, `{
			"password": "NewSecurePass123!",
			"password_confirmation": "NewSecurePass123!"
		}`, string(body))

		fixture := loadFixture(t, "tenantusers_update_password_success.json")
		_, _ = w.Write(fixture)
	})

	resp, err := client.TenantUsers.UpdatePassword(ctx, "tenant-1", "user-1", &TenantUserPasswordUpdateRequest{
		Password:             "NewSecurePass123!",
		PasswordConfirmation: "NewSecurePass123!",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestTenantUsers_ListAvailablePermissions(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("GET /tenants/tenant-1/users/permissions", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fixture := loadFixture(t, "tenantusers_list_available_permissions.json")
		_, _ = w.Write(fixture)
	})
	expectedPermissions := []TenantUserPermission{{
		DisplayName: "View Virtual Machines",
		ID:          1,
		Name:        "allow_view_virtual_machines",
		Type:        "virtual_machines",
	}}

	actualPermissions, resp, err := client.TenantUsers.ListAvailablePermissions(ctx, "tenant-1")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedPermissions, actualPermissions)
}

func TestTenantUsers_UpdatePermissions(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("POST /tenants/tenant-1/users/user-1/permissions", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		var actualRequest TenantUserPermissionsUpdateRequest
		err := json.NewDecoder(r.Body).Decode(&actualRequest)
		assert.NoError(t, err)
		assert.Equal(t, TenantUserPermissionsUpdateRequest{
			ChildTenants: []string{"child-tenant-1"},
			Permissions:  []string{"allow_view_virtual_machines"},
			Roles:        []string{"hq_organization_admin"},
		}, actualRequest)

		fixture := loadFixture(t, "tenantusers_update_permissions_success.json")
		_, _ = w.Write(fixture)
	})

	resp, err := client.TenantUsers.UpdatePermissions(ctx, "tenant-1", "user-1", &TenantUserPermissionsUpdateRequest{
		ChildTenants: []string{"child-tenant-1"},
		Permissions:  []string{"allow_view_virtual_machines"},
		Roles:        []string{"hq_organization_admin"},
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestTenantUsers_MissingData(t *testing.T) {
	setup()
	defer teardown()

	tests := map[string]struct {
		request func() (*TenantUser, *Response, error)
	}{
		"get": {
			request: func() (*TenantUser, *Response, error) {
				return client.TenantUsers.Get(ctx, "tenant-1", "user-1")
			},
		},
		"create": {
			request: func() (*TenantUser, *Response, error) {
				return client.TenantUsers.Create(ctx, "tenant-1", &TenantUserCreateRequest{})
			},
		},
		"update": {
			request: func() (*TenantUser, *Response, error) {
				return client.TenantUsers.Update(ctx, "tenant-1", "user-1", &TenantUserUpdateRequest{})
			},
		},
	}

	mux.HandleFunc("GET /tenants/tenant-1/users/user-1", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"message":"missing data"}`))
	})
	mux.HandleFunc("POST /tenants/tenant-1/users", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"data":null,"message":"missing data"}`))
	})
	mux.HandleFunc("PUT /tenants/tenant-1/users/user-1", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"data":null,"message":"missing data"}`))
	})

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actualUser, resp, err := test.request()

			assert.Nil(t, actualUser)
			assert.NotNil(t, resp)
			assert.EqualError(t, err, "tenant user data is empty")
		})
	}
}

func TestTenantUsers_ValidationErrors(t *testing.T) {
	client := NewClient("auth-token")

	tests := map[string]struct {
		request func() error
		target  error
	}{
		"list empty tenant id": {
			request: func() error {
				_, _, err := client.TenantUsers.List(ctx, "", nil)
				return err
			},
			target: ErrEmptyArgument,
		},
		"get empty tenant id": {
			request: func() error {
				_, _, err := client.TenantUsers.Get(ctx, "", "user-1")
				return err
			},
			target: ErrEmptyArgument,
		},
		"get empty user id": {
			request: func() error {
				_, _, err := client.TenantUsers.Get(ctx, "tenant-1", "")
				return err
			},
			target: ErrEmptyArgument,
		},
		"create empty tenant id": {
			request: func() error {
				_, _, err := client.TenantUsers.Create(ctx, "", &TenantUserCreateRequest{})
				return err
			},
			target: ErrEmptyArgument,
		},
		"create nil payload": {
			request: func() error {
				_, _, err := client.TenantUsers.Create(ctx, "tenant-1", nil)
				return err
			},
			target: ErrEmptyPayloadNotAllowed,
		},
		"update empty tenant id": {
			request: func() error {
				_, _, err := client.TenantUsers.Update(ctx, "", "user-1", &TenantUserUpdateRequest{})
				return err
			},
			target: ErrEmptyArgument,
		},
		"update empty user id": {
			request: func() error {
				_, _, err := client.TenantUsers.Update(ctx, "tenant-1", "", &TenantUserUpdateRequest{})
				return err
			},
			target: ErrEmptyArgument,
		},
		"update nil payload": {
			request: func() error {
				_, _, err := client.TenantUsers.Update(ctx, "tenant-1", "user-1", nil)
				return err
			},
			target: ErrEmptyPayloadNotAllowed,
		},
		"delete empty tenant id": {
			request: func() error {
				_, err := client.TenantUsers.Delete(ctx, "", "user-1")
				return err
			},
			target: ErrEmptyArgument,
		},
		"delete empty user id": {
			request: func() error {
				_, err := client.TenantUsers.Delete(ctx, "tenant-1", "")
				return err
			},
			target: ErrEmptyArgument,
		},
		"restore empty tenant id": {
			request: func() error {
				_, err := client.TenantUsers.Restore(ctx, "", "user-1")
				return err
			},
			target: ErrEmptyArgument,
		},
		"restore empty user id": {
			request: func() error {
				_, err := client.TenantUsers.Restore(ctx, "tenant-1", "")
				return err
			},
			target: ErrEmptyArgument,
		},
		"update password empty tenant id": {
			request: func() error {
				_, err := client.TenantUsers.UpdatePassword(ctx, "", "user-1", &TenantUserPasswordUpdateRequest{})
				return err
			},
			target: ErrEmptyArgument,
		},
		"update password empty user id": {
			request: func() error {
				_, err := client.TenantUsers.UpdatePassword(ctx, "tenant-1", "", &TenantUserPasswordUpdateRequest{})
				return err
			},
			target: ErrEmptyArgument,
		},
		"update password nil payload": {
			request: func() error {
				_, err := client.TenantUsers.UpdatePassword(ctx, "tenant-1", "user-1", nil)
				return err
			},
			target: ErrEmptyPayloadNotAllowed,
		},
		"list available permissions empty tenant id": {
			request: func() error {
				_, _, err := client.TenantUsers.ListAvailablePermissions(ctx, "")
				return err
			},
			target: ErrEmptyArgument,
		},
		"update permissions empty tenant id": {
			request: func() error {
				_, err := client.TenantUsers.UpdatePermissions(ctx, "", "user-1", &TenantUserPermissionsUpdateRequest{})
				return err
			},
			target: ErrEmptyArgument,
		},
		"update permissions empty user id": {
			request: func() error {
				_, err := client.TenantUsers.UpdatePermissions(ctx, "tenant-1", "", &TenantUserPermissionsUpdateRequest{})
				return err
			},
			target: ErrEmptyArgument,
		},
		"update permissions nil payload": {
			request: func() error {
				_, err := client.TenantUsers.UpdatePermissions(ctx, "tenant-1", "user-1", nil)
				return err
			},
			target: ErrEmptyPayloadNotAllowed,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.request()

			assert.Error(t, err)
			assert.True(t, errors.Is(err, test.target))
		})
	}
}
