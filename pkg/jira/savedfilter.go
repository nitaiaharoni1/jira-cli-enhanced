package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// SavedFilter represents a saved Jira filter.
type SavedFilter struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	JQL         string `json:"jql"`
	Owner       struct {
		AccountID   string `json:"accountId"`
		DisplayName string `json:"displayName"`
	} `json:"owner"`
	SharePermissions []struct {
		Type string `json:"type"`
	} `json:"sharePermissions"`
	Favourite bool   `json:"favourite"`
	Self      string `json:"self"`
}

// SavedFilterListResponse represents the response from listing filters.
type SavedFilterListResponse struct {
	MaxResults int            `json:"maxResults"`
	StartAt    int            `json:"startAt"`
	Total      int            `json:"total"`
	IsLast     bool           `json:"isLast"`
	Values     []*SavedFilter `json:"values"`
}

// CreateFilterRequest represents the request to create a filter.
type CreateFilterRequest struct {
	Name             string   `json:"name"`
	Description      string   `json:"description,omitempty"`
	JQL              string   `json:"jql"`
	Favourite        bool     `json:"favourite,omitempty"`
	SharePermissions []string `json:"sharePermissions,omitempty"`
}

// UpdateFilterRequest represents the request to update a filter.
type UpdateFilterRequest struct {
	Name             string   `json:"name,omitempty"`
	Description      string   `json:"description,omitempty"`
	JQL              string   `json:"jql,omitempty"`
	Favourite        *bool   `json:"favourite,omitempty"`
	SharePermissions []string `json:"sharePermissions,omitempty"`
}

// GetFilters fetches user's saved filters from /rest/api/2/filter endpoint.
func (c *Client) GetFilters() ([]*SavedFilter, error) {
	res, err := c.GetV2(context.Background(), "/filter/favourite", nil)
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

	var filters []*SavedFilter
	err = json.NewDecoder(res.Body).Decode(&filters)
	return filters, err
}

// GetAllFilters fetches all filters accessible to the user.
func (c *Client) GetAllFilters(startAt, maxResults int) (*SavedFilterListResponse, error) {
	path := fmt.Sprintf("/filter/search?startAt=%d&maxResults=%d", startAt, maxResults)
	res, err := c.GetV2(context.Background(), path, nil)
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

	var out SavedFilterListResponse
	err = json.NewDecoder(res.Body).Decode(&out)
	return &out, err
}

// GetFilter fetches a specific filter by ID.
func (c *Client) GetFilter(filterID string) (*SavedFilter, error) {
	res, err := c.GetV2(context.Background(), fmt.Sprintf("/filter/%s", filterID), nil)
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

	var filter SavedFilter
	err = json.NewDecoder(res.Body).Decode(&filter)
	return &filter, err
}

// CreateFilter creates a new saved filter.
func (c *Client) CreateFilter(req *CreateFilterRequest) (*SavedFilter, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	res, err := c.PostV2(context.Background(), "/filter", body, Header{
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

	var filter SavedFilter
	err = json.NewDecoder(res.Body).Decode(&filter)
	return &filter, err
}

// UpdateFilter updates an existing filter.
func (c *Client) UpdateFilter(filterID string, req *UpdateFilterRequest) (*SavedFilter, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	res, err := c.PutV2(context.Background(), fmt.Sprintf("/filter/%s", filterID), body, Header{
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

	var filter SavedFilter
	err = json.NewDecoder(res.Body).Decode(&filter)
	return &filter, err
}

// DeleteFilter deletes a filter.
func (c *Client) DeleteFilter(filterID string) error {
	res, err := c.DeleteV2(context.Background(), fmt.Sprintf("/filter/%s", filterID), nil)
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

// ExecuteFilter executes a filter and returns the matching issues.
func (c *Client) ExecuteFilter(filterID string, limit uint) (*SearchResult, error) {
	filter, err := c.GetFilter(filterID)
	if err != nil {
		return nil, err
	}

	// Use the filter's JQL to search for issues
	// SearchV2 takes from and limit parameters
	return c.SearchV2(filter.JQL, 0, limit)
}

