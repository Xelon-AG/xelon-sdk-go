package xelon

// Meta represents an pagination object.
type Meta struct {
	From     int `json:"from,omitempty"`        // From is the starting index on the current page.
	LastPage int `json:"lastPage,omitempty"`    // LastPage is the last available page number.
	Page     int `json:"currentPage,omitempty"` // Page is the current page of the pagination.
	PerPage  int `json:"perPage,omitempty"`     // PerPage is the number of items displayed per page.
	To       int `json:"to,omitempty"`          // To is the ending index on the current page.
	Total    int `json:"total,omitempty"`       // Total is the total number of items.
}

func (v Meta) String() string { return Stringify(v) }

// nilToEmpty returns an empty (but non-nil) slice if s is nil, otherwise s.
// Used in MarshalJSON to ensure empty slices render as "[]" instead of "null".
func nilToEmpty[T any](s []T) []T {
	if s == nil {
		return []T{}
	}
	return s
}
