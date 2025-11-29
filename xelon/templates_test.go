package xelon

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTemplatesService_List(t *testing.T) {
	// Mock API response with v2 structure
	mockResponse := `{
		"data": [
			{
				"identifier": "22a6711de3a8",
				"name": "AlmaLinux 8  64 Bit EN",
				"cloudIdentifier": "0984199ded4e",
				"type": "Linux",
				"description": "AlmaLinux 8 English 64 Bit",
				"category": "Linux",
				"status": 2,
				"active": true,
				"templateType": "linux",
				"createdAt": "2021-08-10T10:00:00Z",
				"updatedAt": "2021-08-10T10:00:00Z"
			},
			{
				"identifier": "33b7822ef4b9",
				"name": "Debian 11  64 Bit EN",
				"cloudIdentifier": "0984199ded4e",
				"type": "Linux",
				"description": "Debian 11 English 64 Bit",
				"category": "Linux",
				"status": 2,
				"active": true,
				"templateType": "linux",
				"createdAt": "2021-09-15T10:00:00Z",
				"updatedAt": "2021-09-15T10:00:00Z"
			}
		],
		"meta": {
			"total": 2,
			"lastPage": 1,
			"perPage": 10,
			"currentPage": 1,
			"from": 1,
			"to": 2
		}
	}`

	// Create test server
	mux := http.NewServeMux()
	mux.HandleFunc("/templates", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockResponse))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	// Create client with test server URL
	client := NewClient("test-token", WithBaseURL(server.URL+"/"))

	// Test List method with nil options (uses defaults)
	templates, _, err := client.Templates.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("Templates.List returned error: %v", err)
	}

	// Verify we got 2 templates
	if len(templates) != 2 {
		t.Errorf("Expected 2 templates, got %d", len(templates))
	}

	// Verify first template
	if templates[0].Identifier != "22a6711de3a8" {
		t.Errorf("Expected identifier '22a6711de3a8', got '%s'", templates[0].Identifier)
	}
	if templates[0].Name != "AlmaLinux 8  64 Bit EN" {
		t.Errorf("Expected name 'AlmaLinux 8  64 Bit EN', got '%s'", templates[0].Name)
	}
	if templates[0].Type != "Linux" {
		t.Errorf("Expected type 'Linux', got '%s'", templates[0].Type)
	}
	if templates[0].CloudIdentifier != "0984199ded4e" {
		t.Errorf("Expected cloudIdentifier '0984199ded4e', got '%s'", templates[0].CloudIdentifier)
	}

	// Verify second template
	if templates[1].Identifier != "33b7822ef4b9" {
		t.Errorf("Expected identifier '33b7822ef4b9', got '%s'", templates[1].Identifier)
	}
	if templates[1].Name != "Debian 11  64 Bit EN" {
		t.Errorf("Expected name 'Debian 11  64 Bit EN', got '%s'", templates[1].Name)
	}
}

