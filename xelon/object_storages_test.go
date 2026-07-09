package xelon

import (
	"encoding/json"
	"errors"
	"io"
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
		ID:                       "00000000-0000-0000-0000-000000000000",
		Name:                     "test-user-0",
		QuotaGB:                  200,
		RegionReplicationEnabled: true,
		S3Endpoints:              []string{"https://ch1-s3.xelon.io"},
		Tenant:                   &Tenant{ID: "000000000", Name: "test-tenant-0"},
		UsedGB:                   10.4,
	}, {
		ID:                       "11111111-1111-1111-1111-111111111111",
		Name:                     "test-user-1",
		QuotaGB:                  500,
		RegionReplicationEnabled: false,
		S3Endpoints:              []string{"https://zh1-s3.xelon.io"},
		Tenant:                   &Tenant{ID: "111111111", Name: "test-tenant-1"},
		UsedGB:                   0,
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

func TestObjectStorages_ListBuckets(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("GET /object-storages/buckets", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "000000000", r.URL.Query().Get("tenantId"))
		assert.Equal(t, "2", r.URL.Query().Get("page"))
		assert.Equal(t, "5", r.URL.Query().Get("perPage"))
		fixture := loadFixture(t, "objectstorages_list_buckets.json")
		_, _ = w.Write(fixture)
	})
	expectedBuckets := []ObjectStorageBucket{{
		CreatedAt:                mustTime(t, "2025-10-27T14:19:56+01:00"),
		ID:                       "zone1.4711.1",
		IPRestrictionsEnabled:    true,
		Name:                     "test-bucket-0",
		ObjectLockEnabled:        false,
		ObjectLockRetentionDays:  0,
		ObjectStorageUserID:      "00000000-0000-0000-0000-000000000000",
		ObjectStorageUserName:    "test-user-0",
		RegionName:               "Aargau",
		RegionReplicationEnabled: true,
		S3Endpoints:              []string{"https://ch1-s3.xelon.io"},
		Tenant:                   &Tenant{ID: "000000000", Name: "test-tenant-0"},
		VersioningEnabled:        true,
	}, {
		CreatedAt:                mustTime(t, "2025-10-29T14:19:56+01:00"),
		ID:                       "zone2.4711.2",
		IPRestrictionsEnabled:    false,
		Name:                     "test-bucket-1",
		ObjectLockEnabled:        true,
		ObjectLockRetentionDays:  90,
		ObjectStorageUserID:      "11111111-1111-1111-1111-111111111111",
		ObjectStorageUserName:    "test-user-1",
		RegionName:               "Zurich",
		RegionReplicationEnabled: false,
		S3Endpoints:              []string{"https://zh1-s3.xelon.io"},
		Tenant:                   &Tenant{ID: "111111111", Name: "test-tenant-1"},
		VersioningEnabled:        false,
	}}

	actualBuckets, resp, err := client.ObjectStorages.ListBuckets(ctx, &ObjectStorageBucketListOptions{
		TenantID: "000000000",
		ListOptions: ListOptions{
			Page:    2,
			PerPage: 5,
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, &Meta{
		Total:    43,
		LastPage: 9,
		PerPage:  5,
		Page:     1,
		From:     1,
		To:       5,
	}, resp.Meta)
	assert.Equal(t, expectedBuckets, actualBuckets)
}

func TestObjectStorages_AllBuckets(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("GET /object-storages/buckets", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "1", r.URL.Query().Get("perPage"))

		switch r.URL.Query().Get("page") {
		case "1":
			fixture := loadFixture(t, "objectstorages_all_buckets_page_1.json")
			_, _ = w.Write(fixture)
		case "2":
			fixture := loadFixture(t, "objectstorages_all_buckets_page_2.json")
			_, _ = w.Write(fixture)
		default:
			t.Fatalf("unexpected page %q", r.URL.Query().Get("page"))
		}
	})

	seq, errFn := client.ObjectStorages.AllBuckets(ctx, &ListOptions{PerPage: 1})

	var actualBuckets []ObjectStorageBucket
	for bucket := range seq {
		actualBuckets = append(actualBuckets, bucket)
	}

	assert.NoError(t, errFn())
	assert.Equal(t, []ObjectStorageBucket{{
		ID:   "zone1.4711.1",
		Name: "test-bucket-0",
	}, {
		ID:   "zone2.4711.2",
		Name: "test-bucket-1",
	}}, actualBuckets)
}

