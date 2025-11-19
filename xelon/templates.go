package xelon

import (
	"context"
	"net/http"
)

// TemplatesService handles communication with the template related methods of the Xelon API v2.
type TemplatesService service

// TemplateV2 represents a template in the Xelon API v2 format.
type TemplateV2 struct {
	Identifier                   string  `json:"identifier"`
	Name                         string  `json:"name"`
	OwnerTenant                  *string `json:"ownerTenant"`
	CloudIdentifier              string  `json:"cloudIdentifier"`
	Status                       int     `json:"status"`
	Description                  string  `json:"description"`
	InternalVersion              string  `json:"internalVersion"`
	Active                       bool    `json:"active"`
	TemplateType                 string  `json:"templateType"`
	ScriptType                   *string `json:"scriptType"`
	NICCountMin                  *int    `json:"nicCountMin"`
	NICCountMax                  *int    `json:"nicCountMax"`
	NICDescription               *string `json:"nicDescription"`
	HDDUnit                      int     `json:"hddUnit"`
	CloudType                    string  `json:"cloudType"`
	HDDMinSizeGB                 *int    `json:"hddMinSizeGb"`
	HasSwap                      *bool   `json:"hasSwap"`
	AllowedSetPwd                *bool   `json:"allowedSetPwd"`
	AllowedSetDisk               *bool   `json:"allowedSetDisk"`
	AllowedUpdateIPGuest         *bool   `json:"allowedUpdateIpGuest"`
	AllowedResizeExistingPartition *bool `json:"allowedResizeExistingPartition"`
	HasServiceUser               *bool   `json:"hasServiceUser"`
	Type                         string  `json:"type"`
	Category                     string  `json:"category"`
	NICUnit                      *int    `json:"nicUnit"`
	SwapUnit                     *int    `json:"swapUnit"`
	VMDKId                       *string `json:"vmdkId"`
	VMDKSwapId                   *string `json:"vmdkSwapId"`
	Instructions                 *string `json:"instructions"`
	PowerOn                      *bool   `json:"powerOn"`
	RegExp                       *string `json:"regExp"`
	RegExpId                     *int    `json:"regExpId"`
	CreatedAt                    string  `json:"createdAt"`
	UpdatedAt                    string  `json:"updatedAt"`
	DeletedAt                    *string `json:"deletedAt"`
}

// TemplateListResponse represents the paginated response from the templates list endpoint.
type TemplateListResponse struct {
	Data []TemplateV2   `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

// PaginationMeta represents pagination metadata.
type PaginationMeta struct {
	Total       int `json:"total"`
	LastPage    int `json:"lastPage"`
	PerPage     int `json:"perPage"`
	CurrentPage int `json:"currentPage"`
	From        int `json:"from"`
	To          int `json:"to"`
}

// TemplateListOptions specifies the optional parameters for listing templates.
type TemplateListOptions struct {
	// Page is the page number to retrieve (set automatically during pagination)
	Page int `url:"page,omitempty"`

	// PerPage is the number of items to return per page (default: 100)
	PerPage int `url:"perPage,omitempty"`

	// Sort specifies the field to sort by: "name", "description", "type", "category"
	Sort string `url:"sort,omitempty"`

	// Search is a search string to filter templates
	Search string `url:"search,omitempty"`

	// TenantIdentifiers filters templates by specific tenant IDs
	TenantIdentifiers []string `url:"tenantIdentifiers[],omitempty"`

	// CloudIdentifier filters templates by cloud ID
	CloudIdentifier string `url:"cloudIdentifier,omitempty"`

	// Type filters by template type: "Linux", "Windows", "pfSense", "Sophos", "Barracuda", "ISO", "OPNsense", "SonicWALL"
	Type string `url:"type,omitempty"`

	// IsActive filters by active status (default: true)
	IsActive *bool `url:"isActive,omitempty"`

	// IsCustom filters custom vs company templates (true=custom only, false=company only, nil=all)
	IsCustom *bool `url:"isCustom,omitempty"`
}

// List provides a paginated list of available templates from the v2 API with optional filtering.
// If opts is nil, defaults to: perPage=100. API defaults isActive=true if not specified.
func (s *TemplatesService) List(ctx context.Context, opts *TemplateListOptions) ([]TemplateV2, *Response, error) {
	// Set defaults
	if opts == nil {
		opts = &TemplateListOptions{}
	}
	if opts.PerPage == 0 {
		opts.PerPage = 100 // Request maximum per page for efficiency
	}
	// Don't set IsActive - let API use its default (true)

	allTemplates := make([]TemplateV2, 0)
	opts.Page = 1

	for {
		// Build URL with query parameters
		path, err := addOptions("templates", opts)
		if err != nil {
			return nil, nil, err
		}

		req, err := s.client.NewRequest(http.MethodGet, path, nil)
		if err != nil {
			return nil, nil, err
		}

		templateResp := new(TemplateListResponse)
		resp, err := s.client.Do(ctx, req, templateResp)
		if err != nil {
			return nil, resp, err
		}

		allTemplates = append(allTemplates, templateResp.Data...)

		// Check if there are more pages
		if opts.Page >= templateResp.Meta.LastPage {
			break
		}
		opts.Page++
	}

	return allTemplates, nil, nil
}

// Get retrieves a specific template by identifier.
func (s *TemplatesService) Get(ctx context.Context, identifier string) (*TemplateV2, *Response, error) {
	// Since there's no direct Get endpoint by identifier, we list all and filter
	// Note: Can't use search parameter as it searches by name, not identifier
	templates, _, err := s.List(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	for _, t := range templates {
		if t.Identifier == identifier {
			return &t, nil, nil
		}
	}

	return nil, nil, nil // Not found
}

// GetByName retrieves a specific template by name using the search parameter for efficiency.
func (s *TemplatesService) GetByName(ctx context.Context, name string, opts *TemplateListOptions) (*TemplateV2, *Response, error) {
	// Use search parameter for server-side filtering
	if opts == nil {
		opts = &TemplateListOptions{}
	}
	opts.Search = name

	templates, _, err := s.List(ctx, opts)
	if err != nil {
		return nil, nil, err
	}

	// Search for exact name match (API search might return partial matches)
	for _, t := range templates {
		if t.Name == name {
			return &t, nil, nil
		}
	}

	return nil, nil, nil // Not found
}
