# Additional Improvements for Jira CLI

Based on deeper code analysis, here are more improvements beyond error handling:

## 1. Performance Optimizations

### 1.1 HTTP Client Reuse
**Current Issue:** HTTP client is created on every request in `pkg/jira/client.go:288`

```go
// Current (inefficient)
httpClient := &http.Client{Transport: c.transport}
return httpClient.Do(req.WithContext(ctx))
```

**Fix:** Reuse HTTP client instance

```go
// pkg/jira/client.go
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

func (c *Client) request(ctx context.Context, method, endpoint string, body []byte, headers Header) (*http.Response, error) {
    // ... existing code ...
    return c.getHTTPClient().Do(req.WithContext(ctx))
}
```

**Benefits:**
- Connection pooling and reuse
- Better performance for multiple requests
- Reduced memory allocation

### 1.2 Configurable Timeout
**Current Issue:** Timeout is hardcoded to 15 seconds in `api/client.go:14`

**Fix:** Make timeout configurable via config file

```go
// api/client.go
func getClientTimeout() time.Duration {
    timeout := viper.GetDuration("timeout")
    if timeout == 0 {
        timeout = 15 * time.Second // default
    }
    return timeout
}

func Client(config jira.Config) *jira.Client {
    // ... existing code ...
    jiraClient = jira.NewClient(
        config,
        jira.WithTimeout(getClientTimeout()),
        jira.WithInsecureTLS(*config.Insecure),
    )
    return jiraClient
}
```

### 1.3 Request Batching
For operations that can be batched (e.g., adding multiple issues to sprint):

```go
// pkg/jira/batch.go
type BatchRequest struct {
    Requests []Request
}

func (c *Client) Batch(ctx context.Context, batch BatchRequest) ([]Response, error) {
    // Implement batching logic
}
```

## 2. Security Improvements

### 2.1 Input Sanitization
Add input validation and sanitization:

```go
// internal/cmdutil/sanitize.go
package cmdutil

import (
    "html"
    "strings"
)

// SanitizeString removes potentially dangerous characters
func SanitizeString(input string) string {
    // Remove null bytes
    input = strings.ReplaceAll(input, "\x00", "")
    // Escape HTML
    input = html.EscapeString(input)
    return strings.TrimSpace(input)
}

// ValidateJQL prevents JQL injection
func ValidateJQL(jql string) error {
    // Check for dangerous patterns
    dangerous := []string{"';", "--", "/*", "*/", "xp_", "sp_"}
    lowerJQL := strings.ToLower(jql)
    for _, pattern := range dangerous {
        if strings.Contains(lowerJQL, pattern) {
            return &jira.ErrValidation{
                Field:   "jql",
                Message: "potentially dangerous JQL pattern detected",
            }
        }
    }
    return nil
}
```

### 2.2 Credential Handling
Improve credential security:

```go
// pkg/jira/credentials.go
package jira

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "io"
)

// SecureTokenStorage encrypts tokens before storage
type SecureTokenStorage struct {
    key []byte
}

func (s *SecureTokenStorage) Encrypt(token string) ([]byte, error) {
    // Implement AES encryption
}

func (s *SecureTokenStorage) Decrypt(encrypted []byte) (string, error) {
    // Implement AES decryption
}
```

### 2.3 Fix TODO: Fail on Invalid Custom Fields
**Current:** Warning only (line 243 in `internal/cmdcommon/create.go`)

**Fix:** Make it fail with proper error

```go
// internal/cmdcommon/create.go
func ValidateCustomFields(fields map[string]string, configuredFields []jira.IssueTypeField) error {
    if len(fields) == 0 {
        return nil
    }

    fieldsMap := make(map[string]string)
    for _, configured := range configuredFields {
        identifier := strings.ReplaceAll(strings.ToLower(strings.TrimSpace(configured.Name)), " ", "-")
        fieldsMap[identifier] = configured.Name
    }

    invalidCustomFields := make([]string, 0, len(fields))
    for key := range fields {
        if _, ok := fieldsMap[key]; !ok {
            invalidCustomFields = append(invalidCustomFields, key)
        }
    }

    if len(invalidCustomFields) > 0 {
        return &jira.ErrValidation{
            Field:   "custom_fields",
            Message: fmt.Sprintf("invalid custom fields: %s", strings.Join(invalidCustomFields, ", ")),
        }
    }
    
    return nil
}
```

## 3. Code Quality Improvements

### 3.1 Reduce Cyclomatic Complexity
**Current:** Functions marked with `//nolint:gocyclo` (e.g., `config/generator.go:118`)

