package xelon

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const snapshotBasePath = deviceBasePath + "/%v/snapshots"

// SnapshotsService handles communication with the snapshots related methods of the Xelon REST API.
type SnapshotsService service

type Snapshot struct {
	Description string     `json:"description,omitempty"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	ID          int        `json:"id,omitempty"`
	Name        string     `json:"name,omitempty"`
}

type SnapshotDeleteRequest struct {
	RemoveChildSnapshots bool `json:"removeChildren"`
}

type SnapshotListOptions struct {
	ListOptions
}

type snapshotsRoot struct {
	Snapshots []Snapshot `json:"data"`
	Meta      *Meta      `json:"meta,omitempty"`
}

func (v Snapshot) String() string { return Stringify(v) }

// List provides a list of all snapshots.
func (s *SnapshotsService) List(ctx context.Context, deviceID string, opts *SnapshotListOptions) ([]Snapshot, *Response, error) {
	if deviceID == "" {
		return nil, nil, errors.New("failed to list snapshots: device id must be supplied")
	}

	path, err := addOptions(fmt.Sprintf(snapshotBasePath, deviceID), opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(snapshotsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}

	return root.Snapshots, resp, err
}

// Delete removes snapshot identified by id.
func (s *SnapshotsService) Delete(ctx context.Context, deviceID string, snapshotID int, deleteRequest *SnapshotDeleteRequest) (*Response, error) {
	if deviceID == "" {
		return nil, errors.New("failed to delete snapshot: device id must be supplied")
	}
	if snapshotID == 0 {
		return nil, errors.New("failed to delete snapshot: id must be supplied")
	}
	if deleteRequest == nil {
		return nil, errors.New("failed to delete snapshot: payload must be supplied")
	}

	path := fmt.Sprintf(snapshotBasePath+"/%v", deviceID, snapshotID)
	req, err := s.client.NewRequest(http.MethodDelete, path, deleteRequest)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
