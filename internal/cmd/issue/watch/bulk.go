package watch

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
	bulkHelpText = `Bulk watch adds multiple issues to watchers.

You can watch up to 50 issues at once. If no watcher is specified, adds the current user.`
	bulkExamples = `# Watch multiple issues (add self)
$ jira issue watch-bulk PROJ-1 PROJ-2 PROJ-3

# Watch multiple issues for a specific user
$ jira issue watch-bulk PROJ-1 PROJ-2 PROJ-3 "John Doe"

# Watch issues from stdin
$ jira issue list --keys-only | jira issue watch-bulk

# Watch issues from JQL
$ jira issue watch-bulk --jql "status = 'To Do'" "John Doe"`
)

// NewCmdWatchBulk is a bulk watch command.
func NewCmdWatchBulk() *cobra.Command {
	cmd := cobra.Command{
		Use:     "watch-bulk ISSUE-KEY... [WATCHER]",
		Short:   "Add multiple issues to watchers",
		Long:    bulkHelpText,
		Example: bulkExamples,
		Aliases: []string{"watch-batch"},
		Args:    cobra.MinimumNArgs(0),
		RunE:    watchBulk,
	}

	cmd.Flags().Bool("stdin", false, "Read issue keys from stdin (one per line)")
	cmd.Flags().String("jql", "", "Apply to all issues matching JQL query")

	return &cmd
}

func watchBulk(cmd *cobra.Command, args []string) error {
	project := viper.GetString("project.key")
	debug, _ := cmd.Flags().GetBool("debug")
	stdin, _ := cmd.Flags().GetBool("stdin")
	jql, _ := cmd.Flags().GetString("jql")
	client := api.DefaultClient(debug)

	var issueKeys []string
	var watcher string

	// Get watcher from args (optional, defaults to current user)
	if len(args) > 0 && !strings.HasPrefix(args[len(args)-1], "PROJ-") && !strings.HasPrefix(args[len(args)-1], "PROJECT-") {
		watcher = args[len(args)-1]
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

	if watcher != "" {
		users, err := api.ProxyUserSearch(client, &jira.UserSearchOptions{
			Query:      watcher,
			Project:    project,
			MaxResults: 100,
		})
		if err != nil {
			return fmt.Errorf("failed to search for user: %w", err)
		}
		if len(users) == 0 {
			return fmt.Errorf("user %q not found", watcher)
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

	s := cmdutil.Info(fmt.Sprintf("Adding %q as watcher to %d issues...", uname, len(issueKeys)))
	defer s.Stop()

	var failed []string
	var succeeded []string

	for _, key := range issueKeys {
		err := api.ProxyWatchIssue(client, key, userObj)
		if err != nil {
			failed = append(failed, key)
			continue
		}
		succeeded = append(succeeded, key)
	}

	s.Stop()

	if len(failed) > 0 {
		if len(succeeded) > 0 {
			cmdutil.Warn("Added watcher to %d issues successfully, %d failed", len(succeeded), len(failed))
			fmt.Printf("Failed: %s\n", strings.Join(failed, ", "))
		} else {
			return fmt.Errorf("failed to add watcher to all issues")
		}
	} else {
		cmdutil.Success("Successfully added %q as watcher to %d issues", uname, len(succeeded))
	}

	return nil
}

