# Quick Test Guide - New Operations

This guide helps you quickly test the new operations against a real Jira instance.

## Quick Start (5 minutes)

### 1. Set Up Environment

```bash
# Export your Jira credentials
export JIRA_API_TOKEN="your-api-token-here"
export JIRA_TEST_SERVER="https://your-instance.atlassian.net"
export JIRA_TEST_LOGIN="your-email@example.com"
export JIRA_TEST_PROJECT="TEST"  # Your project key
```

### 2. Build the CLI

```bash
cd /Users/nitaiaharoni/REPOS/jira-cli
go build -o bin/jira ./cmd/jira
```

### 3. Run Automated Tests

```bash
# Enable integration tests
export JIRA_INTEGRATION_TEST=true

# Run the test script
./scripts/test-integration.sh
```

### 4. Or Test Manually

```bash
# Interactive testing
./scripts/test-manual.sh
```

## Manual Testing Examples

### Test Bulk Assignment

```bash
# Create test issues first
./bin/jira issue create -tTask -s"Test 1" --no-input
./bin/jira issue create -tTask -s"Test 2" --no-input

# Get issue keys (e.g., TEST-1, TEST-2)
# Then assign both
./bin/jira issue assign-bulk TEST-1 TEST-2 $(./bin/jira me | head -1)
```

### Test Bulk Status Transition

```bash
./bin/jira issue move-bulk TEST-1 TEST-2 "In Progress"
```

### Test Estimate

```bash
# Set original estimate
./bin/jira issue estimate TEST-1 "2d 3h"

# Update remaining estimate
./bin/jira issue estimate TEST-1 "1d" --remaining
```

### Test Story Points

```bash
./bin/jira issue story-points TEST-1 5
```

### Test Bulk Custom Fields

```bash
./bin/jira issue custom TEST-1 TEST-2 story-points=8
```

## What to Check

✅ **Bulk Assignment**: All issues assigned to same user  
✅ **Bulk Transition**: All issues moved to same status  
✅ **Estimate**: Time estimate updated correctly  
✅ **Story Points**: Story points set (if field configured)  
✅ **Custom Fields**: Custom fields updated (if configured)  

## Troubleshooting

**"Authentication failed"**
- Check API token is correct
- Verify login email matches your account

**"Project not found"**
- Verify project key is correct
- Check you have access to the project

**"Story points field not found"**
- Story points is a custom field
- Configure it in your Jira config or use `--field` flag

## Cleanup

After testing, delete test issues:

```bash
./bin/jira issue delete TEST-1 --no-input
./bin/jira issue delete TEST-2 --no-input
```

Or let the integration test script handle cleanup automatically.