func TestObjectStorages_GetBucket(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("GET /object-storages/buckets/test-bucket-0/00000000-0000-0000-0000-000000000000", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fixture := loadFixture(t, "objectstorages_get_bucket_success.json")
		_, _ = w.Write(fixture)
	})
	expectedBucket := &ObjectStorageBucket{
		CreatedAt:                mustTime(t, "2025-10-27T14:19:56+01:00"),
		ID:                       "zone1.4711.1",
		IPRestrictionsEnabled:    true,
		Name:                     "test-bucket-0",
		ObjectLockEnabled:        false,
		ObjectLockRetentionDays:  0,
		ObjectStorageUserID:      "00000000-0000-0000-0000-000000000000",
		ObjectStorageUserName:    "test-user-0",
		RegionName:               "Aargau",
		RegionReplicationEnabled: true,
		S3Endpoints:              []string{"https://ch1-s3.xelon.io"},
		Tenant:                   &Tenant{ID: "000000000", Name: "test-tenant-0"},
		VersioningEnabled:        true,
	}

	actualBucket, resp, err := client.ObjectStorages.GetBucket(ctx, "test-bucket-0", "00000000-0000-0000-0000-000000000000")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedBucket, actualBucket)
}

func TestObjectStorages_GetBucket_MissingData(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("GET /object-storages/buckets/test-bucket-0/00000000-0000-0000-0000-000000000000", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		_, _ = w.Write([]byte(`{}`))
	})

	actualBucket, resp, err := client.ObjectStorages.GetBucket(ctx, "test-bucket-0", "00000000-0000-0000-0000-000000000000")

	assert.Nil(t, actualBucket)
	assert.NotNil(t, resp)
	assert.EqualError(t, err, "object storage bucket data is empty")
}

func TestObjectStorages_CreateBucket(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("POST /object-storages/buckets", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		actualRequest, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.JSONEq(t, `{
			"isObjectLock": false,
			"isVersioning": true,
			"name": "test-bucket-0",
			"retentionPeriodDays": 30,
			"s3UserIdentifier": "00000000-0000-0000-0000-000000000000"
		}`, string(actualRequest))

		fixture := loadFixture(t, "objectstorages_create_bucket_success.json")
		_, _ = w.Write(fixture)
	})
	expectedBucket := &ObjectStorageBucket{
		CreatedAt:                mustTime(t, "2025-10-27T14:19:56+01:00"),
		ID:                       "zone1.4711.1",
		IPRestrictionsEnabled:    false,
		Name:                     "test-bucket-0",
		ObjectLockEnabled:        false,
		ObjectLockRetentionDays:  30,
		ObjectStorageUserID:      "00000000-0000-0000-0000-000000000000",
		ObjectStorageUserName:    "test-user-0",
		RegionName:               "Aargau",
		RegionReplicationEnabled: true,
		S3Endpoints:              []string{"https://ch1-s3.xelon.io"},
		Tenant:                   &Tenant{ID: "000000000", Name: "test-tenant-0"},
		VersioningEnabled:        true,
	}

	actualBucket, resp, err := client.ObjectStorages.CreateBucket(ctx, &ObjectStorageBucketCreateRequest{
		Name:                    "test-bucket-0",
		ObjectLockEnabled:       false,
		ObjectLockRetentionDays: 30,
		ObjectStorageUserID:     "00000000-0000-0000-0000-000000000000",
		VersioningEnabled:       true,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedBucket, actualBucket)
}

