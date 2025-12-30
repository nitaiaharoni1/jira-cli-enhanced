# Unwatch Fix - Implementation Complete ✅

## Issue Fixed

The unwatch command was failing because:
1. Cloud instances require `accountId` parameter, not `username`
2. The code was using `username` parameter even for Cloud instances

## Solution Implemented

1. **Created `ProxyUnwatchIssue` function** - Similar to `ProxyWatchIssue`, handles Cloud vs Local
2. **Added `unwatchIssueWithAccountID` method** - Uses correct parameter based on installation type
3. **Updated unwatch command** - Now uses `ProxyUnwatchIssue` to get proper user object with AccountID

## Changes Made

### `api/client.go`
- Added `ProxyUnwatchIssue` function that:
  - Uses `AccountID` for Cloud instances
  - Uses `Name` (username) for Local instances

### `pkg/jira/issue.go`
- Added `unwatchIssueWithAccountID` method that:
  - Uses `accountId` query parameter for Cloud
  - Uses `username` query parameter for Local

### `internal/cmd/issue/unwatch/unwatch.go`
- Updated to use `ProxyUnwatchIssue` instead of direct `UnwatchIssue` call
- Searches for user to get AccountID for Cloud instances
- Handles both self and specified user cases

## Testing

The fix is implemented and ready for testing. The command should now work correctly for Cloud instances by using the `accountId` parameter instead of `username`.

