# Missing Operations & Functionalities Analysis

## 🔍 Comprehensive Analysis of Missing Features

Based on analysis of the codebase and Jira API capabilities, here are missing operations and functionalities:

---

## 📎 **Attachments Management** ❌

**Status:** Not implemented

**Missing Operations:**
- Upload attachments to issues
- Download attachments from issues
- Delete attachments
- List attachments for an issue
- Bulk attachment operations

**Jira API Support:** ✅ Available
- `POST /rest/api/2/issue/{issueIdOrKey}/attachments`
- `GET /rest/api/2/attachment/{id}`
- `DELETE /rest/api/2/attachment/{id}`

**Priority:** 🔴 High (Common use case)

**Example Use Cases:**
```bash
# Upload attachment
jira issue attachment upload PROJ-123 file.pdf

# List attachments
jira issue attachment list PROJ-123

# Download attachment
jira issue attachment download PROJ-123 ATTACHMENT-ID

# Delete attachment
jira issue attachment delete PROJ-123 ATTACHMENT-ID
```

---

## 👁️ **Watch/Unwatch Operations** ⚠️

**Status:** Partial (only watch exists)

**Missing Operations:**
- Unwatch issues
- List watchers
- Bulk watch/unwatch

**Jira API Support:** ✅ Available
- `DELETE /rest/api/2/issue/{issueIdOrKey}/watchers?username={username}`
- `GET /rest/api/2/issue/{issueIdOrKey}/watchers`

**Priority:** 🟡 Medium

**Example Use Cases:**
```bash
# Unwatch issue
jira issue unwatch PROJ-123

# List watchers
jira issue watchers PROJ-123

# Bulk unwatch
jira issue unwatch-bulk PROJ-1 PROJ-2 PROJ-3
```

---

## 💬 **Comment Management** ⚠️

**Status:** Partial (only add exists)

**Missing Operations:**
- List comments
- Edit comments
- Delete comments
- Get specific comment
- Internal comments management

**Jira API Support:** ✅ Available
- `GET /rest/api/2/issue/{issueIdOrKey}/comment`
- `PUT /rest/api/2/issue/{issueIdOrKey}/comment/{id}`
- `DELETE /rest/api/2/issue/{issueIdOrKey}/comment/{id}`

**Priority:** 🟡 Medium

**Example Use Cases:**
```bash
# List comments
jira issue comment list PROJ-123

# Edit comment
jira issue comment edit PROJ-123 COMMENT-ID "Updated text"

# Delete comment
jira issue comment delete PROJ-123 COMMENT-ID
```

---

## ⏱️ **Worklog Management** ⚠️

**Status:** Partial (only add exists)

**Missing Operations:**
- List worklogs
- Update worklog
- Delete worklog
- Get worklog details

**Jira API Support:** ✅ Available
- `GET /rest/api/2/issue/{issueIdOrKey}/worklog`
- `PUT /rest/api/2/issue/{issueIdOrKey}/worklog/{id}`
- `DELETE /rest/api/2/issue/{issueIdOrKey}/worklog/{id}`

**Priority:** 🟡 Medium

**Example Use Cases:**
```bash
# List worklogs
jira issue worklog list PROJ-123

# Update worklog
jira issue worklog update PROJ-123 WORKLOG-ID "2h" "Updated comment"

# Delete worklog
jira issue worklog delete PROJ-123 WORKLOG-ID
```

---

## 📊 **Issue History/Changelog** ❌

**Status:** Not implemented

**Missing Operations:**
- View issue changelog/history
- Filter changelog by field
- View specific change details

**Jira API Support:** ✅ Available
- `GET /rest/api/2/issue/{issueIdOrKey}?expand=changelog`

**Priority:** 🟡 Medium

**Example Use Cases:**
```bash
# View changelog
jira issue history PROJ-123

# Filter by field
jira issue history PROJ-123 --field status
```

---

## 👍 **Voting** ❌

**Status:** Not implemented

**Missing Operations:**
- Vote on issues
- Unvote on issues
- List voters
- Check if user voted

**Jira API Support:** ✅ Available
- `POST /rest/api/2/issue/{issueIdOrKey}/votes`
- `DELETE /rest/api/2/issue/{issueIdOrKey}/votes`

**Priority:** 🟢 Low (Less common)

