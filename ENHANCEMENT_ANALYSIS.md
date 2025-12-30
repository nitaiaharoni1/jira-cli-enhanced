# Jira CLI Enhancement Analysis

Based on Context7 best practices and code review, here are actionable improvements for the jira-cli project.

## 1. Error Handling Improvements

### Current Issues
- Error handling is inconsistent across commands
- Some commands use `Run` instead of `RunE` (can't return errors properly)
- Error messages could be more actionable
- No structured error types for common scenarios

### Recommendations

#### 1.1 Use `RunE` Instead of `Run` for Error Handling
**Current Pattern:**
```go
Run: func(cmd *cobra.Command, args []string) {
    // errors handled internally with os.Exit
}
```

**Recommended Pattern:**
```go
RunE: func(cmd *cobra.Command, args []string) error {
    if err := doSomething(); err != nil {
        return fmt.Errorf("failed to do something: %w", err)
    }
    return nil
}
```

**Benefits:**
- Allows Cobra to handle errors consistently
- Enables better error wrapping and context
- Makes testing easier

#### 1.2 Create Structured Error Types
Add common error types in `pkg/jira/errors.go`:

```go
// ErrAuthentication represents authentication failures
type ErrAuthentication struct {
    Reason string
}

func (e *ErrAuthentication) Error() string {
    return fmt.Sprintf("authentication failed: %s", e.Reason)
}

// ErrNotFound represents resource not found errors
type ErrNotFound struct {
    Resource string
    ID       string
}

func (e *ErrNotFound) Error() string {
    return fmt.Sprintf("%s %q not found", e.Resource, e.ID)
}

// ErrValidation represents input validation errors
type ErrValidation struct {
    Field   string
    Message string
}

func (e *ErrValidation) Error() string {
    return fmt.Sprintf("validation error for %q: %s", e.Field, e.Message)
}
```

#### 1.3 Improve Error Messages with Context
Enhance `cmdutil.ExitIfError` to provide actionable suggestions:

```go
func ExitIfError(err error) {
    if err == nil {
        return
    }

    var msg string
    var suggestion string

    switch e := err.(type) {
    case *jira.ErrAuthentication:
        msg = fmt.Sprintf("Authentication failed: %s", e.Reason)
        suggestion = "Run 'jira init' to reconfigure your credentials"
    case *jira.ErrNotFound:
        msg = fmt.Sprintf("%s not found", e.Resource)
        suggestion = "Verify the ID is correct and you have access"
    case *jira.ErrUnexpectedResponse:
        msg = fmt.Sprintf("Received unexpected response '%s'", e.Status)
        suggestion = "Check your parameters and try again"
    default:
        msg = fmt.Sprintf("Error: %s", err.Error())
    }

    fmt.Fprintf(os.Stderr, "%s\n", msg)
    if suggestion != "" {
        fmt.Fprintf(os.Stderr, "💡 %s\n", suggestion)
    }
    os.Exit(1)
}
```

## 2. Command Structure Improvements

### 2.1 Consistent Argument Validation
Use Cobra's built-in argument validation:

```go
// Instead of manual validation in Run/RunE
Args: cobra.ExactArgs(1), // or MinimumNArgs, MaximumNArgs, etc.

// For custom validation
Args: func(cmd *cobra.Command, args []string) error {
    if len(args) == 0 {
        return fmt.Errorf("issue key is required")
    }
    if !isValidIssueKey(args[0]) {
        return fmt.Errorf("invalid issue key format: %s", args[0])
    }
    return nil
}
```

### 2.2 Mark Required Flags Explicitly
```go
cmd.Flags().StringVarP(&project, "project", "p", "", "Project key")
cmd.MarkFlagRequired("project") // Better than checking manually
```

### 2.3 Better Help Text Organization
Add examples and better descriptions:

```go
var cmdIssue = &cobra.Command{
    Use:   "issue <command>",
    Short: "Manage Jira issues",
    Long: `Manage Jira issues including creation, editing, viewing, and listing.
    
All issue commands support filtering and can be used interactively or
non-interactively. Use --help on any subcommand for more details.`,
    Example: `  # List issues assigned to you
  jira issue list -a$(jira me)
  
  # Create a new bug
  jira issue create -tBug -s"Bug title" -yHigh`,
}
```

## 3. Logging and Debug Improvements

### 3.1 Structured Logging
Create a logger package (`internal/logger/logger.go`):

```go
package logger

import (
    "fmt"
    "os"
    "time"
)

type Level int

const (
    LevelDebug Level = iota
    LevelInfo
    LevelWarn
    LevelError
)

type Logger struct {
    level Level
    debug bool
}

func New(debug bool) *Logger {
    level := LevelInfo
    if debug {
        level = LevelDebug
    }
    return &Logger{level: level, debug: debug}
}

func (l *Logger) Debug(format string, args ...interface{}) {
    if l.level <= LevelDebug {
        l.log("DEBUG", format, args...)
    }
}

func (l *Logger) Info(format string, args ...interface{}) {
    if l.level <= LevelInfo {
        l.log("INFO", format, args...)
    }
}

func (l *Logger) Warn(format string, args ...interface{}) {
    if l.level <= LevelWarn {
        l.log("WARN", format, args...)
    }
}

func (l *Logger) Error(format string, args ...interface{}) {
    if l.level <= LevelError {
        l.log("ERROR", format, args...)
    }
}

func (l *Logger) log(level, format string, args ...interface{}) {
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    msg := fmt.Sprintf(format, args...)
    fmt.Fprintf(os.Stderr, "[%s] %s: %s\n", timestamp, level, msg)
}
```

### 3.2 Request/Response Logging
Enhance debug output in `pkg/jira/client.go`:

```go
func (c *Client) request(ctx context.Context, method, endpoint string, body []byte, headers Header) (*http.Response, error) {
    if c.debug {
        logger.Debug("Request: %s %s", method, endpoint)
        if len(body) > 0 {
            logger.Debug("Request body: %s", string(body))
        }
    }
    
    // ... existing code ...
    
    if c.debug && res != nil {
        logger.Debug("Response: %s %s", res.Status, res.Header.Get("Content-Type"))
    }
    
    return res, err
}
```

## 4. Testing Improvements

### 4.1 Add Table-Driven Tests
For utility functions, use table-driven tests:

```go
func TestGetJiraIssueKey(t *testing.T) {
    tests := []struct {
        name    string
        project string
        key     string
        want    string
    }{
        {
            name:    "numeric key with project",
            project: "PROJ",
            key:     "123",
            want:    "PROJ-123",
        },
        {
            name:    "non-numeric key",
            project: "PROJ",
            key:     "ABC",
            want:    "ABC",
        },
        {
            name:    "no project",
            project: "",
            key:     "123",
            want:    "123",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := cmdutil.GetJiraIssueKey(tt.project, tt.key)
            if got != tt.want {
                t.Errorf("GetJiraIssueKey() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### 4.2 Mock HTTP Client for Testing
Create a test helper for mocking HTTP responses:

```go
// internal/testutil/http.go
package testutil

import (
    "net/http"
    "net/http/httptest"
)

func NewMockServer(handler http.HandlerFunc) *httptest.Server {
    return httptest.NewServer(handler)
}

func MockJiraResponse(statusCode int, body string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(statusCode)
        w.Write([]byte(body))
    }
}
```

## 5. User Experience Enhancements

### 5.1 Progress Indicators
Enhance progress feedback:

```go
// Instead of just spinner, show percentage or steps
func InfoWithProgress(msg string, total int) *ProgressSpinner {
    // Show "Fetching issues... (1/10)"
}
```

### 5.2 Better Validation Messages
Provide clearer validation errors:

```go
func ValidateIssueKey(key string) error {
    if key == "" {
        return &ErrValidation{
            Field:   "issue-key",
            Message: "issue key cannot be empty",
        }
    }
    if !regexp.MustCompile(`^[A-Z]+-\d+$`).MatchString(key) {
        return &ErrValidation{
            Field:   "issue-key",
            Message: fmt.Sprintf("invalid format: %s (expected: PROJECT-123)", key),
        }
    }
    return nil
}
```

### 5.3 Retry Logic for Transient Errors
Add retry logic for network errors:

```go
func (c *Client) requestWithRetry(ctx context.Context, method, endpoint string, body []byte, headers Header, maxRetries int) (*http.Response, error) {
    var lastErr error
    for i := 0; i < maxRetries; i++ {
        res, err := c.request(ctx, method, endpoint, body, headers)
        if err == nil {
            return res, nil
        }
        
        // Only retry on network errors or 5xx errors
        if isRetryableError(err, res) {
            lastErr = err
            time.Sleep(time.Duration(i+1) * time.Second) // exponential backoff
            continue
        }
        
        return res, err
    }
    return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}
