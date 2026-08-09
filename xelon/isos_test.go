package xelon

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestISOs_List(t *testing.T) {
	setup(t)
	defer teardown()

	mux.HandleFunc("/isos", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fixture := loadFixture(t, "isos_list.json")
		_, _ = w.Write(fixture)
	})
	expectedISOs := []ISO{
		{ID: "abc123", Name: "iso-1", Status: true},
		{ID: "def456", Name: "iso-2", Status: false},
	}
	expectedMeta := &Meta{Total: 20, LastPage: 2, PerPage: 10, Page: 1, From: 1, To: 10}

	actualISOs, resp, err := client.ISOs.List(ctx, nil)

	assert.NoError(t, err)
	assert.Equal(t, expectedISOs, actualISOs)
	assert.Equal(t, expectedMeta, resp.Meta)
}

func TestISOs_Upload(t *testing.T) {
	setup(t)
	defer teardown()

	fileBytes := []byte{0x00, 0x01, 0x02, 0xff, 'I', 'S', 'O'}
	mux.HandleFunc("POST /isos/upload", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/isos/upload", r.URL.Path)
		assert.Positive(t, r.ContentLength)
		assert.Empty(t, r.TransferEncoding)

		mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		assert.NoError(t, err)
		assert.Equal(t, "multipart/form-data", mediaType)
		assert.NotEmpty(t, params["boundary"])
		assert.NoError(t, r.ParseMultipartForm(1<<20))
		assert.Equal(t, []string{"test-iso"}, r.MultipartForm.Value["name"])
		assert.Equal(t, []string{"7"}, r.MultipartForm.Value["categoryId"])
		assert.Equal(t, []string{"cloud-123"}, r.MultipartForm.Value["cloudIdentifier"])
		assert.Equal(t, []string{"small test ISO"}, r.MultipartForm.Value["description"])
		assert.Equal(t, []string{"tenant-123"}, r.MultipartForm.Value["tenantIdentifier"])
		_, hasFilenameField := r.MultipartForm.Value["filename"]
		assert.False(t, hasFilenameField)

		files := r.MultipartForm.File["file"]
		if assert.Len(t, files, 1) {
			fileHeader := files[0]
			assert.Equal(t, "source.iso", fileHeader.Filename)
			assert.Equal(t, "application/octet-stream", fileHeader.Header.Get("Content-Type"))
			assert.Equal(t, `form-data; name="file"; filename="source.iso"`, fileHeader.Header.Get("Content-Disposition"))

			file, err := fileHeader.Open()
			assert.NoError(t, err)
			if err == nil {
				defer func() { _ = file.Close() }()
				actualFile, err := io.ReadAll(file)
				assert.NoError(t, err)
				assert.Equal(t, fileBytes, actualFile)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		fixture := loadFixture(t, "isos_upload_success.json")
		_, _ = w.Write(fixture)
	})

	iso, resp, err := client.ISOs.Upload(ctx, &ISOUploadRequest{
		CategoryID:  7,
		CloudID:     "cloud-123",
		Description: "small test ISO",
		File:        bytes.NewReader(fileBytes),
		Filename:    "/tmp/nested/source.iso",
		Name:        "test-iso",
		TenantID:    "tenant-123",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, &ISO{
		Active:      true,
		Category:    "Windows Server 2012/2016/2019",
		Cloud:       &Cloud{ID: "cloud-123", Name: "Example Cloud", Type: "public"},
		CreatedAt:   mustTime(t, "2026-07-27T17:58:24+02:00"),
		Description: "",
		ID:          "iso-123",
		Name:        "test-iso",
		Owner:       "Example Tenant",
		Status:      true,
	}, iso)
}

func TestISOs_UploadOmitsEmptyOptionalFields(t *testing.T) {
	setup(t)
	defer teardown()

	mux.HandleFunc("POST /isos/upload", func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(t, r.ParseMultipartForm(1<<20))
		_, hasDescription := r.MultipartForm.Value["description"]
		_, hasTenantID := r.MultipartForm.Value["tenantIdentifier"]
		assert.False(t, hasDescription)
		assert.False(t, hasTenantID)
		_, _ = fmt.Fprint(w, `{"data":{"identifier":"iso-123"}}`)
	})

	iso, _, err := client.ISOs.Upload(ctx, validISOUploadRequest())

	assert.NoError(t, err)
	assert.Equal(t, "iso-123", iso.ID)
}

func TestISOs_UploadValidation(t *testing.T) {
	requests := 0
	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		requests++
		return newTestResponse(req, http.StatusOK, `{"data":{"identifier":"unexpected"}}`), nil
	})}
	client := NewClient("auth-token", WithHTTPClient(httpClient))

	tests := map[string]struct {
		request *ISOUploadRequest
		target  error
	}{
		"nil payload": {
			request: nil,
			target:  ErrEmptyPayloadNotAllowed,
		},
		"nil file": {
			request: func() *ISOUploadRequest {
				request := validISOUploadRequest()
				request.File = nil
				return request
			}(),
			target: ErrEmptyArgument,
		},
		"missing filename": {
			request: func() *ISOUploadRequest {
				request := validISOUploadRequest()
				request.Filename = ""
				return request
			}(),
			target: ErrEmptyArgument,
		},
		"missing name": {
			request: func() *ISOUploadRequest {
				request := validISOUploadRequest()
				request.Name = ""
				return request
			}(),
			target: ErrEmptyArgument,
		},
		"missing cloud id": {
			request: func() *ISOUploadRequest {
				request := validISOUploadRequest()
				request.CloudID = ""
				return request
			}(),
			target: ErrEmptyArgument,
		},
		"missing category id": {
			request: func() *ISOUploadRequest {
				request := validISOUploadRequest()
				request.CategoryID = 0
				return request
			}(),
			target: ErrEmptyArgument,
		},
		"negative category id": {
			request: func() *ISOUploadRequest {
				request := validISOUploadRequest()
				request.CategoryID = -1
				return request
			}(),
			target: ErrEmptyArgument,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, _, err := client.ISOs.Upload(context.Background(), test.request)
			assert.Error(t, err)
			assert.ErrorIs(t, err, test.target)
		})
	}
	assert.Zero(t, requests)
}

