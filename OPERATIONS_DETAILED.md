# Detailed Operations Reference - What's Available

## ✅ Already Implemented Operations

### 📊 Planning & Estimation

#### **Time Estimation**
```bash
# Set original estimate when creating issue
jira issue create -tStory -s"Feature" --original-estimate "2d 3h"

# Update remaining estimate when logging work
jira issue worklog add PROJ-123 "1h" --new-estimate "1d 2h"
```

**Available:**
- ✅ `--original-estimate` flag for issue creation
- ✅ `--new-estimate` flag for worklog entries
- ✅ Time format: `2d 3h 30m`, `10m`, etc.

#### **Story Points (via Custom Fields)**
```bash
# Set story points when creating issue
jira issue create -tStory -s"Feature" --custom story-points=5

# Set story points when editing issue
jira issue edit PROJ-123 --custom story-points=8
```

**Available:**
- ✅ Custom fields support via `--custom` flag
- ✅ Story points can be set as custom field
- ✅ Works for both create and edit operations

**Note:** Story points field name varies by Jira instance. You need to configure it in your config file first. See: https://github.com/ankitpokhrel/jira-cli/discussions/346

### 🔗 Linking

#### **Issue Linking**
```bash
# Link two issues
jira issue link PROJ-123 PROJ-456 Blocks

# Link types: Blocks, Relates to, Clones, etc.
jira issue link PROJ-123 PROJ-456 "Relates to"
```

**Available:**
- ✅ Link issues with various link types
- ✅ Unlink issues
- ✅ Add remote web links

### 👥 Assignees

#### **Setting Assignees**
```bash
# Assign to user
jira issue assign PROJ-123 "John Doe"

# Assign to self
jira issue assign PROJ-123 $(jira me)

# Unassign
jira issue assign PROJ-123 x

# Assign during transition
jira issue move PROJ-123 "In Progress" -a$(jira me)
```

**Available:**
- ✅ Assign/unassign issues
- ✅ Assign during issue creation
- ✅ Assign during transition
- ✅ Assign during edit

### 📍 Statuses

#### **Status Management**
```bash
# Transition issue to new status
jira issue move PROJ-123 "In Progress"

# Transition with comment and resolution
jira issue move PROJ-123 Done -RFixed --comment "Completed"

# Filter by status
jira issue list -s"In Progress"
jira issue list -s~Done  # Not Done
```

**Available:**
- ✅ Transition between statuses
- ✅ Set resolution during transition
- ✅ Add comments during transition
- ✅ Filter issues by status
- ✅ Status categories (open, closed, etc.)

### 🏃 Sprints

#### **Sprint Management**
```bash
# List sprints
jira sprint list

# List current sprint issues
jira sprint list --current

# Add issues to sprint
jira sprint add 123 PROJ-456 PROJ-789

# Close sprint
jira sprint close 123

# Filter sprint issues
jira sprint list 123 -a$(jira me) -yHigh
```

**Available:**
- ✅ List sprints
- ✅ List issues in sprint
- ✅ Add issues to sprint (up to 50 at once)
- ✅ Close sprints
- ✅ Filter by sprint state (active, closed, future)
- ✅ Filter sprint issues by all issue filters

### 📋 Epics

#### **Epic Management**
```bash
# List epics
jira epic list

# List issues in epic
jira epic list EPIC-1

# Create epic
jira epic create -n"Epic Name" -s"Summary"

# Add issues to epic
jira epic add EPIC-1 PROJ-123 PROJ-456

# Remove issues from epic
jira epic remove PROJ-123 PROJ-456
```

**Available:**
- ✅ Create epics
- ✅ List epics
- ✅ List epic issues
- ✅ Add/remove issues from epics (up to 50 at once)
- ✅ Filter epic issues by all issue filters

### 🎯 Custom Fields

#### **Setting Custom Fields**
```bash
# Create issue with custom fields
jira issue create -tStory -s"Feature" --custom story-points=5,epic-link=EPIC-1

# Edit issue custom fields
jira issue edit PROJ-123 --custom story-points=8

# Create epic with custom fields
jira epic create -n"Epic" --custom story-points=13
```