**Example Use Cases:**
```bash
# Vote on issue
jira issue vote PROJ-123

# Unvote
jira issue unvote PROJ-123

# List voters
jira issue voters PROJ-123
```

---

## 🔍 **Saved Filters** ❌

**Status:** Not implemented

**Missing Operations:**
- List saved filters
- Create filter
- Update filter
- Delete filter
- Execute saved filter
- Share filter

**Jira API Support:** ✅ Available
- `GET /rest/api/2/filter`
- `POST /rest/api/2/filter`
- `PUT /rest/api/2/filter/{id}`
- `DELETE /rest/api/2/filter/{id}`

**Priority:** 🟡 Medium

**Example Use Cases:**
```bash
# List filters
jira filter list

# Create filter
jira filter create "My Filter" "project = PROJ"

# Execute filter
jira filter execute FILTER-ID

# Delete filter
jira filter delete FILTER-ID
```

---

## 📈 **Dashboards** ❌

**Status:** Not implemented

**Missing Operations:**
- List dashboards
- View dashboard
- Create dashboard
- Update dashboard
- Delete dashboard

**Jira API Support:** ✅ Available (Jira Cloud)
- `GET /rest/api/3/dashboard`
- `POST /rest/api/3/dashboard`
- `PUT /rest/api/3/dashboard/{id}`
- `DELETE /rest/api/3/dashboard/{id}`

**Priority:** 🟢 Low (Less common in CLI)

---

## 🏃 **Sprint Management** ⚠️

**Status:** Partial (list, add, close exist)

**Missing Operations:**
- Create sprint
- Update sprint
- Start sprint
- Activate sprint
- Get sprint details

**Jira API Support:** ✅ Available
- `POST /rest/agile/1.0/sprint`
- `PUT /rest/agile/1.0/sprint/{sprintId}`
- `POST /rest/agile/1.0/sprint/{sprintId}` (start)

**Priority:** 🟡 Medium

**Example Use Cases:**
```bash
# Create sprint
jira sprint create "Sprint 1" --start-date "2025-01-01" --end-date "2025-01-14"

# Start sprint
jira sprint start SPRINT-ID

# Update sprint
jira sprint update SPRINT-ID --name "Updated Sprint"
```

---

## 🏗️ **Project Management** ⚠️

**Status:** Partial (only list exists)

**Missing Operations:**
- Create project
- Update project
- Delete project
- Get project details
- Project roles/permissions

**Jira API Support:** ⚠️ Limited (Admin only)
- `POST /rest/api/2/project` (Admin)
- `PUT /rest/api/2/project/{projectIdOrKey}` (Admin)
- `DELETE /rest/api/2/project/{projectIdOrKey}` (Admin)

**Priority:** 🟢 Low (Admin operations)

---

## 📋 **Board Management** ⚠️

**Status:** Partial (only list exists)

**Missing Operations:**
- Create board
- Update board
- Delete board
- Get board configuration
- Board columns configuration

**Jira API Support:** ⚠️ Limited (Admin/Project Admin)
- `POST /rest/agile/1.0/board`
- `PUT /rest/agile/1.0/board/{boardId}`
- `DELETE /rest/agile/1.0/board/{boardId}`

**Priority:** 🟢 Low (Admin operations)

---

## 🚀 **Release/Version Management** ⚠️

**Status:** Partial (only list exists)

**Missing Operations:**
- Create release/version
- Update release/version
- Delete release/version
- Release/version details
- Merge versions

**Jira API Support:** ✅ Available
- `POST /rest/api/2/version`
- `PUT /rest/api/2/version/{id}`
- `DELETE /rest/api/2/version/{id}`

**Priority:** 🟡 Medium

**Example Use Cases:**
```bash
# Create release
jira release create "v1.0.0" --project PROJ --start-date "2025-01-01" --release-date "2025-01-15"

# Update release
jira release update VERSION-ID --name "v1.0.1"

# Delete release
jira release delete VERSION-ID
```

---

## 🔄 **Bulk Operations** ⚠️

**Status:** Partial (we added some, but missing more)

**Already Implemented:**
- ✅ Bulk assignment (`assign-bulk`)
- ✅ Bulk status transition (`move-bulk`)

