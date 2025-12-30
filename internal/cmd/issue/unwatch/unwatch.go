package unwatch

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

const (
	helpText = `Unwatch removes user from issue watchers.`
	examples = `# Unwatch issue (remove self)
$ jira issue unwatch PROJ-123

# Unwatch specific user
$ jira issue unwatch PROJ-123 "John Doe"`
)

// NewCmdUnwatch is an unwatch command.
func NewCmdUnwatch() *cobra.Command {
	return &cobra.Command{
		Use:     "unwatch ISSUE-KEY [USER]",
		Short:   "Remove user from issue watchers",
		Long:    helpText,
		Example: examples,
		Aliases: []string{"unwat"},
		RunE:    unwatch,
	}
}

func unwatch(cmd *cobra.Command, args []string) error {
	project := viper.GetString("project.key")
	var issueKey, user string

	if len(args) >= 1 {
		issueKey = cmdutil.GetJiraIssueKey(project, args[0])
	}
	if len(args) >= 2 {
		user = args[1]
	}

	debug, _ := cmd.Flags().GetBool("debug")
	client := api.DefaultClient(debug)

	if issueKey == "" {
		var ans string
		qs := &survey.Question{
			Name:     "key",
			Prompt:   &survey.Input{Message: "Issue key"},
			Validate: survey.Required,
		}
		if err := survey.Ask([]*survey.Question{qs}, &ans); err != nil {
			return err
		}
		issueKey = cmdutil.GetJiraIssueKey(project, ans)
	}

	// Get user object if user specified, otherwise use self
	var userObj *jira.User
	var uname string

	if user != "" {
		// Search for user
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
		// Use self - need to get AccountID for Cloud instances
		me, err := client.Me()
		if err != nil {
			return fmt.Errorf("failed to get current user: %w", err)
		}
		
		// Search for self to get AccountID (needed for Cloud)
		users, err := api.ProxyUserSearch(client, &jira.UserSearchOptions{
			Query:      me.Login,
			Project:    project,
			MaxResults: 10,
		})
		if err == nil && len(users) > 0 {
			// Found user with AccountID
			userObj = users[0]
			uname = getQueryableName(userObj.Name, userObj.DisplayName)
		} else {
			// Fallback: create User object from Me (may not have AccountID)
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

	s := cmdutil.Info(fmt.Sprintf("Removing %q from watchers of issue %q...", uname, issueKey))
	err := api.ProxyUnwatchIssue(client, issueKey, userObj)
	s.Stop()

	if err != nil {
		return err
	}

	cmdutil.Success("User %q removed from watchers of issue %q", uname, issueKey)
	fmt.Printf("%s\n", cmdutil.GenerateServerBrowseURL(viper.GetString("server"), issueKey))

	return nil
}

func getQueryableName(name, displayName string) string {
	if name != "" {
		return name
	}
	return displayName
}

