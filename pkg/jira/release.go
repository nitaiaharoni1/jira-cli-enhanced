package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Release fetches response from /project/{projectIdOrKey}/version endpoint.
func (c *Client) Release(project string) ([]*ProjectVersion, error) {
	path := fmt.Sprintf("/project/%s/versions", project)
	res, err := c.Get(context.Background(), path, nil)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, ErrEmptyResponse
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return nil, formatUnexpectedResponse(res)
	}

	var out []*ProjectVersion

	err = json.NewDecoder(res.Body).Decode(&out)

	return out, err
}

// CreateVersionRequest represents the request to create a version.
type CreateVersionRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	Archived     bool   `json:"archived,omitempty"`
	Released     bool   `json:"released,omitempty"`
	ReleaseDate  string `json:"releaseDate,omitempty"`
	StartDate    string `json:"startDate,omitempty"`
	ProjectID    string `json:"projectId,omitempty"`
	Project      string `json:"project,omitempty"`
}

// UpdateVersionRequest represents the request to update a version.
type UpdateVersionRequest struct {
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	Archived     *bool  `json:"archived,omitempty"`
	Released     *bool  `json:"released,omitempty"`
	ReleaseDate  string `json:"releaseDate,omitempty"`
	StartDate    string `json:"startDate,omitempty"`
}

// CreateVersion creates a new project version.
func (c *Client) CreateVersion(req *CreateVersionRequest) (*ProjectVersion, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	res, err := c.PostV2(context.Background(), "/version", body, Header{
		"Content-Type": "application/json",
	})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, ErrEmptyResponse
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		return nil, formatUnexpectedResponse(res)
	}

	var version ProjectVersion
	err = json.NewDecoder(res.Body).Decode(&version)
	return &version, err
}

// UpdateVersion updates an existing project version.
func (c *Client) UpdateVersion(versionID string, req *UpdateVersionRequest) (*ProjectVersion, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	res, err := c.PutV2(context.Background(), fmt.Sprintf("/version/%s", versionID), body, Header{
		"Content-Type": "application/json",
	})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, ErrEmptyResponse
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return nil, formatUnexpectedResponse(res)
	}

	var version ProjectVersion
	err = json.NewDecoder(res.Body).Decode(&version)
	return &version, err
}

// DeleteVersion deletes a project version.
func (c *Client) DeleteVersion(versionID string, moveFixIssuesTo string, moveAffectedIssuesTo string) error {
	path := fmt.Sprintf("/version/%s", versionID)
	if moveFixIssuesTo != "" || moveAffectedIssuesTo != "" {
		path += "?"
		params := []string{}
		if moveFixIssuesTo != "" {
			params = append(params, fmt.Sprintf("moveFixIssuesTo=%s", moveFixIssuesTo))
		}
		if moveAffectedIssuesTo != "" {
			params = append(params, fmt.Sprintf("moveAffectedIssuesTo=%s", moveAffectedIssuesTo))
		}
		path += fmt.Sprintf("%s", params[0])
		if len(params) > 1 {
			path += "&" + params[1]
		}
	}

	res, err := c.DeleteV2(context.Background(), path, nil)
	if err != nil {
		return err
	}
	if res == nil {
		return ErrEmptyResponse
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusNoContent && res.StatusCode != http.StatusOK {
		return formatUnexpectedResponse(res)
	}

	return nil
}