func TestISOs_UploadPreparationErrorsDoNotSendRequest(t *testing.T) {
	readerFailure := errors.New("synthetic ISO read failure")
	requests := 0
	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		requests++
		return newTestResponse(req, http.StatusOK, `{"data":{"identifier":"unexpected"}}`), nil
	})}
	client := NewClient("auth-token", WithHTTPClient(httpClient))

	request := validISOUploadRequest()
	request.File = failingReader{err: readerFailure}
	_, _, err := client.ISOs.Upload(context.Background(), request)
	assert.ErrorIs(t, err, readerFailure)
	assert.Contains(t, err.Error(), "copy ISO file")

	client.baseURL, _ = url.Parse("https://example.test/api/v2")
	_, _, err = client.ISOs.Upload(context.Background(), validISOUploadRequest())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "create upload request")
	assert.Zero(t, requests)
}

func TestISOs_UploadMultipartWriteError(t *testing.T) {
	writeFailure := errors.New("synthetic multipart write failure")
	writer := multipart.NewWriter(failingWriter{err: writeFailure})

	err := writeISOUploadMultipart(context.Background(), writer, validISOUploadRequest())

	assert.ErrorIs(t, err, writeFailure)
	assert.Contains(t, err.Error(), "write multipart field")
}

func TestISOs_UploadResponseErrors(t *testing.T) {
	tests := map[string]struct {
		body       string
		statusCode int
		validate   func(*testing.T, *Response, error)
	}{
		"non-2xx": {
			body:       `{"message":"invalid upload"}`,
			statusCode: http.StatusUnprocessableEntity,
			validate: func(t *testing.T, resp *Response, err error) {
				var responseError *ErrorResponse
				assert.ErrorAs(t, err, &responseError)
				assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
				assert.Contains(t, err.Error(), "upload ISO")
			},
		},
		"decode failure": {
			body:       `{"data":`,
			statusCode: http.StatusOK,
			validate: func(t *testing.T, resp *Response, err error) {
				assert.NotNil(t, resp)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "upload ISO")
			},
		},
		"missing data": {
			body:       `{"message":"ISO uploaded successfully."}`,
			statusCode: http.StatusOK,
			validate: func(t *testing.T, resp *Response, err error) {
				assert.NotNil(t, resp)
				assert.EqualError(t, err, "iso data is empty")
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			setup(t)
			defer teardown()
			mux.HandleFunc("POST /isos/upload", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(test.statusCode)
				_, _ = io.WriteString(w, test.body)
			})

			iso, resp, err := client.ISOs.Upload(context.Background(), validISOUploadRequest())
			assert.Nil(t, iso)
			test.validate(t, resp, err)
		})
	}
}

