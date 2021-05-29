package mangodex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	CreateScanGroupPath   = "group"
	ViewScanGroupPath     = "group/%s"
	UpdateScanGroupPath   = ViewScanGroupPath
	DeleteScanGroupPath   = ViewScanGroupPath
	FollowScanGroupPath   = "group/%s/follow"
	UnfollowScanGroupPath = FollowScanGroupPath
)

type ScanGroupList struct {
	Results []ScanGroupResponse `json:"results"`
	Limit   int                 `json:"limit"`
	Offset  int                 `json:"offset"`
	Total   int                 `json:"total"`
}

type ScanGroupResponse struct {
	Result        string         `json:"result"`
	Data          ScanGroup      `json:"data"`
	Relationships []Relationship `json:"relationships"`
}

func (r *ScanGroupResponse) GetResult() string {
	return r.Result
}

type ScanGroup struct {
	ID         string              `json:"id"`
	Type       string              `json:"string"`
	Attributes ScanGroupAttributes `json:"attributes"`
}

type ScanGroupAttributes struct {
	Name      string `json:"name"`
	Leader    User   `json:"leader"`
	Version   int    `json:"version"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// CreateScanGroup : Create a scanlation group.
// https://api.mangadex.org/docs.html#operation/post-group
func (dc *DexClient) CreateScanGroup(name, leader string, members []string, version int) (*ScanGroupResponse, error) {
	return dc.CreateScanGroupContext(context.Background(), name, leader, members, version)
}

// CreateScanGroupContext : CreateScanGroup with custom context.
func (dc *DexClient) CreateScanGroupContext(ctx context.Context, name, leader string, members []string, version int) (*ScanGroupResponse, error) {
	// Create request body.
	req := struct {
		Name    string
		Leader  string
		Members []string
		Version int
	}{Name: name, Leader: leader, Members: members, Version: version}
	rBytes, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	var r ScanGroupResponse
	err = dc.responseOp(ctx, http.MethodPost, CreateScanGroupPath, bytes.NewBuffer(rBytes), &r)
	return &r, err
}

// ViewScanGroup : View a scanlation group.
// https://api.mangadex.org/docs.html#operation/get-group-id
func (dc *DexClient) ViewScanGroup(id string) (*ScanGroupResponse, error) {
	return dc.ViewScanGroupContext(context.Background(), id)
}

// ViewScanGroupContext : ViewScanGroup with custom context.
func (dc *DexClient) ViewScanGroupContext(ctx context.Context, id string) (*ScanGroupResponse, error) {
	var r ScanGroupResponse
	err := dc.responseOp(ctx, http.MethodGet, fmt.Sprintf(ViewScanGroupPath, id), nil, &r)
	return &r, err
}

// UpdateScanGroup : Update a scanlation group.
// https://api.mangadex.org/docs.html#operation/put-group-id
func (dc *DexClient) UpdateScanGroup(id, name, leader string, members []string, version int) (*ScanGroupResponse, error) {
	return dc.UpdateScanGroupContext(context.Background(), id, name, leader, members, version)
}

// UpdateScanGroupContext : UpdateScanGroup with custom context.
func (dc *DexClient) UpdateScanGroupContext(ctx context.Context, id, name, leader string, members []string, version int) (*ScanGroupResponse, error) {
	// Create request body.
	req := struct {
		Name    string
		Leader  string
		Members []string
		Version int
	}{Name: name, Leader: leader, Members: members, Version: version}
	rBytes, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	var r ScanGroupResponse
	err = dc.responseOp(ctx, http.MethodPut, fmt.Sprintf(UpdateScanGroupPath, id), bytes.NewBuffer(rBytes), &r)
	return &r, err
}

// DeleteScanGroup : Delete a scanlation group.
// https://api.mangadex.org/docs.html#operation/delete-group-id
func (dc *DexClient) DeleteScanGroup(id string) error {
	return dc.DeleteScanGroupContext(context.Background(), id)
}

// DeleteScanGroupContext : DeleteScanGroup with custom context.
func (dc *DexClient) DeleteScanGroupContext(ctx context.Context, id string) error {
	return dc.responseOp(ctx, http.MethodDelete, fmt.Sprintf(DeleteScanGroupPath, id), nil, nil)
}

// FollowScanGroup : Follow a scanlation group.
// https://api.mangadex.org/docs.html#operation/post-group-id-follow
func (dc *DexClient) FollowScanGroup(id string) error {
	return dc.FollowScanGroupContext(context.Background(), id)
}

// FollowScanGroupContext : FollowScanGroup with custom context.
func (dc *DexClient) FollowScanGroupContext(ctx context.Context, id string) error {
	return dc.responseOp(ctx, http.MethodPost, fmt.Sprintf(FollowScanGroupPath, id), nil, nil)
}

// UnfollowScanGroup : Unfollow a scanlation group.
// https://api.mangadex.org/docs.html#operation/delete-group-id-follow
func (dc *DexClient) UnfollowScanGroup(id string) error {
	return dc.UnfollowScanGroupContext(context.Background(), id)
}

// UnfollowScanGroupContext : UnfollowScanGroup with custom context.
func (dc *DexClient) UnfollowScanGroupContext(ctx context.Context, id string) error {
	return dc.responseOp(ctx, http.MethodDelete, fmt.Sprintf(UnfollowScanGroupPath, id), nil, nil)
}
