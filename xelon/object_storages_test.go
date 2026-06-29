package xelon

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObjectStorages_ListUsers(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("GET /object-storages/users", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fixture := loadFixture(t, "objectstorages_list_users.json")
		_, _ = w.Write(fixture)
	})
	expectedUsers := []ObjectStorageUser{{
		ID:                     "00000000-0000-0000-0000-000000000000",
		Name:                   "test-user-0",
		QuotaGB:                200,
		S3Endpoints:            []string{"https://ch1-s3.xelon.io"},
		Tenant:                 &Tenant{ID: "000000000", Name: "test-tenant-0"},
		UsedGB:                 10.4,
		ZoneReplicationEnabled: true,
	}, {
		ID:                     "11111111-1111-1111-1111-111111111111",
		Name:                   "test-user-1",
		QuotaGB:                500,
		S3Endpoints:            []string{"https://zh1-s3.xelon.io"},
		Tenant:                 &Tenant{ID: "111111111", Name: "test-tenant-1"},
		UsedGB:                 0,
		ZoneReplicationEnabled: false,
	}}

	actualUsers, resp, err := client.ObjectStorages.ListUsers(ctx, nil)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedUsers, actualUsers)
}

func TestObjectStorages_CreateUser(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("POST /object-storages/users", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		fixture := loadFixture(t, "objectstorages_create_user_success.json")
		_, _ = w.Write(fixture)
	})
	expectedUser := &ObjectStorageUser{
		ID:   "00000000-0000-0000-0000-000000000000",
		Name: "test-user-0",
		Tokens: []ObjectStorageUserToken{{
			AccessKey: "ak_test_1234567890",
			CreatedAt: mustTime(t, "2025-10-27T14:19:56+01:00"),
			ID:        "000000000000",
			SecretKey: "sk_test_1234567890abcdef",
		}},
	}

	actualUser, resp, err := client.ObjectStorages.CreateUser(ctx, &ObjectStorageUserCreateRequest{
		Name:    "test-user-0",
		QuotaGB: 200,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedUser, actualUser)
}

func TestObjectStorages_CreateUser_MissingData(t *testing.T) {
	setup()
	defer teardown()

	type testCase struct {
		responseBody string
	}
	tests := map[string]testCase{
		"missing data": {
			responseBody: `{"message":"S3 user successfully created"}`,
		},
		"null data": {
			responseBody: `{"data":null,"message":"S3 user successfully created"}`,
		},
	}

	var responseBody string
	mux.HandleFunc("POST /object-storages/users", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		_, _ = w.Write([]byte(responseBody))
	})

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			responseBody = test.responseBody

			actualUser, resp, err := client.ObjectStorages.CreateUser(ctx, &ObjectStorageUserCreateRequest{
				Name:    "test-user-0",
				QuotaGB: 200,
			})

			assert.Nil(t, actualUser)
			assert.NotNil(t, resp)
			assert.EqualError(t, err, "failed to create object storage user: response data is empty")
		})
	}
}

func TestObjectStorages_UpdateUser(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("PUT /object-storages/users/00000000-0000-0000-0000-000000000000", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		fixture := loadFixture(t, "objectstorages_update_user_success.json")
		_, _ = w.Write(fixture)
	})
	expectedUser := &ObjectStorageUser{
		ID:   "00000000-0000-0000-0000-000000000000",
		Name: "test-user-0__updated",
	}

	actualUser, resp, err := client.ObjectStorages.UpdateUser(ctx, "00000000-0000-0000-0000-000000000000", &ObjectStorageUserUpdateRequest{
		Name:    "test-user-0__updated",
		QuotaGB: 200,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedUser, actualUser)
}

func TestObjectStorages_UpdateUser_MissingData(t *testing.T) {
	setup()
	defer teardown()

	type testCase struct {
		responseBody string
	}
	tests := map[string]testCase{
		"missing data": {
			responseBody: `{"message":"S3 user successfully edited"}`,
		},
		"null data": {
			responseBody: `{"data":null,"message":"S3 user successfully edited"}`,
		},
	}

	var responseBody string
	mux.HandleFunc("PUT /object-storages/users/00000000-0000-0000-0000-000000000000", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		_, _ = w.Write([]byte(responseBody))
	})

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			responseBody = test.responseBody

			actualUser, resp, err := client.ObjectStorages.UpdateUser(ctx, "00000000-0000-0000-0000-000000000000", &ObjectStorageUserUpdateRequest{
				Name:    "test-user-0__updated",
				QuotaGB: 200,
			})

			assert.Nil(t, actualUser)
			assert.NotNil(t, resp)
			assert.EqualError(t, err, "failed to update object storage user: response data is empty")
		})
	}
}

func TestObjectStorages_CreateUserToken(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("POST /object-storages/users/00000000-0000-0000-0000-000000000000/tokens", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		fixture := loadFixture(t, "objectstorages_create_user_token_success.json")
		_, _ = w.Write(fixture)
	})
	expectedToken := &ObjectStorageUserToken{
		AccessKey: "ak_test_1234567890",
		CreatedAt: mustTime(t, "2025-10-27T14:19:56+01:00"),
		ID:        "000000000000",
		SecretKey: "sk_test_1234567890abcdef",
	}

	actualToken, resp, err := client.ObjectStorages.CreateUserToken(ctx, "00000000-0000-0000-0000-000000000000")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedToken, actualToken)
}

func TestObjectStorages_CreateUserToken_MissingData(t *testing.T) {
	setup()
	defer teardown()

	type testCase struct {
		responseBody string
	}
	tests := map[string]testCase{
		"missing data": {
			responseBody: `{"message":"S3 user token successfully created"}`,
		},
		"null data": {
			responseBody: `{"data":null,"message":"S3 user token successfully created"}`,
		},
	}

	var responseBody string
	mux.HandleFunc("POST /object-storages/users/00000000-0000-0000-0000-000000000000/tokens", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		_, _ = w.Write([]byte(responseBody))
	})

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			responseBody = test.responseBody

			actualToken, resp, err := client.ObjectStorages.CreateUserToken(ctx, "00000000-0000-0000-0000-000000000000")

			assert.Nil(t, actualToken)
			assert.NotNil(t, resp)
			assert.EqualError(t, err, "failed to create user token: response data is empty")
		})
	}
}
