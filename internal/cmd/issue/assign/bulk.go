package assign

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

const (
	bulkHelpText = `Bulk assign assigns multiple issues to a user.

You can assign up to 50 issues at once. All issues will be assigned to the same user.`
	bulkExamples = `# Assign multiple issues to a user
$ jira issue assign-bulk PROJ-1 PROJ-2 PROJ-3 "John Doe"

# Assign multiple issues to self
$ jira issue assign-bulk PROJ-1 PROJ-2 PROJ-3 $(jira me)

# Unassign multiple issues
$ jira issue assign-bulk PROJ-1 PROJ-2 PROJ-3 x`
)

// NewCmdAssignBulk is a bulk assign command.
func NewCmdAssignBulk() *cobra.Command {
	return &cobra.Command{
		Use:     "assign-bulk ISSUE-KEY... ASSIGNEE",
		Short:   "Assign multiple issues to a user",
		Long:    bulkHelpText,
		Example: bulkExamples,
		Aliases: []string{"assign-batch", "asg-bulk"},
		Args:    cobra.MinimumNArgs(2),
		RunE:    assignBulk,
	}
}

func assignBulk(cmd *cobra.Command, args []string) error {
	project := viper.GetString("project.key")

	// Last argument is the assignee
	assignee := args[len(args)-1]
	issueKeys := args[:len(args)-1]

	// Normalize issue keys
	normalizedKeys := make([]string, 0, len(issueKeys))
	for _, key := range issueKeys {
		normalizedKeys = append(normalizedKeys, cmdutil.GetJiraIssueKey(project, key))
	}

	debug, _ := cmd.Flags().GetBool("debug")
	client := api.DefaultClient(debug)

	// Handle special cases
	lu := strings.ToLower(assignee)
	var user *jira.User
	var assigneeValue string

	switch {
	case lu == "x" || lu == strings.ToLower(optionNone):
		assigneeValue = jira.AssigneeNone
	case lu == strings.ToLower(optionDefault):
		assigneeValue = jira.AssigneeDefault
	default:
		// Search for user
		users, err := api.ProxyUserSearch(client, &jira.UserSearchOptions{
			Query:      assignee,
			Project:    project,
			MaxResults: maxResults,
		})
		if err != nil {
			return fmt.Errorf("failed to search for user: %w", err)
		}

		if len(users) == 0 {
			return fmt.Errorf("user %q not found", assignee)
		}

		// Find exact match
		for _, u := range users {
			name := strings.ToLower(getQueryableName(u.Name, u.DisplayName))
			if name == lu || strings.ToLower(u.Email) == lu {
				user = u
				break
			}
		}

		if user == nil {
			// Use first result if no exact match
			user = users[0]
		}
	}

	// Assign all issues
	var assigneeName string
	if assigneeValue == jira.AssigneeNone {
		assigneeName = "unassigned"
	} else if assigneeValue == jira.AssigneeDefault {
		assigneeName = "default assignee"
	} else {
		assigneeName = getQueryableName(user.Name, user.DisplayName)
	}

	s := cmdutil.Info(fmt.Sprintf("Assigning %d issues to %q...", len(normalizedKeys), assigneeName))
	defer s.Stop()

	var failed []string
	var succeeded []string

	for _, key := range normalizedKeys {
		err := api.ProxyAssignIssue(client, key, user, assigneeValue)
		if err != nil {
			failed = append(failed, key)
			continue
		}
		succeeded = append(succeeded, key)
	}

	s.Stop()

	if len(failed) > 0 {
		if len(succeeded) > 0 {
			cmdutil.Warn("Assigned %d issues successfully, %d failed", len(succeeded), len(failed))
			fmt.Printf("Failed: %s\n", strings.Join(failed, ", "))
		} else {
			return fmt.Errorf("failed to assign all issues")
		}
	} else {
		if assigneeValue == jira.AssigneeNone {
			cmdutil.Success("Successfully unassigned %d issues", len(succeeded))
		} else {
			cmdutil.Success("Successfully assigned %d issues to %q", len(succeeded), assigneeName)
		}
	}

	return nil
}


