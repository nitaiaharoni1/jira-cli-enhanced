# Integration Testing Guide

This guide explains how to test jira-cli operations against a real Jira API instance.

## Overview

Integration tests verify that the CLI works correctly with real Jira instances. They create actual issues, perform operations, and clean up afterward.

## Prerequisites

1. **Jira Instance**: Access to a Jira Cloud or Server instance
2. **API Token**: Generate an API token from your Jira profile
3. **Test Project**: A project key where test issues can be created

## Setup

### 1. Generate API Token

**For Jira Cloud:**
1. Go to https://id.atlassian.com/manage-profile/security/api-tokens
2. Click "Create API token"
3. Copy the token

**For Jira Server:**
- Use your password for basic auth, or
- Generate a Personal Access Token from your profile

### 2. Set Environment Variables

```bash
export JIRA_INTEGRATION_TEST=true
export JIRA_TEST_SERVER="https://your-instance.atlassian.net"
export JIRA_TEST_LOGIN="your-email@example.com"
export JIRA_TEST_API_TOKEN="your-api-token"
export JIRA_TEST_PROJECT="TEST"  # Your project key
export JIRA_TEST_DEBUG="true"   # Optional: for debug output
```

### 3. Run Integration Tests

```bash
# Make script executable
chmod +x scripts/test-integration.sh

# Run integration tests
./scripts/test-integration.sh
```

## What Gets Tested

The integration test script tests:

1. **Issue Creation** - Creates 3 test issues
2. **Bulk Assignment** - Assigns all issues to a user
3. **Bulk Status Transition** - Transitions issues to new status
4. **Estimate Update** - Updates time estimate
5. **Story Points** - Sets story points (if configured)
6. **Custom Fields** - Updates custom fields (if configured)

## Manual Testing

You can also test individual commands manually:

### Test Bulk Operations

```bash
# Set up environment
export JIRA_API_TOKEN="your-token"
export JIRA_CONFIG_FILE="$HOME/.config/.jira/.config.yml"

# Create test issues
jira issue create -tTask -s"Test 1" --no-input
jira issue create -tTask -s"Test 2" --no-input
jira issue create -tTask -s"Test 3" --no-input

# Test bulk assignment
jira issue assign-bulk PROJ-1 PROJ-2 PROJ-3 $(jira me)

# Test bulk transition
jira issue move-bulk PROJ-1 PROJ-2 PROJ-3 "In Progress"

# Test estimate
jira issue estimate PROJ-1 "2d"

# Test story points
jira issue story-points PROJ-1 5

# Cleanup
jira issue delete PROJ-1 --no-input
jira issue delete PROJ-2 --no-input
jira issue delete PROJ-3 --no-input
```

## Test Script Features

- **Automatic Cleanup**: Test issues are deleted after tests
- **Error Handling**: Continues on non-critical failures
- **Debug Output**: Set `JIRA_TEST_DEBUG=true` for verbose output
- **Safe**: Only runs when explicitly enabled

## Running Go Integration Tests

You can also run Go integration tests:

```bash
# Build with integration tag
go test -tags=integration ./internal/cmd/issue/move/...

# Or run all integration tests
go test -tags=integration ./...
```

## Test Data

The integration tests create temporary issues that are automatically cleaned up. Make sure:

1. You have permission to create/delete issues in the test project
2. The project has the necessary issue types (Task, Story, etc.)
3. Workflows allow the transitions being tested

## Troubleshooting

### "Authentication failed"
- Check your API token is correct
- Verify login email matches your Jira account
- For Server: ensure auth type is correct (basic/bearer)

### "Project not found"
- Verify project key is correct
- Check you have access to the project

### "Status transition failed"
- The target status may not be available
- Check your workflow configuration
- Try a different status

### "Story points field not found"
- Story points is a custom field
- Configure it in your Jira config file first
- Or use `--field` flag to specify field name

## Safety

- Tests create real issues (but clean them up)
- Use a test project, not production
- Tests are idempotent (can run multiple times)
- Cleanup runs even if tests fail

## CI/CD Integration

To run integration tests in CI:

```yaml
# .github/workflows/integration.yml
name: Integration Tests

on:
  workflow_dispatch:
  schedule:
    - cron: '0 0 * * 0'  # Weekly

jobs:
  integration:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '^1.24.1'
      - name: Run integration tests
        env:
          JIRA_INTEGRATION_TEST: true
          JIRA_TEST_SERVER: ${{ secrets.JIRA_TEST_SERVER }}
          JIRA_TEST_LOGIN: ${{ secrets.JIRA_TEST_LOGIN }}
          JIRA_TEST_API_TOKEN: ${{ secrets.JIRA_TEST_API_TOKEN }}
          JIRA_TEST_PROJECT: ${{ secrets.JIRA_TEST_PROJECT }}
        run: ./scripts/test-integration.sh
```

## Best Practices

1. **Use Test Project**: Never test against production
2. **Isolated Tests**: Each test should be independent
3. **Cleanup**: Always clean up test data
4. **Idempotent**: Tests should be safe to run multiple times
5. **Documentation**: Document any manual setup required


