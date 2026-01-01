package jira

import "fmt"

// ErrAuthentication represents authentication failures.
type ErrAuthentication struct {
	Reason string
}

func (e *ErrAuthentication) Error() string {
	return fmt.Sprintf("authentication failed: %s", e.Reason)
}

// ErrNotFound represents resource not found errors.
type ErrNotFound struct {
	Resource string
	ID       string
}

func (e *ErrNotFound) Error() string {
	if e.ID != "" {
		return fmt.Sprintf("%s %q not found", e.Resource, e.ID)
	}
	return fmt.Sprintf("%s not found", e.Resource)
}

// ErrValidation represents input validation errors.
type ErrValidation struct {
	Field   string
	Message string
}

func (e *ErrValidation) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation error for %q: %s", e.Field, e.Message)
	}
	return fmt.Sprintf("validation error: %s", e.Message)
}

// ErrRateLimit represents rate limiting errors.
type ErrRateLimit struct {
	RetryAfter int // seconds
}

func (e *ErrRateLimit) Error() string {
	if e.RetryAfter > 0 {
		return fmt.Sprintf("rate limit exceeded, retry after %d seconds", e.RetryAfter)
	}
	return "rate limit exceeded"
}

// ErrNetwork represents network-related errors.
type ErrNetwork struct {
	Underlying error
}

func (e *ErrNetwork) Error() string {
	return fmt.Sprintf("network error: %v", e.Underlying)
}

func (e *ErrNetwork) Unwrap() error {
	return e.Underlying
}

// IsRetryableError checks if an error is retryable.
func IsRetryableError(err error) bool {
	switch err.(type) {
	case *ErrNetwork:
		return true
	case *ErrRateLimit:
		return true
	case *ErrUnexpectedResponse:
		// Retry on 5xx errors
		if e, ok := err.(*ErrUnexpectedResponse); ok {
			return e.StatusCode >= 500 && e.StatusCode < 600
		}
	}
	return false
}


