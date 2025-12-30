package cmdutil

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

var (
	issueKeyRegex = regexp.MustCompile(`^[A-Z]+-\d+$`)
)

// ValidateIssueKey validates an issue key format.
func ValidateIssueKey(key string) error {
	if key == "" {
		return &jira.ErrValidation{
			Field:   "issue-key",
			Message: "issue key cannot be empty",
		}
	}
	if !issueKeyRegex.MatchString(key) {
		return &jira.ErrValidation{
			Field:   "issue-key",
			Message: fmt.Sprintf("invalid format: %s (expected: PROJECT-123)", key),
		}
	}
	return nil
}

// ValidateProjectKey validates a project key format.
func ValidateProjectKey(key string) error {
	if key == "" {
		return &jira.ErrValidation{
			Field:   "project-key",
			Message: "project key cannot be empty",
		}
	}
	if len(key) < 2 || len(key) > 10 {
		return &jira.ErrValidation{
			Field:   "project-key",
			Message: fmt.Sprintf("project key must be 2-10 characters: %s", key),
		}
	}
	if !regexp.MustCompile(`^[A-Z]+$`).MatchString(key) {
		return &jira.ErrValidation{
			Field:   "project-key",
			Message: fmt.Sprintf("project key must contain only uppercase letters: %s", key),
		}
	}
	return nil
}

// ValidateServerURL validates a server URL format.
func ValidateServerURL(url string) error {
	if url == "" {
		return &jira.ErrValidation{
			Field:   "server",
			Message: "server URL cannot be empty",
		}
	}
	url = strings.TrimSuffix(url, "/")
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return &jira.ErrValidation{
			Field:   "server",
			Message: fmt.Sprintf("server URL must start with http:// or https://: %s", url),
		}
	}
	return nil
}

// ValidateSprintID validates a sprint ID.
func ValidateSprintID(id string) error {
	if id == "" {
		return &jira.ErrValidation{
			Field:   "sprint-id",
			Message: "sprint ID cannot be empty",
		}
	}
	return nil
}

