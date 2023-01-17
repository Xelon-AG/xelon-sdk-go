package xelon

// APIResponse is a generic Xelon API response.
type APIResponse struct {
	Message           string             `json:"message,omitempty"`
	PersistentStorage *PersistentStorage `json:"persistentStorage,omitempty"`
}

type Meta struct {
	Page    int `json:"current_page,omitempty"`
	PerPage int `json:"per_page,omitempty"`
	Total   int `json:"total,omitempty"`
}