func TestObjectStorages_CreateBucket_ZeroObjectLockRetentionDaysOmitted(t *testing.T) {
	actualRequest, err := json.Marshal(ObjectStorageBucketCreateRequest{
		Name:                    "test-bucket-0",
		ObjectLockEnabled:       true,
		ObjectLockRetentionDays: 0,
		ObjectStorageUserID:     "00000000-0000-0000-0000-000000000000",
		VersioningEnabled:       true,
	})
	assert.NoError(t, err)

	var payload map[string]any
	err = json.Unmarshal(actualRequest, &payload)
	assert.NoError(t, err)
	assert.NotContains(t, payload, "retentionPeriodDays")
}

func TestObjectStorages_CreateBucket_MissingData(t *testing.T) {
	setup()
	defer teardown()

	type testCase struct {
		responseBody string
	}
	tests := map[string]testCase{
		"missing data": {
			responseBody: `{"message":"Bucket successfully created"}`,
		},
		"null data": {
			responseBody: `{"data":null,"message":"Bucket successfully created"}`,
		},
	}

	var responseBody string
	mux.HandleFunc("POST /object-storages/buckets", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		_, _ = w.Write([]byte(responseBody))
	})

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			responseBody = test.responseBody

			actualBucket, resp, err := client.ObjectStorages.CreateBucket(ctx, &ObjectStorageBucketCreateRequest{
				Name:                "test-bucket-0",
				ObjectStorageUserID: "00000000-0000-0000-0000-000000000000",
				VersioningEnabled:   true,
			})

			assert.Nil(t, actualBucket)
			assert.NotNil(t, resp)
			assert.EqualError(t, err, "object storage bucket data is empty")
		})
	}
}

func TestObjectStorages_UpdateBucket(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("PUT /object-storages/buckets/zone1.4711.1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)

		var actualRequest ObjectStorageBucketUpdateRequest
		err := json.NewDecoder(r.Body).Decode(&actualRequest)
		assert.NoError(t, err)
		assert.Equal(t, ObjectStorageBucketUpdateRequest{Name: "test-bucket-0-updated"}, actualRequest)

		_, _ = w.Write([]byte(`{"message":"Bucket successfully edited"}`))
	})

	resp, err := client.ObjectStorages.UpdateBucket(ctx, "zone1.4711.1", &ObjectStorageBucketUpdateRequest{Name: "test-bucket-0-updated"})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestObjectStorages_DeleteBucket(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("DELETE /object-storages/buckets/test-bucket-0/00000000-0000-0000-0000-000000000000", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		_, _ = w.Write([]byte(`{"message":"Bucket successfully deleted"}`))
	})

	resp, err := client.ObjectStorages.DeleteBucket(ctx, "test-bucket-0", "00000000-0000-0000-0000-000000000000")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestObjectStorages_UpdateBucketVersioning(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("PUT /object-storages/buckets/test-bucket-0/00000000-0000-0000-0000-000000000000/versioning", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)

		var actualRequest ObjectStorageBucketVersioningUpdateRequest
		err := json.NewDecoder(r.Body).Decode(&actualRequest)
		assert.NoError(t, err)
		assert.Equal(t, ObjectStorageBucketVersioningUpdateRequest{VersioningEnabled: false}, actualRequest)

		_, _ = w.Write([]byte(`{"message":"Bucket versioning has been disabled"}`))
	})

	resp, err := client.ObjectStorages.UpdateBucketVersioning(ctx, "test-bucket-0", "00000000-0000-0000-0000-000000000000", &ObjectStorageBucketVersioningUpdateRequest{
		VersioningEnabled: false,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestObjectStorages_GetBucketIPRestrictions(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("GET /object-storages/buckets/test-bucket-0/00000000-0000-0000-0000-000000000000/ip-restrictions", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fixture := loadFixture(t, "objectstorages_get_bucket_ip_restrictions_success.json")
		_, _ = w.Write(fixture)
	})
	expectedRestrictions := &ObjectStorageBucketIPRestrictions{
		AllowedIPs: []string{"192.168.1.0/24", "10.0.0.1"},
		Enabled:    true,
	}

	actualRestrictions, resp, err := client.ObjectStorages.GetBucketIPRestrictions(ctx, "test-bucket-0", "00000000-0000-0000-0000-000000000000")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedRestrictions, actualRestrictions)
}