```

## 6. Code Organization

### 6.1 Extract Common Patterns
Create reusable command patterns:

```go
// internal/cmdcommon/base.go
type BaseCommand struct {
    requiresAuth bool
    requiresProject bool
}

func (b *BaseCommand) PreRun(cmd *cobra.Command, args []string) error {
    if b.requiresAuth {
        if err := validateAuth(); err != nil {
            return err
        }
    }
    if b.requiresProject {
        if err := validateProject(); err != nil {
            return err
        }
    }
    return nil
}
```

### 6.2 Configuration Validation
Add config validation on startup:

```go
func ValidateConfig(cfg *jira.Config) error {
    if cfg.Server == "" {
        return fmt.Errorf("server URL is required")
    }
    if !strings.HasPrefix(cfg.Server, "http") {
        return fmt.Errorf("server URL must start with http:// or https://")
    }
    // ... more validation
    return nil
}
```

## 7. Performance Improvements

### 7.1 Request Timeout Configuration
Make timeouts configurable:

```go
type Config struct {
    // ... existing fields
    Timeout time.Duration `mapstructure:"timeout"`
}

func NewClient(c Config, opts ...ClientFunc) *Client {
    timeout := c.Timeout
    if timeout == 0 {
        timeout = 30 * time.Second // default
    }
    // ... use timeout
}
```

### 7.2 Connection Pooling
Reuse HTTP client:

```go
type Client struct {
    // ... existing fields
    httpClient *http.Client
}

