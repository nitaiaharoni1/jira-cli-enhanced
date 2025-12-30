package jira

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// Attachment represents a Jira attachment.
type Attachment struct {
	ID          string `json:"id"`
	Self        string `json:"self"`
	Filename    string `json:"filename"`
	Author      User   `json:"author"`
	Created     string `json:"created"`
	Size        int64  `json:"size"`
	MimeType    string `json:"mimeType"`
	Content     string `json:"content"`
	Thumbnail   string `json:"thumbnail,omitempty"`
}

// AttachmentList represents a list of attachments.
type AttachmentList struct {
	Attachments []*Attachment `json:"attachments"`
}

// UploadAttachment uploads a file as an attachment to an issue.
func (c *Client) UploadAttachment(key string, filePath string) ([]*Attachment, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create multipart form
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add file to form
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	// Create request with multipart form data
	endpoint := fmt.Sprintf("%s/rest/api/2/issue/%s/attachments", c.server, key)
	req, err := http.NewRequest(http.MethodPost, endpoint, &requestBody)
	if err != nil {
		return nil, &ErrNetwork{Underlying: err}
	}

	// Set headers
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Atlassian-Token", "no-check")

	// Set authentication
	if c.authType == nil {
		basic := AuthTypeBasic
		c.authType = &basic
	}

	switch c.authType.String() {
	case string(AuthTypeMTLS):
		if c.token != "" {
			req.Header.Add("Authorization", "Bearer "+c.token)
		}
	case string(AuthTypeBearer):
		req.Header.Add("Authorization", "Bearer "+c.token)
	case string(AuthTypeBasic):
		req.SetBasicAuth(c.login, c.token)
	}

	// Execute request
	httpClient := c.getHTTPClient()
	res, err := httpClient.Do(req.WithContext(context.Background()))
	if err != nil {
		return nil, &ErrNetwork{Underlying: err}
	}

	if res == nil {
		return nil, ErrEmptyResponse
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return nil, formatUnexpectedResponse(res)
	}

	var attachments []*Attachment
	err = json.NewDecoder(res.Body).Decode(&attachments)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return attachments, nil
}

// GetAttachments retrieves all attachments for an issue.
func (c *Client) GetAttachments(key string) ([]*Attachment, error) {
	issue, err := c.GetIssueV2(key)
	if err != nil {
		return nil, err
	}

	return issue.Fields.Attachments, nil
}

// DeleteAttachment deletes an attachment from an issue.
func (c *Client) DeleteAttachment(attachmentID string) error {
	endpoint := fmt.Sprintf("/attachment/%s", attachmentID)
	res, err := c.DeleteV2(context.Background(), endpoint, Header{
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

// DownloadAttachment downloads an attachment to a local file.
func (c *Client) DownloadAttachment(attachmentID, filePath string) error {
	endpoint := fmt.Sprintf("/attachment/%s", attachmentID)
	res, err := c.GetV2(context.Background(), endpoint, Header{
		"Accept": "application/json",
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

	var attachment Attachment
	err = json.NewDecoder(res.Body).Decode(&attachment)
	if err != nil {
		return fmt.Errorf("failed to decode attachment: %w", err)
	}

	// Download the actual file content
	downloadRes, err := c.GetV2(context.Background(), fmt.Sprintf("/attachment/content/%s", attachmentID), Header{})
	if err != nil {
		return err
	}
	if downloadRes == nil {
		return ErrEmptyResponse
	}
	defer func() { _ = downloadRes.Body.Close() }()

	if downloadRes.StatusCode != http.StatusOK {
		return formatUnexpectedResponse(downloadRes)
	}

	// Create output file
	outFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, downloadRes.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

