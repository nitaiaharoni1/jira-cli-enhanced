package unwatch

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

const (
	bulkHelpText = `Bulk unwatch removes multiple issues from watchers.

You can unwatch up to 50 issues at once. If no user is specified, removes the current user.`
	bulkExamples = `# Unwatch multiple issues (remove self)
$ jira issue unwatch-bulk PROJ-1 PROJ-2 PROJ-3

# Unwatch multiple issues for a specific user
$ jira issue unwatch-bulk PROJ-1 PROJ-2 PROJ-3 "John Doe"

# Unwatch issues from stdin
$ jira issue list --keys-only | jira issue unwatch-bulk

# Unwatch issues from JQL
$ jira issue unwatch-bulk --jql "status = 'Done'"`
)

// NewCmdUnwatchBulk is a bulk unwatch command.
func NewCmdUnwatchBulk() *cobra.Command {
	cmd := cobra.Command{
		Use:     "unwatch-bulk ISSUE-KEY... [USER]",
		Short:   "Remove multiple issues from watchers",
		Long:    bulkHelpText,
		Example: bulkExamples,
		Aliases: []string{"unwatch-batch"},
		Args:    cobra.MinimumNArgs(0),
		RunE:    unwatchBulk,
	}

	cmd.Flags().Bool("stdin", false, "Read issue keys from stdin (one per line)")
	cmd.Flags().String("jql", "", "Apply to all issues matching JQL query")

	return &cmd
}

func unwatchBulk(cmd *cobra.Command, args []string) error {
	project := viper.GetString("project.key")
	debug, _ := cmd.Flags().GetBool("debug")
	stdin, _ := cmd.Flags().GetBool("stdin")
	jql, _ := cmd.Flags().GetString("jql")
	client := api.DefaultClient(debug)

	var issueKeys []string
	var user string

	// Get user from args (optional, defaults to current user)
	if len(args) > 0 && !strings.HasPrefix(args[len(args)-1], "PROJ-") && !strings.HasPrefix(args[len(args)-1], "PROJECT-") {
		user = args[len(args)-1]
		args = args[:len(args)-1]
	}

	// Get issue keys
	if stdin {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			key := strings.TrimSpace(scanner.Text())
			if key != "" {
				issueKeys = append(issueKeys, cmdutil.GetJiraIssueKey(project, key))
			}
		}
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("failed to read from stdin: %w", err)
		}
	} else if jql != "" {
		result, err := api.ProxySearch(client, jql, 0, 1000)
		if err != nil {
			return fmt.Errorf("failed to search issues: %w", err)
		}
		for _, issue := range result.Issues {
			issueKeys = append(issueKeys, issue.Key)
		}
	} else {
		if len(args) == 0 {
			return fmt.Errorf("no issue keys provided")
		}
		for _, key := range args {
			issueKeys = append(issueKeys, cmdutil.GetJiraIssueKey(project, key))
		}
	}

	if len(issueKeys) == 0 {
		return fmt.Errorf("no issues found")
	}

	// Get user object
	var userObj *jira.User
	var uname string

	if user != "" {
		users, err := api.ProxyUserSearch(client, &jira.UserSearchOptions{
			Query:      user,
			Project:    project,
			MaxResults: 100,
		})
		if err != nil {
			return fmt.Errorf("failed to search for user: %w", err)
		}
		if len(users) == 0 {
			return fmt.Errorf("user %q not found", user)
		}
		userObj = users[0]
		uname = getQueryableName(userObj.Name, userObj.DisplayName)
	} else {
		// Use current user
		me, err := api.ProxyMe(client)
		if err != nil {
			return fmt.Errorf("failed to get current user: %w", err)
		}
		users, err := api.ProxyUserSearch(client, &jira.UserSearchOptions{
			Query:      me.Login,
			Project:    project,
			MaxResults: 10,
		})
		if err == nil && len(users) > 0 {
			userObj = users[0]
			uname = getQueryableName(userObj.Name, userObj.DisplayName)
		} else {
			userObj = &jira.User{
				Name:        me.Login,
				DisplayName: me.Name,
				Email:       me.Email,
			}
			uname = me.Name
			if uname == "" {
				uname = me.Login
			}
		}
	}

	s := cmdutil.Info(fmt.Sprintf("Removing %q from watchers of %d issues...", uname, len(issueKeys)))
	defer s.Stop()

	var failed []string
	var succeeded []string

	for _, key := range issueKeys {
		err := api.ProxyUnwatchIssue(client, key, userObj)
		if err != nil {
			failed = append(failed, key)
			continue
		}
		succeeded = append(succeeded, key)
	}

	s.Stop()

	if len(failed) > 0 {
		if len(succeeded) > 0 {
			cmdutil.Warn("Removed watcher from %d issues successfully, %d failed", len(succeeded), len(failed))
			fmt.Printf("Failed: %s\n", strings.Join(failed, ", "))
		} else {
			return fmt.Errorf("failed to remove watcher from all issues")
		}
	} else {
		cmdutil.Success("Successfully removed %q from watchers of %d issues", uname, len(succeeded))
	}

	return nil
}

