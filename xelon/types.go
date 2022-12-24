package xelon

// APIResponse is a generic Xelon API response.
type APIResponse struct {
	Message           string             `json:"message,omitempty"`
	PersistentStorage *PersistentStorage `json:"persistentStorage,omitempty"`
}
