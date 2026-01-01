# New Operations Implemented ✅

All missing operations have been successfully implemented!

## ✅ New Commands Added

### 1. Bulk Status Transitions
**Command:** `jira issue move-bulk`
**Aliases:** `move-batch`, `transition-bulk`

**Usage:**
```bash
# Transition multiple issues to same status
jira issue move-bulk PROJ-1 PROJ-2 PROJ-3 "In Progress"

# With comment and resolution
jira issue move-bulk PROJ-1 PROJ-2 PROJ-3 Done --comment "All completed" -RFixed

# With assignee change
jira issue move-bulk PROJ-1 PROJ-2 "In Progress" -a$(jira me)
```

**Features:**
- Transition up to 50 issues at once
- Supports all transition flags (comment, resolution, assignee)
- Shows success/failure summary
- Continues on partial failures

**File:** `internal/cmd/issue/move/bulk.go`

---

### 2. Bulk Assignment
**Command:** `jira issue assign-bulk`
**Aliases:** `assign-batch`, `asg-bulk`

**Usage:**
```bash
# Assign multiple issues to a user
jira issue assign-bulk PROJ-1 PROJ-2 PROJ-3 "John Doe"

# Assign to self
jira issue assign-bulk PROJ-1 PROJ-2 PROJ-3 $(jira me)

# Unassign multiple issues
jira issue assign-bulk PROJ-1 PROJ-2 PROJ-3 x
```

**Features:**
- Assign up to 50 issues at once
- Supports unassign (x) and default assignee
- Auto-searches for users
- Shows success/failure summary

**File:** `internal/cmd/issue/assign/bulk.go`

---

### 3. Direct Estimate Command
**Command:** `jira issue estimate`
**Aliases:** `est`

**Usage:**
```bash
# Set original estimate
jira issue estimate PROJ-123 "2d 3h"

# Update remaining estimate
jira issue estimate PROJ-123 "1d" --remaining

# Set estimate for multiple issues
jira issue estimate PROJ-123 PROJ-456 PROJ-789 "3d"
```

**Features:**
- Set original estimate (via edit API)
- Update remaining estimate (via worklog API)
- Supports multiple issues
- Time format: `2d 3h 30m`, `10m`, `1w`, etc.

**File:** `internal/cmd/issue/estimate/estimate.go`

**Implementation Notes:**
- Original estimate: Uses edit API with `timetracking.originalEstimate` field
- Remaining estimate: Uses worklog API with `adjustEstimate=new&newEstimate=...`
- Added `OriginalEstimate` field to `EditRequest` struct
- Added `TimeTracking` support to edit API

---

### 4. Direct Story Points Command
**Command:** `jira issue story-points`
**Aliases:** `sp`, `points`

**Usage:**
```bash
# Set story points for an issue
jira issue story-points PROJ-123 5

# Set story points for multiple issues
jira issue story-points PROJ-123 PROJ-456 PROJ-789 8

# Remove story points (set to 0)
jira issue story-points PROJ-123 0

# Use specific field name
jira issue story-points PROJ-123 5 --field "Story Point Estimate"
```

**Features:**
- Auto-detects story points field from config
- Supports custom field name via `--field` flag
- Supports multiple issues
- Validates field configuration

**File:** `internal/cmd/issue/storypoints/storypoints.go`

**Implementation Notes:**
- Searches configured custom fields for story points
- Looks for keywords: "story point", "storypoint", "story-point"
- Falls back to manual field specification if not found

---

### 5. Bulk Custom Field Updates
**Command:** `jira issue custom`
**Aliases:** `cf`

**Usage:**
```bash
# Set a custom field for an issue
jira issue custom PROJ-123 story-points=5

# Set multiple custom fields
jira issue custom PROJ-123 story-points=5,epic-link=EPIC-1

# Set custom fields for multiple issues
jira issue custom PROJ-123 PROJ-456 PROJ-789 story-points=8
```

**Features:**
- Update multiple custom fields at once
- Supports multiple issues
- Validates custom fields against configuration
- Format: `FIELD=VALUE` or `FIELD1=VALUE1,FIELD2=VALUE2`

**File:** `internal/cmd/issue/custom/custom.go`

---

## 📊 Summary

### Files Created
1. `internal/cmd/issue/move/bulk.go` - Bulk status transitions
2. `internal/cmd/issue/assign/bulk.go` - Bulk assignment
3. `internal/cmd/issue/estimate/estimate.go` - Direct estimate command
4. `internal/cmd/issue/storypoints/storypoints.go` - Direct story points command
5. `internal/cmd/issue/custom/custom.go` - Bulk custom field updates

### Files Modified
1. `internal/cmd/issue/issue.go` - Registered new commands
2. `pkg/jira/edit.go` - Added `OriginalEstimate` support to `EditRequest`

### New Operations Available