func TestObjectStorages_GetBucketIPRestrictions_MissingData(t *testing.T) {
	setup()
	defer teardown()

	type testCase struct {
		responseBody string
	}
	tests := map[string]testCase{
		"missing data": {
			responseBody: `{"message":"IP restrictions fetched"}`,
		},
		"null data": {
			responseBody: `{"data":null}`,
		},
	}

	var responseBody string
	mux.HandleFunc("GET /object-storages/buckets/test-bucket-0/00000000-0000-0000-0000-000000000000/ip-restrictions", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		_, _ = w.Write([]byte(responseBody))
	})

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			responseBody = test.responseBody

			actualRestrictions, resp, err := client.ObjectStorages.GetBucketIPRestrictions(ctx, "test-bucket-0", "00000000-0000-0000-0000-000000000000")

			assert.Nil(t, actualRestrictions)
			assert.NotNil(t, resp)
			assert.EqualError(t, err, "object storage bucket ip restrictions data is empty")
		})
	}
}

func TestObjectStorages_UpdateBucketIPRestrictions(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("PUT /object-storages/buckets/test-bucket-0/00000000-0000-0000-0000-000000000000/ip-restrictions", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)

		var actualRequest ObjectStorageBucketIPRestrictionsUpdateRequest
		err := json.NewDecoder(r.Body).Decode(&actualRequest)
		assert.NoError(t, err)
		assert.Equal(t, ObjectStorageBucketIPRestrictionsUpdateRequest{
			AllowedIPs: []string{"192.168.1.0/24", "10.0.0.1"},
			Enabled:    true,
		}, actualRequest)

		_, _ = w.Write([]byte(`{"message":"IP restrictions will updated soon"}`))
	})

	resp, err := client.ObjectStorages.UpdateBucketIPRestrictions(ctx, "test-bucket-0", "00000000-0000-0000-0000-000000000000", &ObjectStorageBucketIPRestrictionsUpdateRequest{
		AllowedIPs: []string{"192.168.1.0/24", "10.0.0.1"},
		Enabled:    true,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestObjectStorages_UpdateBucketIPRestrictions_DisableSendsEnabledFalse(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("PUT /object-storages/buckets/test-bucket-0/00000000-0000-0000-0000-000000000000/ip-restrictions", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)

		var actualRequest map[string]any
		err := json.NewDecoder(r.Body).Decode(&actualRequest)
		assert.NoError(t, err)
		assert.Equal(t, map[string]any{"enabled": false}, actualRequest)

		_, _ = w.Write([]byte(`{"message":"IP restrictions will updated soon"}`))
	})

	resp, err := client.ObjectStorages.UpdateBucketIPRestrictions(ctx, "test-bucket-0", "00000000-0000-0000-0000-000000000000", &ObjectStorageBucketIPRestrictionsUpdateRequest{
		Enabled: false,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestObjectStorages_BucketValidation(t *testing.T) {
	tests := map[string]struct {
		err    error
		target error
	}{
		"get bucket missing bucket name": {
			err:    errorFromBucketGet(client.ObjectStorages.GetBucket(ctx, "", "00000000-0000-0000-0000-000000000000")),
			target: ErrEmptyArgument,
		},
		"get bucket missing user id": {
			err:    errorFromBucketGet(client.ObjectStorages.GetBucket(ctx, "test-bucket-0", "")),
			target: ErrEmptyArgument,
		},
		"create bucket missing payload": {
			err:    errorFromBucketCreate(client.ObjectStorages.CreateBucket(ctx, nil)),
			target: ErrEmptyPayloadNotAllowed,
		},
		"update bucket missing bucket id": {
			err:    errorFromResponse(client.ObjectStorages.UpdateBucket(ctx, "", &ObjectStorageBucketUpdateRequest{Name: "test-bucket-0"})),
			target: ErrEmptyArgument,
		},
		"update bucket missing payload": {
			err:    errorFromResponse(client.ObjectStorages.UpdateBucket(ctx, "zone1.4711.1", nil)),
			target: ErrEmptyPayloadNotAllowed,
		},
		"delete bucket missing bucket name": {
			err:    errorFromResponse(client.ObjectStorages.DeleteBucket(ctx, "", "00000000-0000-0000-0000-000000000000")),
			target: ErrEmptyArgument,
		},
		"delete bucket missing user id": {
			err:    errorFromResponse(client.ObjectStorages.DeleteBucket(ctx, "test-bucket-0", "")),
			target: ErrEmptyArgument,
		},
		"update versioning missing bucket name": {
			err: errorFromResponse(client.ObjectStorages.UpdateBucketVersioning(ctx, "", "00000000-0000-0000-0000-000000000000", &ObjectStorageBucketVersioningUpdateRequest{
				VersioningEnabled: true,
			})),
			target: ErrEmptyArgument,
		},
		"update versioning missing user id": {
			err: errorFromResponse(client.ObjectStorages.UpdateBucketVersioning(ctx, "test-bucket-0", "", &ObjectStorageBucketVersioningUpdateRequest{
				VersioningEnabled: true,
			})),
			target: ErrEmptyArgument,
		},
		"update versioning missing payload": {
			err:    errorFromResponse(client.ObjectStorages.UpdateBucketVersioning(ctx, "test-bucket-0", "00000000-0000-0000-0000-000000000000", nil)),
			target: ErrEmptyPayloadNotAllowed,
		},
		"get ip restrictions missing bucket name": {
			err:    errorFromIPRestrictionsGet(client.ObjectStorages.GetBucketIPRestrictions(ctx, "", "00000000-0000-0000-0000-000000000000")),
			target: ErrEmptyArgument,
		},
		"get ip restrictions missing user id": {
			err:    errorFromIPRestrictionsGet(client.ObjectStorages.GetBucketIPRestrictions(ctx, "test-bucket-0", "")),
			target: ErrEmptyArgument,
		},
		"update ip restrictions missing bucket name": {
			err:    errorFromResponse(client.ObjectStorages.UpdateBucketIPRestrictions(ctx, "", "00000000-0000-0000-0000-000000000000", &ObjectStorageBucketIPRestrictionsUpdateRequest{Enabled: true})),
			target: ErrEmptyArgument,
		},
		"update ip restrictions missing user id": {
			err:    errorFromResponse(client.ObjectStorages.UpdateBucketIPRestrictions(ctx, "test-bucket-0", "", &ObjectStorageBucketIPRestrictionsUpdateRequest{Enabled: true})),
			target: ErrEmptyArgument,
		},
		"update ip restrictions missing payload": {
			err:    errorFromResponse(client.ObjectStorages.UpdateBucketIPRestrictions(ctx, "test-bucket-0", "00000000-0000-0000-0000-000000000000", nil)),
			target: ErrEmptyPayloadNotAllowed,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.True(t, errors.Is(test.err, test.target))
		})
	}
}

func errorFromBucketGet(_ *ObjectStorageBucket, _ *Response, err error) error {
	return err
}

func errorFromBucketCreate(_ *ObjectStorageBucket, _ *Response, err error) error {
	return err
}

func errorFromIPRestrictionsGet(_ *ObjectStorageBucketIPRestrictions, _ *Response, err error) error {
	return err
}

func errorFromResponse(_ *Response, err error) error {
	return err
}
