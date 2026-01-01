package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ankitpokhrel/jira-cli/pkg/md"
)

// Worklog represents a Jira worklog entry.
type Worklog struct {
	ID          string `json:"id"`
	Author      User   `json:"author"`
	Comment     string `json:"comment"`
	Created     string `json:"created"`
	Updated     string `json:"updated,omitempty"`
	Started     string `json:"started"`
	TimeSpent   string `json:"timeSpent"`
	TimeSpentSeconds int `json:"timeSpentSeconds"`
}

// WorklogList represents a list of worklogs.
type WorklogList struct {
	Worklogs   []*Worklog `json:"worklogs"`
	Total      int        `json:"total"`
	MaxResults int        `json:"maxResults"`
	StartAt    int        `json:"startAt"`
}

// GetWorklogs retrieves all worklogs for an issue.
func (c *Client) GetWorklogs(key string) ([]*Worklog, error) {
	path := fmt.Sprintf("/issue/%s/worklog", key)
	res, err := c.GetV2(context.Background(), path, Header{
		"Accept": "application/json",
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

	var worklogList WorklogList
	err = json.NewDecoder(res.Body).Decode(&worklogList)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return worklogList.Worklogs, nil
}

// UpdateWorklog updates an existing worklog entry.
func (c *Client) UpdateWorklog(key, worklogID, started, timeSpent, comment string) error {
	updateReq := struct {
		Started   string `json:"started,omitempty"`
		TimeSpent string `json:"timeSpent,omitempty"`
		Comment   string `json:"comment,omitempty"`
	}{
		TimeSpent: timeSpent,
		Comment:   md.ToJiraMD(comment),
	}

	if started != "" {
		updateReq.Started = started
	}

	bodyBytes, err := json.Marshal(updateReq)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/issue/%s/worklog/%s", key, worklogID)
	res, err := c.PutV2(context.Background(), path, bodyBytes, Header{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	})
	if err != nil {
		return err
	}
	if res == nil {
		return ErrEmptyResponse
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return formatUnexpectedResponse(res)
	}

	return nil
}

// DeleteWorklog deletes a worklog entry from an issue.
func (c *Client) DeleteWorklog(key, worklogID string) error {
	path := fmt.Sprintf("/issue/%s/worklog/%s", key, worklogID)
	res, err := c.DeleteV2(context.Background(), path, Header{
		"Accept": "application/json",
	})
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