func TestTemplatesService_List_Pagination(t *testing.T) {
	// Create test server that serves multiple pages
	mux := http.NewServeMux()
	mux.HandleFunc("/templates", func(w http.ResponseWriter, r *http.Request) {
		pageParam := r.URL.Query().Get("page")

		var response string
		if pageParam == "1" || pageParam == "" {
			response = `{
				"data": [
					{"identifier": "id1", "name": "Template 1", "cloudIdentifier": "cloud1", "type": "Linux", "status": 2, "active": true, "templateType": "linux", "createdAt": "2021-01-01T00:00:00Z", "updatedAt": "2021-01-01T00:00:00Z"}
				],
				"meta": {"total": 2, "lastPage": 2, "perPage": 1, "currentPage": 1, "from": 1, "to": 1}
			}`
		} else if pageParam == "2" {
			response = `{
				"data": [
					{"identifier": "id2", "name": "Template 2", "cloudIdentifier": "cloud1", "type": "Windows", "status": 2, "active": true, "templateType": "windows", "createdAt": "2021-01-02T00:00:00Z", "updatedAt": "2021-01-02T00:00:00Z"}
				],
				"meta": {"total": 2, "lastPage": 2, "perPage": 1, "currentPage": 2, "from": 2, "to": 2}
			}`
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient("test-token", WithBaseURL(server.URL+"/"))

	// Test pagination - should fetch both pages
	templates, _, err := client.Templates.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("Templates.List returned error: %v", err)
	}

	// Should have collected templates from both pages
	if len(templates) != 2 {
		t.Errorf("Expected 2 templates from pagination, got %d", len(templates))
	}

	// Verify we got both templates
	if templates[0].Identifier != "id1" {
		t.Errorf("Expected first template identifier 'id1', got '%s'", templates[0].Identifier)
	}
	if templates[1].Identifier != "id2" {
		t.Errorf("Expected second template identifier 'id2', got '%s'", templates[1].Identifier)
	}
}

func TestTemplatesService_Get(t *testing.T) {
	mockResponse := `{
		"data": [
			{
				"identifier": "22a6711de3a8",
				"name": "AlmaLinux 8  64 Bit EN",
				"cloudIdentifier": "0984199ded4e",
				"type": "Linux",
				"description": "AlmaLinux 8 English 64 Bit",
				"category": "Linux",
				"status": 2,
				"active": true,
				"templateType": "linux",
				"createdAt": "2021-08-10T10:00:00Z",
				"updatedAt": "2021-08-10T10:00:00Z"
			}
		],
		"meta": {"total": 1, "lastPage": 1, "perPage": 10, "currentPage": 1, "from": 1, "to": 1}
	}`

	mux := http.NewServeMux()
	mux.HandleFunc("/templates", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockResponse))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient("test-token", WithBaseURL(server.URL+"/"))

	// Test Get method
	template, _, err := client.Templates.Get(context.Background(), "22a6711de3a8")
	if err != nil {
		t.Fatalf("Templates.Get returned error: %v", err)
	}

	if template == nil {
		t.Fatal("Expected template, got nil")
	}

	if template.Identifier != "22a6711de3a8" {
		t.Errorf("Expected identifier '22a6711de3a8', got '%s'", template.Identifier)
	}
	if template.Name != "AlmaLinux 8  64 Bit EN" {
		t.Errorf("Expected name 'AlmaLinux 8  64 Bit EN', got '%s'", template.Name)
	}
}

func TestTemplatesService_List_WithFilters(t *testing.T) {
	// Create test server that validates query parameters
	mux := http.NewServeMux()
	mux.HandleFunc("/templates", func(w http.ResponseWriter, r *http.Request) {
		// Verify query parameters are set correctly
		query := r.URL.Query()

		if query.Get("cloudIdentifier") != "0984199ded4e" {
			t.Errorf("Expected cloudIdentifier=0984199ded4e, got %s", query.Get("cloudIdentifier"))
		}
		if query.Get("type") != "Linux" {
			t.Errorf("Expected type=Linux, got %s", query.Get("type"))
		}
		if query.Get("perPage") != "100" {
			t.Errorf("Expected perPage=100, got %s", query.Get("perPage"))
		}
		if query.Get("search") != "ubuntu" {
			t.Errorf("Expected search=ubuntu, got %s", query.Get("search"))
		}

		response := `{
			"data": [{
				"identifier": "id1",
				"name": "Ubuntu 22",
				"cloudIdentifier": "0984199ded4e",
				"type": "Linux",
				"status": 2,
				"active": true,
				"templateType": "linux",
				"createdAt": "2021-01-01T00:00:00Z",
				"updatedAt": "2021-01-01T00:00:00Z"
			}],
			"meta": {"total": 1, "lastPage": 1, "perPage": 100, "currentPage": 1, "from": 1, "to": 1}
		}`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient("test-token", WithBaseURL(server.URL+"/"))

	// Test List with filters
	opts := &TemplateListOptions{
		CloudIdentifier: "0984199ded4e",
		Type:            "Linux",
		Search:          "ubuntu",
		PerPage:         100,
	}

	templates, _, err := client.Templates.List(context.Background(), opts)
	if err != nil {
		t.Fatalf("Templates.List returned error: %v", err)
	}

	if len(templates) != 1 {
		t.Errorf("Expected 1 template, got %d", len(templates))
	}

	if templates[0].Name != "Ubuntu 22" {
		t.Errorf("Expected name 'Ubuntu 22', got '%s'", templates[0].Name)
	}
}

func TestTemplatesService_GetByName(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/templates", func(w http.ResponseWriter, r *http.Request) {
		// Verify search parameter is used
		query := r.URL.Query()
		if query.Get("search") != "ubuntu 22" {
			t.Errorf("Expected search parameter, got %s", query.Get("search"))
		}

		response := `{
			"data": [{
				"identifier": "0b8e53f7147a",
				"name": "ubuntu 22",
				"cloudIdentifier": "0984199ded4e",
				"type": "Linux",
				"status": 2,
				"active": true,
				"templateType": "linux",
				"createdAt": "2021-01-01T00:00:00Z",
				"updatedAt": "2021-01-01T00:00:00Z"
			}],
			"meta": {"total": 1, "lastPage": 1, "perPage": 100, "currentPage": 1, "from": 1, "to": 1}
		}`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient("test-token", WithBaseURL(server.URL+"/"))

	template, _, err := client.Templates.GetByName(context.Background(), "ubuntu 22", nil)
	if err != nil {
		t.Fatalf("Templates.GetByName returned error: %v", err)
	}

	if template == nil {
		t.Fatal("Expected template, got nil")
	}

	if template.Name != "ubuntu 22" {
		t.Errorf("Expected name 'ubuntu 22', got '%s'", template.Name)
	}
}

func TestTemplateV2_JSONUnmarshal(t *testing.T) {
	jsonData := `{
		"identifier": "22a6711de3a8",
		"name": "AlmaLinux 8  64 Bit EN",
		"ownerTenant": null,
		"cloudIdentifier": "0984199ded4e",
		"status": 2,
		"description": "AlmaLinux 8 English 64 Bit",
		"internalVersion": "1.0",
		"active": true,
		"templateType": "linux",
		"scriptType": null,
		"type": "Linux",
		"category": "Linux",
		"createdAt": "2021-08-10T10:00:00Z",
		"updatedAt": "2021-08-10T10:00:00Z",
		"deletedAt": null
	}`

	var template TemplateV2
	err := json.Unmarshal([]byte(jsonData), &template)
	if err != nil {
		t.Fatalf("Failed to unmarshal TemplateV2: %v", err)
	}

	if template.Identifier != "22a6711de3a8" {
		t.Errorf("Expected identifier '22a6711de3a8', got '%s'", template.Identifier)
	}
	if template.Name != "AlmaLinux 8  64 Bit EN" {
		t.Errorf("Expected name 'AlmaLinux 8  64 Bit EN', got '%s'", template.Name)
	}
	if template.CloudIdentifier != "0984199ded4e" {
		t.Errorf("Expected cloudIdentifier '0984199ded4e', got '%s'", template.CloudIdentifier)
	}
	if template.Type != "Linux" {
		t.Errorf("Expected type 'Linux', got '%s'", template.Type)
	}
	if template.Category != "Linux" {
		t.Errorf("Expected category 'Linux', got '%s'", template.Category)
	}
	if !template.Active {
		t.Error("Expected active to be true")
	}
	if template.Status != 2 {
		t.Errorf("Expected status 2, got %d", template.Status)
	}
}