**Available:**
- ✅ Set custom fields on create
- ✅ Update custom fields on edit
- ✅ Support for various field types (number, option, array, project)
- ✅ Validation of custom fields

**Supported Custom Field Types:**
- Number fields (e.g., story points)
- Option fields (dropdowns)
- Array fields (multi-select)
- Project fields

---

## 🔍 What's Available - Complete List

### Issue Operations
| Operation | Command | Status |
|-----------|---------|--------|
| Create | `jira issue create` | ✅ |
| Edit | `jira issue edit` | ✅ |
| View | `jira issue view` | ✅ |
| List/Search | `jira issue list` | ✅ |
| Delete | `jira issue delete` | ✅ |
| Clone | `jira issue clone` | ✅ |
| Assign | `jira issue assign` | ✅ |
| Transition/Move | `jira issue move` | ✅ |
| Link | `jira issue link` | ✅ |
| Unlink | `jira issue unlink` | ✅ |
| Remote Link | `jira issue link remote` | ✅ |
| Comment | `jira issue comment add` | ✅ |
| Worklog | `jira issue worklog add` | ✅ |
| Watch | `jira issue watch` | ✅ |

### Planning Operations
| Operation | Command | Status |
|-----------|---------|--------|
| Set Original Estimate | `--original-estimate` | ✅ |
| Update Remaining Estimate | `--new-estimate` (worklog) | ✅ |
| Set Story Points | `--custom story-points=X` | ✅ |
| Set Epic | `--parent EPIC-KEY` or `--custom epic-link=EPIC-KEY` | ✅ |
| Add to Sprint | `jira sprint add` | ✅ |
| Add to Epic | `jira epic add` | ✅ |

### Status & Workflow Operations
| Operation | Command | Status |
|-----------|---------|--------|
| Transition Status | `jira issue move` | ✅ |
| Set Resolution | `-R` flag with move | ✅ |
| Filter by Status | `-s` flag | ✅ |
| Filter by Resolution | `-R` flag | ✅ |

### Assignment Operations
| Operation | Command | Status |
|-----------|---------|--------|
| Assign User | `jira issue assign` | ✅ |
| Unassign | `jira issue assign ISSUE-KEY x` | ✅ |
| Assign on Create | `-a` flag | ✅ |
| Assign on Transition | `-a` flag with move | ✅ |
| Filter by Assignee | `-a` flag | ✅ |

### Sprint Operations
| Operation | Command | Status |
|-----------|---------|--------|
| List Sprints | `jira sprint list` | ✅ |
| List Sprint Issues | `jira sprint list SPRINT-ID` | ✅ |
| Add to Sprint | `jira sprint add` | ✅ |
| Close Sprint | `jira sprint close` | ✅ |
| Filter by Sprint State | `--state` flag | ✅ |

### Epic Operations
| Operation | Command | Status |
|-----------|---------|--------|
| Create Epic | `jira epic create` | ✅ |
| List Epics | `jira epic list` | ✅ |
| List Epic Issues | `jira epic list EPIC-KEY` | ✅ |
| Add to Epic | `jira epic add` | ✅ |
| Remove from Epic | `jira epic remove` | ✅ |

### Linking Operations
| Operation | Command | Status |
|-----------|---------|--------|
| Link Issues | `jira issue link` | ✅ |
| Unlink Issues | `jira issue unlink` | ✅ |
| Add Remote Link | `jira issue link remote` | ✅ |
| View Linked Issues | `jira issue view` (shows links) | ✅ |

---

## ❌ What's NOT Available (Potential Enhancements)

### Missing Operations

1. **Bulk Status Transitions**
   - Currently: One issue at a time
   - Could add: `jira issue move ISSUE-1 ISSUE-2 ISSUE-3 "Done"`

2. **Bulk Assignment**
   - Currently: One issue at a time
   - Could add: `jira issue assign ISSUE-1 ISSUE-2 ISSUE-3 "John Doe"`