func (c *Client) getHTTPClient() *http.Client {
    if c.httpClient == nil {
        c.httpClient = &http.Client{
            Transport: c.transport,
            Timeout:   c.timeout,
        }
    }
    return c.httpClient
}
```

## 8. Documentation Improvements

### 8.1 Add GoDoc Comments
Ensure all exported functions have proper documentation:

```go
// GetJiraIssueKey constructs an issue key from a project and key.
// If project is empty, returns the key as-is.
// If key is numeric, returns "PROJECT-KEY", otherwise returns uppercase key.
func GetJiraIssueKey(project, key string) string {
    // ...
}
```

### 8.2 Add Examples to README
Include more real-world examples in README for common workflows.

## Priority Implementation Order

1. **High Priority:**
   - Convert `Run` to `RunE` for better error handling
   - Add structured error types
   - Improve error messages with suggestions
   - Add argument validation using Cobra's Args

2. **Medium Priority:**
   - Structured logging
   - Better progress indicators
   - Configuration validation
   - Retry logic for transient errors

3. **Low Priority:**
   - Performance optimizations
   - Enhanced testing utilities
   - Code organization refactoring

## Next Steps

1. Create a feature branch: `git checkout -b enhance/error-handling`
2. Start with high-priority items
3. Add tests for new error handling
4. Update documentation
5. Submit PR with improvements