| Operation | Command | Status |
|-----------|---------|--------|
| Bulk Status Transitions | `jira issue move-bulk` | ✅ |
| Bulk Assignment | `jira issue assign-bulk` | ✅ |
| Direct Estimate | `jira issue estimate` | ✅ |
| Direct Story Points | `jira issue story-points` | ✅ |
| Bulk Custom Fields | `jira issue custom` | ✅ |

---

## 🎯 Usage Examples

### Complete Planning Workflow (New!)
```bash
# 1. Create issues with estimates
jira issue create -tStory -s"Feature 1" --original-estimate "2d"
jira issue create -tStory -s"Feature 2" --original-estimate "3d"

# 2. Set story points for multiple issues
jira issue story-points PROJ-1 PROJ-2 5

# 3. Add to epic
jira epic add EPIC-1 PROJ-1 PROJ-2

# 4. Add to sprint
jira sprint add 123 PROJ-1 PROJ-2

# 5. Assign all to developer
jira issue assign-bulk PROJ-1 PROJ-2 "John Doe"

# 6. Transition all to In Progress
jira issue move-bulk PROJ-1 PROJ-2 "In Progress" --comment "Starting work"

# 7. Update estimates if needed
jira issue estimate PROJ-1 "3d"  # Update original estimate
jira issue estimate PROJ-1 "2d" --remaining  # Update remaining

# 8. Transition all to Done
jira issue move-bulk PROJ-1 PROJ-2 Done -RFixed --comment "Completed"
```

### Sprint Planning (Enhanced!)
```bash
# 1. View current sprint
jira sprint list --current

# 2. Add multiple issues to sprint
jira sprint add 123 PROJ-1 PROJ-2 PROJ-3 PROJ-4 PROJ-5

# 3. Set story points for all sprint issues
jira issue story-points PROJ-1 PROJ-2 PROJ-3 PROJ-4 PROJ-5 8

# 4. Assign all to team members
jira issue assign-bulk PROJ-1 PROJ-2 PROJ-3 "Developer A"
jira issue assign-bulk PROJ-4 PROJ-5 "Developer B"

# 5. Transition all to In Progress
jira issue move-bulk PROJ-1 PROJ-2 PROJ-3 PROJ-4 PROJ-5 "In Progress"
```

### Bulk Operations
```bash
# Bulk status transition
jira issue move-bulk PROJ-1 PROJ-2 PROJ-3 "Code Review"

# Bulk assignment
jira issue assign-bulk PROJ-1 PROJ-2 PROJ-3 $(jira me)

# Bulk estimate update
jira issue estimate PROJ-1 PROJ-2 PROJ-3 "2d"

# Bulk story points
jira issue story-points PROJ-1 PROJ-2 PROJ-3 5

# Bulk custom fields
jira issue custom PROJ-1 PROJ-2 PROJ-3 story-points=5,epic-link=EPIC-1
```

---

## 🔧 Technical Details

### Bulk Operations Implementation
- All bulk operations process issues sequentially
- Continue on errors (partial success handling)
- Show summary of successes and failures
- Limit: Up to 50 issues (Jira API limit)

### Estimate Implementation
- **Original Estimate**: Uses `PUT /issue/{key}` with `fields.timetracking.originalEstimate`
- **Remaining Estimate**: Uses `POST /issue/{key}/worklog` with `adjustEstimate=new&newEstimate=...`
- Added `OriginalEstimate` field to `EditRequest` struct
- Added `TimeTracking` support to edit request structure

### Story Points Implementation
- Auto-detects story points field from configuration
- Searches for common field name patterns
- Falls back to manual specification via `--field` flag
- Validates field exists and is configured

### Error Handling
- All commands use `RunE` for proper error handling
- Partial success reporting for bulk operations
- Clear error messages with suggestions
- Validation before processing

---

## ✅ Testing

All code compiles successfully:
```bash
go build ./...
# ✅ No errors
```

Commands are registered and available:
```bash
jira issue --help
# Shows: move-bulk, assign-bulk, estimate, story-points, custom
```

---

## 📝 Notes

1. **Story Points Field**: Must be configured in your Jira config file. The command will auto-detect common field names, but you can override with `--field`.

2. **Original Estimate**: Uses the standard Jira API `timetracking` field. Works for both Cloud and On-Premise installations.

3. **Bulk Operations**: Process issues sequentially to avoid overwhelming the API. Consider rate limits for very large batches.

4. **Error Handling**: Bulk operations continue on errors and report both successes and failures at the end.

---

## 🎉 Result

All missing operations are now implemented! The jira-cli tool now supports:

✅ Bulk status transitions  
✅ Bulk assignment  
✅ Direct estimate commands  
✅ Direct story points commands  
✅ Bulk custom field updates  

The tool is now even more powerful for managing Jira issues efficiently!