func TestISOs_UploadTimeoutPolicy(t *testing.T) {
	t.Run("SDK-owned client uses upload fallback", func(t *testing.T) {
		client := NewClient("auth-token")
		client.httpClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
			deadline, ok := req.Context().Deadline()
			assert.True(t, ok)
			assert.InDelta(t, defaultISOUploadTimeout.Seconds(), time.Until(deadline).Seconds(), 0.5)
			return newTestResponse(req, http.StatusOK, `{"data":{"identifier":"iso-123"}}`), nil
		})

		iso, _, err := client.ISOs.Upload(context.Background(), validISOUploadRequest())
		assert.NoError(t, err)
		assert.Equal(t, "iso-123", iso.ID)
	})

	t.Run("caller deadline longer than fallback is preserved", func(t *testing.T) {
		client := NewClient("auth-token")
		deadline := time.Now().Add(5 * time.Minute)
		ctx, cancel := context.WithDeadline(context.Background(), deadline)
		defer cancel()
		client.httpClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
			got, ok := req.Context().Deadline()
			assert.True(t, ok)
			assert.True(t, got.Equal(deadline), "deadline = %v, want %v", got, deadline)
			return newTestResponse(req, http.StatusOK, `{"data":{"identifier":"iso-123"}}`), nil
		})

		_, _, err := client.ISOs.Upload(ctx, validISOUploadRequest())
		assert.NoError(t, err)
	})

	t.Run("custom zero-timeout client suppresses upload fallback", func(t *testing.T) {
		httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			_, ok := req.Context().Deadline()
			assert.False(t, ok)
			return newTestResponse(req, http.StatusOK, `{"data":{"identifier":"iso-123"}}`), nil
		})}
		client := NewClient("auth-token", WithHTTPClient(httpClient))

		_, _, err := client.ISOs.Upload(context.Background(), validISOUploadRequest())
		assert.NoError(t, err)
	})

	t.Run("already-cancelled context skips preparation and request", func(t *testing.T) {
		requests := 0
		client := NewClient("auth-token")
		client.httpClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
			requests++
			return newTestResponse(req, http.StatusOK, `{"data":{"identifier":"unexpected"}}`), nil
		})
		reader := new(countingReader)
		uploadRequest := validISOUploadRequest()
		uploadRequest.File = reader
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		started := time.Now()
		_, _, err := client.ISOs.Upload(ctx, uploadRequest)
		assert.ErrorIs(t, err, context.Canceled)
		assert.Less(t, time.Since(started), time.Second)
		assert.Zero(t, reader.reads)
		assert.Zero(t, requests)
	})

	t.Run("cancellation during buffering stops before request", func(t *testing.T) {
		requests := 0
		client := NewClient("auth-token")
		client.httpClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
			requests++
			return newTestResponse(req, http.StatusOK, `{"data":{"identifier":"unexpected"}}`), nil
		})
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		reader := &cancellingReader{cancel: cancel}
		uploadRequest := validISOUploadRequest()
		uploadRequest.File = reader

		_, _, err := client.ISOs.Upload(ctx, uploadRequest)
		assert.ErrorIs(t, err, context.Canceled)
		assert.Contains(t, err.Error(), "copy ISO file")
		assert.Equal(t, 1, reader.reads)
		assert.Zero(t, requests)
	})
}

func validISOUploadRequest() *ISOUploadRequest {
	return &ISOUploadRequest{
		CategoryID: 1,
		CloudID:    "cloud-123",
		File:       strings.NewReader("ISO"),
		Filename:   "source.iso",
		Name:       "test-iso",
	}
}

type failingReader struct {
	err error
}

func (r failingReader) Read([]byte) (int, error) {
	return 0, r.err
}

type failingWriter struct {
	err error
}

func (w failingWriter) Write([]byte) (int, error) {
	return 0, w.err
}

type countingReader struct {
	reads int
}

func (r *countingReader) Read([]byte) (int, error) {
	r.reads++
	return 0, io.EOF
}

type cancellingReader struct {
	cancel context.CancelFunc
	reads  int
}

func (r *cancellingReader) Read(p []byte) (int, error) {
	r.reads++
	n := copy(p, "partial ISO")
	r.cancel()
	return n, nil
}