3. **Bulk Estimation**
   - Currently: Set per issue
   - Could add: `jira issue estimate ISSUE-1 ISSUE-2 ISSUE-3 "2d"`

4. **Sprint Planning Commands**
   - Currently: Manual sprint add
   - Could add: `jira sprint plan SPRINT-ID` (interactive planning)

5. **Epic Planning Commands**
   - Currently: Manual epic add
   - Could add: `jira epic plan EPIC-KEY` (interactive planning)

6. **Estimate Adjustment**
   - Currently: Only via worklog
   - Could add: `jira issue estimate PROJ-123 "3d"` (direct estimate update)

7. **Story Points Direct Command**
   - Currently: Via custom fields
   - Could add: `jira issue story-points PROJ-123 5` (dedicated command)

8. **Bulk Custom Field Updates**
   - Currently: One issue at a time
   - Could add: `jira issue custom ISSUE-1 ISSUE-2 story-points=5`

9. **Sprint Capacity Planning**
   - Currently: Not available
   - Could add: `jira sprint capacity SPRINT-ID` (show capacity vs. planned)

10. **Velocity Tracking**
    - Currently: Not available
    - Could add: `jira sprint velocity` (show team velocity)

---

## 💡 Usage Examples for Available Operations

### Complete Planning Workflow
```bash
# 1. Create issue with estimate and story points
jira issue create -tStory -s"New Feature" \
  --original-estimate "3d" \
  --custom story-points=5

# 2. Add to epic
jira epic add EPIC-1 PROJ-123

# 3. Add to sprint
jira sprint add 456 PROJ-123

# 4. Assign to developer
jira issue assign PROJ-123 "John Doe"

# 5. Transition to In Progress
jira issue move PROJ-123 "In Progress" --comment "Starting work"

# 6. Log work and update estimate
jira issue worklog add PROJ-123 "2h" --new-estimate "2d 6h"

# 7. Update story points if needed
jira issue edit PROJ-123 --custom story-points=8

# 8. Transition to Done
jira issue move PROJ-123 Done -RFixed --comment "Completed"
```

### Sprint Planning
```bash
# 1. View current sprint
jira sprint list --current

# 2. Add multiple issues to sprint
jira sprint add 123 PROJ-1 PROJ-2 PROJ-3 PROJ-4 PROJ-5

# 3. View sprint issues assigned to me
jira sprint list 123 -a$(jira me)

# 4. Close sprint when done
jira sprint close 123
```

### Epic Planning
```bash
# 1. Create epic
jira epic create -n"User Authentication" -s"Epic for auth features"

# 2. Add issues to epic
jira epic add EPIC-1 PROJ-1 PROJ-2 PROJ-3

# 3. View epic issues
jira epic list EPIC-1

# 4. Filter epic issues by status
jira epic list EPIC-1 -s"In Progress"
```

---

## 📝 Summary

### ✅ Fully Supported
- **Linking**: ✅ Issue linking, unlinking, remote links
- **Planning**: ✅ Epics, sprints, custom fields
- **Estimating**: ✅ Time estimates, story points (via custom fields)
- **Assignees**: ✅ Assign, unassign, filter by assignee
- **Statuses**: ✅ Transition, filter, set resolution
- **Sprints**: ✅ List, add issues, close, filter

### ⚠️ Partially Supported (via workarounds)
- **Story Points**: Via custom fields (requires config)
- **Bulk Operations**: Can add up to 50 issues at once to sprint/epic

### ❌ Not Supported (potential enhancements)
- Direct story points command
- Bulk status transitions
- Bulk assignments
- Sprint capacity planning
- Velocity tracking
- Interactive sprint/epic planning

---

## 🚀 Recommendations

If you need features that aren't available:

1. **Story Points**: Use `--custom story-points=X` (configure field name in config)
2. **Bulk Operations**: Use shell scripts to loop through issues
3. **Planning**: Use interactive TUI for visual planning, then add to sprint/epic

Most common operations are fully supported! The tool is quite comprehensive for day-to-day Jira management.