**Fix:** Break down complex functions:

```go
// Instead of one large Generate() function, split into:
func (c *JiraCLIConfigGenerator) Generate() (string, error) {
    cfgFile, err := c.determineConfigFile()
    if err != nil {
        return "", err
    }

    if err := c.validateAndPrepare(); err != nil {
        return "", err
    }

    if err := c.configureInstallation(); err != nil {
        return "", err
    }

    if err := c.configureAuthentication(); err != nil {
        return "", err
    }

    if err := c.configureProjectAndBoard(); err != nil {
        return "", err
    }

    return c.writeConfig(cfgFile)
}
```

### 3.2 Add Structured Logging
Create a proper logger package:

```go
// internal/logger/logger.go
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

### 3.3 Configuration Validation
Add config validation on startup:

```go
// internal/config/validator.go
package config

import (
    "fmt"
    "net/url"
    "strings"

    "github.com/ankitpokhrel/jira-cli/pkg/jira"
)

func ValidateConfig() error {
    server := viper.GetString("server")
    if server == "" {
        return &jira.ErrValidation{
            Field:   "server",
            Message: "server URL is required",
        }
    }

    if err := validateServerURL(server); err != nil {
        return err
    }

    login := viper.GetString("login")
    if login == "" {
        return &jira.ErrValidation{
            Field:   "login",
            Message: "login is required",
        }
    }

    // Check if API token exists (via env, netrc, or keyring)
    if !hasAPIToken() {
        return &jira.ErrAuthentication{
            Reason: "API token not found",
        }
    }

    return nil
}

func validateServerURL(serverURL string) error {
    u, err := url.Parse(serverURL)
    if err != nil {
        return &jira.ErrValidation{
            Field:   "server",
            Message: fmt.Sprintf("invalid URL: %s", err),
        }
    }

    if u.Scheme != "http" && u.Scheme != "https" {
        return &jira.ErrValidation{
            Field:   "server",
            Message: "URL must use http:// or https://",
        }
    }

    return nil
}

func hasAPIToken() bool {
    // Check env var
    if os.Getenv("JIRA_API_TOKEN") != "" {
        return true
    }
    // Check netrc
    // Check keyring
    return false
}
```

## 4. Testing Improvements

### 4.1 Add Test Coverage Reporting
Update CI to include coverage:

```yaml
# .github/workflows/ci.yml
- name: Run tests with coverage
  run: |
    go test -race -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html

- name: Upload coverage
  uses: codecov/codecov-action@v3
  with:
    file: ./coverage.out
```

### 4.2 Add Integration Tests
Create integration test helpers:

```go
// internal/testutil/integration.go
package testutil

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/ankitpokhrel/jira-cli/pkg/jira"
)

func NewTestServer(handler http.HandlerFunc) *httptest.Server {
    return httptest.NewServer(handler)
}

