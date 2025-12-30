package cmdutil

import (
	"html"
	"strings"

	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

var (
	// dangerousJQLPatterns are patterns that could be used for JQL injection.
	dangerousJQLPatterns = []string{"';", "--", "/*", "*/", "xp_", "sp_"}
)

// SanitizeString removes potentially dangerous characters and escapes HTML.
func SanitizeString(input string) string {
	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")
	// Escape HTML
	input = html.EscapeString(input)
	return strings.TrimSpace(input)
}

// ValidateJQL validates JQL to prevent injection attacks.
func ValidateJQL(jql string) error {
	if jql == "" {
		return nil
	}

	lowerJQL := strings.ToLower(jql)
	for _, pattern := range dangerousJQLPatterns {
		if strings.Contains(lowerJQL, pattern) {
			return &jira.ErrValidation{
				Field:   "jql",
				Message: "potentially dangerous JQL pattern detected",
			}
		}
	}

	return nil
}

// SanitizeIssueKey sanitizes an issue key input.
func SanitizeIssueKey(key string) string {
	// Remove any whitespace
	key = strings.TrimSpace(key)
	// Convert to uppercase
	key = strings.ToUpper(key)
	return key
}

