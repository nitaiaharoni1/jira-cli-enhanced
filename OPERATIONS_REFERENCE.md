# Jira CLI Operations Reference

Complete list of all available operations in jira-cli.

## 📋 Main Commands

### 🔹 Issue Operations

#### List Issues
```bash
jira issue list [flags]
```
**Operations:**
- List/search issues with advanced filtering
- Filter by: status, priority, assignee, reporter, labels, components, dates, etc.
- Output formats: interactive table, plain text, CSV, JSON
- Supports JQL queries

**Key Flags:**
- `-s, --status` - Filter by status
- `-y, --priority` - Filter by priority
- `-a, --assignee` - Filter by assignee
- `-r, --reporter` - Filter by reporter
- `-l, --label` - Filter by label
- `-t, --type` - Filter by issue type
- `--created` - Filter by creation date
- `--updated` - Filter by update date
- `-q, --jql` - Raw JQL query
- `--plain` - Plain text output
- `--csv` - CSV output
- `--raw` - JSON output

#### Create Issue
```bash
jira issue create [flags]
```
**Operations:**
- Create new issues interactively or non-interactively
- Support for custom fields
- Markdown support for descriptions
- Template support

**Key Flags:**
- `-t, --type` - Issue type (Bug, Story, Task, etc.)
- `-s, --summary` - Issue summary
- `-b, --body` - Issue description
- `-y, --priority` - Priority
- `-l, --label` - Labels
- `-C, --component` - Components
- `-P, --parent` - Parent epic
- `--template` - Load from template file
- `--custom` - Custom fields
- `--no-input` - Skip interactive prompts

#### Edit Issue
```bash
jira issue edit ISSUE-KEY [flags]
```
**Operations:**
- Update issue fields
- Add/remove labels, components, fix versions
- Update summary, description, priority, etc.

**Key Flags:**
- `-s, --summary` - Update summary
- `-b, --body` - Update description
- `-y, --priority` - Update priority
- `--label` - Add/remove labels (use `-` prefix to remove)
- `--component` - Add/remove components
- `--fix-version` - Add/remove fix versions
- `--no-input` - Skip interactive prompts

#### View Issue
```bash
jira issue view ISSUE-KEY [flags]
```
**Operations:**
- Display issue details in terminal
- Show description, comments, linked issues
- Markdown rendering

**Key Flags:**
- `--comments N` - Show N comments
- `--plain` - Plain text mode
- `--raw` - Raw JSON output

#### Assign Issue
```bash
jira issue assign ISSUE-KEY [USER]
```
**Operations:**
- Assign issue to user
- Unassign issue
- Assign to self

**Examples:**
- `jira issue assign PROJ-123 "John Doe"`
- `jira issue assign PROJ-123 $(jira me)` - Assign to self
- `jira issue assign PROJ-123 x` - Unassign

#### Move/Transition Issue
```bash
jira issue move ISSUE-KEY [STATUS] [flags]
```
**Operations:**
- Transition issue between statuses
- Add comments during transition
- Set resolution

**Key Flags:**
- `--comment` - Add comment during transition
- `-R, --resolution` - Set resolution
- `-a, --assignee` - Change assignee during transition

#### Clone Issue
```bash
jira issue clone ISSUE-KEY [flags]
```
**Operations:**
- Clone existing issue
- Modify fields during clone
- Replace text in summary/description

**Key Flags:**
- `-s, --summary` - New summary
- `-y, --priority` - New priority
- `-a, --assignee` - New assignee
- `-H, --replace` - Replace text (format: "find:replace")

#### Delete Issue
```bash
jira issue delete ISSUE-KEY [flags]
```
**Operations:**
- Delete issue
- Cascade delete (delete subtasks)

**Key Flags:**
- `--cascade` - Delete with subtasks

#### Link Issues
```bash
jira issue link ISSUE-1 ISSUE-2 [LINK-TYPE]
```
**Operations:**
- Link two issues
- Support for various link types (Blocks, Relates to, etc.)

#### Unlink Issues
```bash
jira issue unlink ISSUE-1 ISSUE-2
```
**Operations:**
- Remove link between issues