func NewTestClient(server *httptest.Server) *jira.Client {
    return jira.NewClient(jira.Config{
        Server:   server.URL,
        Login:    "test",
        APIToken: "test-token",
    })
}
```

### 4.3 Add Table-Driven Tests
For utility functions:

```go
// internal/cmdutil/validation_test.go
func TestValidateIssueKey(t *testing.T) {
    tests := []struct {
        name    string
        key     string
        wantErr bool
    }{
        {
            name:    "valid issue key",
            key:     "PROJ-123",
            wantErr: false,
        },
        {
            name:    "invalid format",
            key:     "invalid",
            wantErr: true,
        },
        {
            name:    "empty key",
            key:     "",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateIssueKey(tt.key)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateIssueKey() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## 5. CI/CD Enhancements

### 5.1 Enhanced CI Workflow
Add more checks:

```yaml
# .github/workflows/ci.yml
name: CI

on:
  pull_request:
    types: [opened, synchronize, reopened]
  push:
    branches: [main]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '^1.24.1'
      - run: make lint

  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '^1.24.1'
      - run: make deps
      - run: make test
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.out

  security:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run Gosec Security Scanner
        uses: securego/gosec@master
        with:
          args: './...'

  build:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        os: [linux, darwin, windows]
        arch: [amd64, arm64]
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '^1.24.1'
      - name: Build
        run: |
          GOOS=${{ matrix.os }} GOARCH=${{ matrix.arch }} go build -o jira-${{ matrix.os }}-${{ matrix.arch }} ./cmd/jira
```

### 5.2 Add Pre-commit Hooks
Create `.githooks/pre-commit`:

```bash
#!/bin/bash
set -e

echo "Running linter..."
make lint

echo "Running tests..."
make test

echo "Checking for TODO/FIXME comments..."
if git diff --cached --name-only | xargs grep -l "TODO\|FIXME"; then
    echo "Warning: Found TODO/FIXME comments in staged files"
fi
```

## 6. Documentation Improvements

### 6.1 Add GoDoc Comments
Ensure all exported functions have documentation:

```go
// GetJiraIssueKey constructs an issue key from a project and key.
//
// If project is empty, returns the key as-is.
// If key is numeric, returns "PROJECT-KEY", otherwise returns uppercase key.
//
// Example:
//   GetJiraIssueKey("PROJ", "123")  // Returns "PROJ-123"
//   GetJiraIssueKey("PROJ", "ABC")  // Returns "ABC"
//   GetJiraIssueKey("", "PROJ-123") // Returns "PROJ-123"
func GetJiraIssueKey(project, key string) string {
    // ...
}
```

### 6.2 Add Examples to Commands
Add more examples to command help:

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
  jira issue create -tBug -s"Bug title" -yHigh
  
  # View an issue
  jira issue view PROJ-123
  
  # Transition an issue
  jira issue move PROJ-123 "In Progress"`,
}
```

## 7. User Experience Enhancements

### 7.1 Better Progress Indicators
Show progress for long operations:

```go
// internal/cmdutil/progress.go
type ProgressBar struct {
    total   int
    current int
    message string
}

func NewProgressBar(total int, message string) *ProgressBar {
    return &ProgressBar{
        total:   total,
        current: 0,
        message: message,
    }
}

func (p *ProgressBar) Update(current int) {
    p.current = current
    percentage := (float64(current) / float64(p.total)) * 100
    fmt.Fprintf(os.Stderr, "\r%s: %d/%d (%.1f%%)", p.message, current, p.total, percentage)
}

func (p *ProgressBar) Finish() {
    fmt.Fprintf(os.Stderr, "\r%s: Complete\n", p.message)
}
```

### 7.2 Add Retry Logic
Implement retry for transient errors:

```go
// pkg/jira/retry.go
func (c *Client) requestWithRetry(ctx context.Context, method, endpoint string, body []byte, headers Header, maxRetries int) (*http.Response, error) {
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        res, err := c.request(ctx, method, endpoint, body, headers)
        if err == nil {
            // Check for retryable status codes
            if res.StatusCode < 500 {
                return res, nil
            }
        }
        
        if jira.IsRetryableError(err) || (res != nil && res.StatusCode >= 500) {
            lastErr = err
            if i < maxRetries-1 {
                backoff := time.Duration(i+1) * time.Second
                time.Sleep(backoff)
                continue
            }
        }
        
        return res, err
    }
    
    return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}
```

### 7.3 Add Autocomplete Improvements
Enhance shell completion:

```go
// internal/cmd/completion/completion.go
func generateBashCompletion() string {
    return `# Custom completions for issue keys
_jira_issue_key_completion() {
    local cur="${COMP_WORDS[COMP_CWORD]}"
    # Fetch recent issue keys from Jira
    COMPREPLY=($(jira issue list --plain --columns key --no-headers 2>/dev/null | grep "^$cur"))
}
complete -F _jira_issue_key_completion jira issue view
complete -F _jira_issue_key_completion jira issue edit
`
}
```

## 8. Dependency Management

### 8.1 Check for Updates
Add script to check for outdated dependencies:

```bash
#!/bin/bash
# scripts/check-deps.sh

echo "Checking for outdated dependencies..."
go list -u -m all | grep -v "^go " | grep "\[" || echo "All dependencies are up to date"
```

### 8.2 Add Dependabot
Create `.github/dependabot.yml`:

```yaml
version: 2
updates:
  - package-ecosystem: "gomod"
    directory: "/"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 10
```

## Priority Order

1. **High Priority:**
   - HTTP client reuse (performance)
   - Configuration validation (security)
   - Fix TODO for custom fields (code quality)
   - Add test coverage reporting (testing)

2. **Medium Priority:**
   - Structured logging
   - Reduce cyclomatic complexity
   - Add retry logic
   - Enhanced CI workflow

3. **Low Priority:**
   - Request batching
   - Enhanced autocomplete
   - Documentation improvements
   - Dependency management automation

## Implementation Notes

- Start with high-priority items
- Test each change thoroughly
- Update documentation as you go
- Consider backward compatibility
- Add tests for new features


