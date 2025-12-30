package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ankitpokhrel/jira-cli/pkg/md"
)

// Comment represents a Jira comment.
type Comment struct {
	ID       string      `json:"id"`
	Author   User        `json:"author"`
	Body     interface{} `json:"body"` // string in v2, adf.ADF in v3
	Created  string      `json:"created"`
	Updated  string      `json:"updated,omitempty"`
	Visibility struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"visibility,omitempty"`
}

// CommentList represents a list of comments.
type CommentList struct {
	Comments []*Comment `json:"comments"`
	Total    int        `json:"total"`
	MaxResults int     `json:"maxResults"`
	StartAt  int        `json:"startAt"`
}

// GetComments retrieves all comments for an issue.
func (c *Client) GetComments(key string) ([]*Comment, error) {
	path := fmt.Sprintf("/issue/%s/comment", key)
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

	var commentList CommentList
	err = json.NewDecoder(res.Body).Decode(&commentList)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return commentList.Comments, nil
}

// UpdateComment updates an existing comment.
func (c *Client) UpdateComment(key, commentID, body string, internal bool) error {
	var bodyContent interface{}
	if internal {
		bodyContent = md.ToJiraMD(body)
	} else {
		bodyContent = md.ToJiraMD(body)
	}

	updateReq := struct {
		Body interface{} `json:"body"`
		Visibility *struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"visibility,omitempty"`
	}{
		Body: bodyContent,
	}

	if internal {
		updateReq.Visibility = &struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		}{
			Type:  "role",
			Value: "Administrators",
		}
	}

	bodyBytes, err := json.Marshal(updateReq)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/issue/%s/comment/%s", key, commentID)
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

// DeleteComment deletes a comment from an issue.
func (c *Client) DeleteComment(key, commentID string) error {
	path := fmt.Sprintf("/issue/%s/comment/%s", key, commentID)
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