**Missing Bulk Operations:**
- Bulk delete
- Bulk clone
- Bulk watch/unwatch
- Bulk link/unlink
- Bulk comment
- Bulk worklog

**Priority:** 🟡 Medium

**Example Use Cases:**
```bash
# Bulk delete
jira issue delete-bulk PROJ-1 PROJ-2 PROJ-3

# Bulk clone
jira issue clone-bulk PROJ-1 PROJ-2 PROJ-3 --project NEW-PROJ

# Bulk watch
jira issue watch-bulk PROJ-1 PROJ-2 PROJ-3
```

---

## 📝 **Issue Templates** ⚠️

**Status:** Partial (template support exists but limited)

**Missing Operations:**
- List templates
- Create template
- Update template
- Delete template
- Template variables/substitution

**Priority:** 🟢 Low

---

## 🔗 **Advanced Linking** ⚠️

**Status:** Partial (basic link/unlink exists)

**Missing Operations:**
- List all links for issue
- Get link details
- Update link type
- Link hierarchy visualization

**Priority:** 🟢 Low

---

## 📊 **Reporting & Analytics** ❌

**Status:** Not implemented

**Missing Operations:**
- Velocity reports
- Burndown charts (data export)
- Sprint reports
- Issue statistics
- Time tracking reports

**Jira API Support:** ⚠️ Limited (some via Agile API)

**Priority:** 🟢 Low (Better suited for web UI)

---

## 👥 **User Management** ❌

**Status:** Not implemented

**Missing Operations:**
- List users
- Search users (exists but not as command)
- Get user details
- User groups management

**Jira API Support:** ⚠️ Limited (Admin only)

**Priority:** 🟢 Low (Admin operations)

---

## 🔐 **Permissions & Security** ❌

**Status:** Not implemented

**Missing Operations:**
- View permissions
- Permission schemes
- Security levels

**Jira API Support:** ⚠️ Limited (Admin only)

**Priority:** 🟢 Low (Admin operations)

---

## 📋 **Summary by Priority**

### 🔴 High Priority (Common Use Cases)
1. **Attachments Management** - Upload/download/delete/list
2. **Comment Management** - List/edit/delete (beyond add)
3. **Worklog Management** - List/update/delete (beyond add)

### 🟡 Medium Priority (Useful Features)
4. **Watch/Unwatch** - Unwatch, list watchers
5. **Issue History/Changelog** - View change history
6. **Saved Filters** - Create/manage filters
7. **Sprint Creation** - Create/update sprints
8. **Release Management** - Create/update releases
9. **Bulk Operations** - More bulk commands

### 🟢 Low Priority (Less Common or Admin)
10. **Voting** - Vote/unvote operations
11. **Dashboards** - Dashboard management
12. **Project Management** - Create/update projects (admin)
13. **Board Management** - Create/update boards (admin)
14. **Templates** - Template management
15. **Advanced Linking** - Link visualization
16. **Reporting** - Analytics and reports
17. **User Management** - User operations (admin)
18. **Permissions** - Permission management (admin)

---

## 🎯 **Recommended Next Steps**

1. **Start with High Priority:**
   - Attachments management
   - Comment list/edit/delete
   - Worklog list/update/delete

2. **Then Medium Priority:**
   - Unwatch operations
   - Issue history/changelog
   - Sprint creation
   - Release management

3. **Consider User Feedback:**
   - Survey users on most needed features
   - Prioritize based on actual usage patterns

---

## 📈 **Implementation Complexity**

| Feature | Complexity | API Support | Estimated Effort |
|---------|-----------|-------------|------------------|
| Attachments | Medium | ✅ Full | 2-3 days |
| Comment Management | Low | ✅ Full | 1-2 days |
| Worklog Management | Low | ✅ Full | 1-2 days |
| Unwatch | Low | ✅ Full | 0.5 days |
| Issue History | Low | ✅ Full | 1 day |
| Sprint Creation | Medium | ✅ Full | 1-2 days |
| Release Management | Medium | ✅ Full | 1-2 days |
| Saved Filters | Medium | ✅ Full | 2-3 days |
| Bulk Operations | Low-Medium | ✅ Full | 1-2 days each |

---

This analysis provides a comprehensive view of missing operations. The high-priority items would significantly enhance the CLI's usability for daily workflows.