#### Add Remote Link
```bash
jira issue link remote ISSUE-KEY URL [TEXT]
```
**Operations:**
- Add external web link to issue

#### Comment Operations
```bash
jira issue comment add ISSUE-KEY [BODY] [flags]
```
**Operations:**
- Add comments to issues
- Internal comments
- Markdown support
- Template support

**Key Flags:**
- `--internal` - Create internal comment
- `--template` - Load from template

#### Worklog Operations
```bash
jira issue worklog add ISSUE-KEY TIME [flags]
```
**Operations:**
- Add worklog entries
- Time tracking

**Key Flags:**
- `--comment` - Add comment with worklog
- `--no-input` - Skip prompts

**Time Format:**
- `2d 3h 30m` - 2 days, 3 hours, 30 minutes
- `10m` - 10 minutes

#### Watch Issue
```bash
jira issue watch ISSUE-KEY [USER]
```
**Operations:**
- Add watcher to issue
- Watch as current user

---

### 🔹 Epic Operations

#### List Epics
```bash
jira epic list [EPIC-KEY] [flags]
```
**Operations:**
- List epics
- List issues in epic
- Filter epic issues

**Key Flags:**
- `--table` - Table view instead of explorer
- All issue list filters apply

#### Create Epic
```bash
jira epic create [flags]
```
**Operations:**
- Create new epic
- Similar to issue creation

**Key Flags:**
- `-n, --name` - Epic name
- `-s, --summary` - Summary
- `-b, --body` - Description
- Other flags same as issue create

#### Add Issues to Epic
```bash
jira epic add EPIC-KEY ISSUE-1 [ISSUE-2 ...]
```
**Operations:**
- Add up to 50 issues to epic

#### Remove Issues from Epic
```bash
jira epic remove ISSUE-1 [ISSUE-2 ...]
```
**Operations:**
- Remove up to 50 issues from epic

---

### 🔹 Sprint Operations

#### List Sprints
```bash
jira sprint list [SPRINT-ID] [flags]
```
**Operations:**
- List sprints
- List issues in sprint
- Filter by sprint state

**Key Flags:**
- `--current` - Current active sprint
- `--prev` - Previous sprint
- `--next` - Next sprint
- `--state` - Filter by state (active, closed, future)
- `--table` - Table view
- All issue list filters apply

#### Add Issues to Sprint
```bash
jira sprint add SPRINT-ID ISSUE-1 [ISSUE-2 ...]
```
**Operations:**
- Add up to 50 issues to sprint

#### Close Sprint
```bash
jira sprint close SPRINT-ID
```
**Operations:**
- Close a sprint

---

### 🔹 Board Operations

#### List Boards
```bash
jira board list [flags]
```
**Operations:**
- List boards in project

---

### 🔹 Project Operations

#### List Projects
```bash
jira project list
```
**Operations:**
- List all accessible projects

---

### 🔹 Release Operations

#### List Releases
```bash
jira release list [flags]
```
**Operations:**
- List project versions/releases

**Key Flags:**
- `--project` - Filter by project

---

## 🛠️ Utility Commands

### Init
```bash
jira init [flags]
```
**Operations:**
- Initialize configuration
- Interactive setup wizard
- Configure server, credentials, project, board

**Key Flags:**
- `--force` - Overwrite existing config

### Me
```bash
jira me
```
**Operations:**
- Display current user information
- Returns account ID/username

### Open
```bash
jira open [ISSUE-KEY]
```
**Operations:**
- Open issue in browser
- Open project in browser (if no issue key)

### Server Info
```bash
jira serverinfo
```
**Operations:**
- Display Jira server information
- Version, deployment type, etc.

### Completion
```bash
jira completion [bash|zsh]
```
**Operations:**
- Generate shell completion scripts

### Version
```bash
jira version
```
**Operations:**
- Display version information

### Help
```bash
jira help [COMMAND]
```
**Operations:**
- Display help for commands

### Man Pages
```bash
jira man
```
**Operations:**
- Generate man pages

---

## 🎯 Common Operations Summary

