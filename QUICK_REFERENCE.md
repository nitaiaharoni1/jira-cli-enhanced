# Quick Reference: Using the Enhanced Error Handling

## Error Types

### Creating Errors

```go
// Authentication error
return &jira.ErrAuthentication{
    Reason: "invalid API token",
}

// Not found error
return &jira.ErrNotFound{
    Resource: "issue",
    ID:       "PROJ-123",
}

// Validation error
return &jira.ErrValidation{
    Field:   "issue-key",
    Message: "invalid format: expected PROJECT-123",
}

// Rate limit error
return &jira.ErrRateLimit{
    RetryAfter: 60, // seconds
}

// Network error
return &jira.ErrNetwork{
    Underlying: err,
}
```

### Checking Error Types

```go
// Check if error is retryable
if jira.IsRetryableError(err) {
    // Retry logic
}

// Type assertion
if e, ok := err.(*jira.ErrNotFound); ok {
    fmt.Printf("Resource: %s, ID: %s\n", e.Resource, e.ID)
}
```

## Validation Functions

```go
// Validate issue key
if err := cmdutil.ValidateIssueKey("PROJ-123"); err != nil {
    return err
}

// Validate project key
if err := cmdutil.ValidateProjectKey("PROJ"); err != nil {
    return err
}

// Validate server URL
if err := cmdutil.ValidateServerURL("https://example.atlassian.net"); err != nil {
    return err
}
```

## Command Patterns

### Using RunE

```go
var cmd = &cobra.Command{
    Use:   "command ISSUE-KEY",
    Short: "Command description",
    Args: func(cmd *cobra.Command, args []string) error {
        if len(args) == 0 {
            return fmt.Errorf("issue key is required")
        }
        return cmdutil.ValidateIssueKey(args[0])
    },
    RunE: func(cmd *cobra.Command, args []string) error {
        // Return errors instead of calling ExitIfError
        if err := doSomething(); err != nil {
            return fmt.Errorf("failed to do something: %w", err)
        }
        return nil
    },
}
```

### Error Wrapping

```go
// Good: Wrap with context
if err != nil {
    return fmt.Errorf("failed to fetch issue %q: %w", key, err)
}

// Better: Convert to structured error
if e, ok := err.(*jira.ErrUnexpectedResponse); ok {
    if e.StatusCode == 404 {
        return &jira.ErrNotFound{
            Resource: "issue",
            ID:       key,
        }
    }
    return fmt.Errorf("unexpected response: %w", err)
}
```

## Common Patterns

### Converting HTTP Errors

```go
resp, err := client.Get(ctx, path, headers)
if err != nil {
    return &jira.ErrNetwork{Underlying: err}
}

if resp.StatusCode == 404 {
    return &jira.ErrNotFound{
        Resource: "issue",
        ID:       key,
    }
}

if resp.StatusCode == 401 {
    return &jira.ErrAuthentication{
        Reason: "invalid credentials",
    }
}

if resp.StatusCode == 429 {
    retryAfter := parseRetryAfter(resp.Header)
    return &jira.ErrRateLimit{
        RetryAfter: retryAfter,
    }
}
```

### Retry Logic

```go
func doWithRetry(fn func() error, maxRetries int) error {
    var lastErr error
    for i := 0; i < maxRetries; i++ {
        err := fn()
        if err == nil {
            return nil
        }
        
        if !jira.IsRetryableError(err) {
            return err
        }
        
        lastErr = err
        time.Sleep(time.Duration(i+1) * time.Second)
    }
    return fmt.Errorf("max retries exceeded: %w", lastErr)
}
```

## Best Practices

1. **Always use RunE** for commands that can fail
2. **Wrap errors with context** - include what operation failed and relevant IDs
3. **Use structured error types** - enables better error handling
4. **Validate early** - use Args validation in Cobra
5. **Provide suggestions** - errors should tell users how to fix the issue

## Migration Checklist

- [ ] Convert `Run` to `RunE`
- [ ] Add `Args` validation
- [ ] Replace `cmdutil.ExitIfError()` with `return err`
- [ ] Wrap errors with context
- [ ] Convert HTTP errors to structured types
- [ ] Test error messages are helpful

