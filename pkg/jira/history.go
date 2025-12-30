package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// HistoryEntry represents a single changelog entry.
type HistoryEntry struct {
	ID      string `json:"id"`
	Author  User   `json:"author"`
	Created string `json:"created"`
	Items   []HistoryItem `json:"items"`
}

// HistoryItem represents a single field change in a changelog entry.
type HistoryItem struct {
	Field      string `json:"field"`
	FieldType  string `json:"fieldtype"`
	From       interface{} `json:"from"`
	FromString string `json:"fromString"`
	To         interface{} `json:"to"`
	ToString   string `json:"toString"`
}

// GetIssueHistory retrieves the changelog/history for an issue.
func (c *Client) GetIssueHistory(key string) ([]HistoryEntry, error) {
	path := fmt.Sprintf("/issue/%s?expand=changelog", key)
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

	var issue struct {
		Changelog struct {
			Histories []HistoryEntry `json:"histories"`
		} `json:"changelog"`
	}

	err = json.NewDecoder(res.Body).Decode(&issue)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Flatten history entries with items
	var result []HistoryEntry
	for _, history := range issue.Changelog.Histories {
		for _, item := range history.Items {
			result = append(result, HistoryEntry{
				ID:      history.ID,
				Author:  history.Author,
				Created: history.Created,
				Items:   []HistoryItem{item},
			})
		}
	}

	return result, nil
}

// GetIssueHistoryFlat returns flattened history entries (one per field change).
func (c *Client) GetIssueHistoryFlat(key string) ([]struct {
	ID         string
	Author     User
	Created    string
	Field      string
	FromString string
	ToString   string
}, error) {
	history, err := c.GetIssueHistory(key)
	if err != nil {
		return nil, err
	}

	var result []struct {
		ID         string
		Author     User
		Created    string
		Field      string
		FromString string
		ToString   string
	}

	for _, h := range history {
		for _, item := range h.Items {
			result = append(result, struct {
				ID         string
				Author     User
				Created    string
				Field      string
				FromString string
				ToString   string
			}{
				ID:         h.ID,
				Author:     h.Author,
				Created:    h.Created,
				Field:      item.Field,
				FromString: item.FromString,
				ToString:   item.ToString,
			})
		}
	}

	return result, nil
}