### Issue Management
- ✅ Create, edit, delete, clone issues
- ✅ List/search with advanced filtering
- ✅ View issue details
- ✅ Assign/unassign
- ✅ Transition between statuses
- ✅ Link/unlink issues
- ✅ Add comments
- ✅ Add worklogs
- ✅ Watch issues

### Epic Management
- ✅ Create epics
- ✅ List epics and epic issues
- ✅ Add/remove issues from epics

### Sprint Management
- ✅ List sprints and sprint issues
- ✅ Add issues to sprints
- ✅ Close sprints

### Project Management
- ✅ List projects
- ✅ List boards
- ✅ List releases/versions

### Configuration
- ✅ Initialize configuration
- ✅ View current user
- ✅ Open in browser
- ✅ View server info

---

## 📊 Output Formats

### Interactive Mode (Default)
- TUI (Terminal User Interface)
- Navigate with arrow keys
- View, edit, transition from UI

### Plain Mode
- `--plain` flag
- Tab-separated values
- Suitable for scripting

### CSV Mode
- `--csv` flag
- Comma-separated values
- Import to spreadsheets

### JSON Mode
- `--raw` flag
- Raw API responses
- Programmatic processing

---

## 🔍 Filtering Capabilities

### Status Filtering
- `-s"To Do"` - Exact match
- `-s~Done` - Not operator
- `-sopen` - Status category

### Date Filtering
- `--created week` - Last week
- `--created -7d` - Last 7 days
- `--created month` - This month
- `--created-before -24w` - Before 24 weeks
- `--updated -30m` - Updated in last 30 minutes

### User Filtering
- `-a$(jira me)` - Assigned to me
- `-ax` - Unassigned
- `-a~x` - Assigned (not unassigned)
- `-r$(jira me)` - Reported by me

### Priority Filtering
- `-yHigh` - High priority
- `-y~Low` - Not low priority

### Label Filtering
- `-lbackend` - Has label "backend"
- `-lbackend -lurgent` - Has both labels

### Component Filtering
- `-CBackend` - Has component "Backend"

### Type Filtering
- `-tBug` - Bug type
- `-t~Task` - Not Task type

### JQL Queries
- `-q"summary ~ cli"` - Custom JQL
- `-q"project IS NOT EMPTY"` - All projects

---

## 💡 Usage Examples

### Daily Workflow
```bash
# Check my assigned issues
jira issue list -a$(jira me)

# View an issue
jira issue view PROJ-123

# Transition to In Progress
jira issue move PROJ-123 "In Progress" --comment "Starting work"

# Add worklog
jira issue worklog add PROJ-123 "2h" --comment "Fixed bug"

# Add comment
jira issue comment add PROJ-123 "Fixed in commit abc123"

# Transition to Done
jira issue move PROJ-123 Done -RFixed
```

### Sprint Management
```bash
# View current sprint issues
jira sprint list --current -a$(jira me)

# Add issue to sprint
jira sprint add 123 PROJ-456 PROJ-789

# Close sprint
jira sprint close 123
```

### Epic Management
```bash
# List epic issues
jira epic list EPIC-1

# Add issues to epic
jira epic add EPIC-1 PROJ-123 PROJ-456
```

### Bulk Operations
```bash
# List high priority bugs assigned to me
jira issue list -a$(jira me) -tBug -yHigh

# List issues created this week
jira issue list --created week

# List issues I'm watching
jira issue list -w
```

---

## 🎨 Interactive UI Features

### Navigation
- Arrow keys or `j/k/h/l` - Navigate
- `g` - Go to top
- `G` - Go to bottom
- `CTRL+f` - Page down
- `CTRL+b` - Page up
- `v` - View issue details
- `m` - Transition issue
- `ENTER` - Open in browser
- `c` - Copy URL to clipboard
- `CTRL+k` - Copy issue key
- `w/TAB` - Toggle sidebar focus
- `q/ESC/CTRL+c` - Quit
- `?` - Help

---

## 📝 Notes

- All commands support `--help` for detailed information
- Use `--debug` flag for verbose output
- Configuration file: `~/.config/.jira/.config.yml`
- Can override config with `--config` flag or `JIRA_CONFIG_FILE` env var
- Supports both Cloud and On-Premise Jira installations


