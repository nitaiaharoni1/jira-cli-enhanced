package assign

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/internal/query"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

const (
	helpText = `Assign issue to a user.`
	examples = `$ jira issue assign ISSUE-1 jon@domain.tld

# Assignee name or email needs to be an exact match
$ jira issue assign ISSUE-1 "Jon Doe"

# Assign to self
$ jira issue assign ISSUE-1 $(jira me)

# Assign to default assignee
$ jira issue assign ISSUE-1 default

# Unassign
$ jira issue assign ISSUE-1 x`

	maxResults = 100
	lineBreak  = "----------"

	optionSearch  = "[Search...]"
	optionDefault = "Default"
	optionNone    = "No-one (Unassign)"
	optionCancel  = "Cancel"
)

// NewCmdAssign is an assign command.
func NewCmdAssign() *cobra.Command {
	cmd := cobra.Command{
		Use:     "assign ISSUE-KEY... ASSIGNEE",
		Short:   "Assign issue to a user",
		Long:    helpText,
		Example: examples,
		Aliases: []string{"asg"},
		Annotations: map[string]string{
			"help:args": `ISSUE-KEY	Issue key(s), eg: ISSUE-1
ASSIGNEE	Email or display name of the user to assign the issue to`,
		},
		RunE: assign,
	}

	cmd.Flags().Bool("stdin", false, "Read issue keys from stdin (one per line)")
	cmd.Flags().String("jql", "", "Apply to all issues matching JQL query")

	return &cmd
}

func assign(cmd *cobra.Command, args []string) error {
	project := viper.GetString("project.key")
	
	// Check for stdin or JQL flags
	stdin, _ := cmd.Flags().GetBool("stdin")
	jql, _ := cmd.Flags().GetString("jql")
	
	if stdin || jql != "" {
		// Use the bulk assign command handler
		return assignBulk(cmd, args)
	}
	
	params := parseArgsAndFlags(cmd.Flags(), args, project)
	client := api.DefaultClient(params.debug)
	ac := assignCmd{
		client: client,
		users:  nil,
		params: params,
	}
	lu := strings.ToLower(ac.params.user)

	if err := ac.setIssueKey(project); err != nil {
		return err
	}

	if lu != strings.ToLower(optionNone) && lu != "x" && lu != jira.AssigneeDefault {
		if err := ac.setAvailableUsers(project); err != nil {
			return err
		}
		if err := ac.setAssignee(project); err != nil {
			return err
		}

		lu = strings.ToLower(ac.params.user)
	}

	u, err := ac.verifyAssignee()
	if err != nil {
		return fmt.Errorf("error: %s", err.Error())
	}

	var assignee, uname string

	switch {
	case u != nil:
		uname = getQueryableName(u.DisplayName, u.Name)
	case lu == strings.ToLower(optionNone) || lu == "x":
		assignee = jira.AssigneeNone
		uname = "unassigned"
	case lu == strings.ToLower(optionDefault):
		assignee = jira.AssigneeDefault
		uname = assignee
	}

	err = func() error {
		var s *spinner.Spinner
		if uname == "unassigned" {
			s = cmdutil.Info(fmt.Sprintf("Unassigning user from issue %q...", ac.params.key))
		} else {
			s = cmdutil.Info(fmt.Sprintf("Assigning issue %q to user %q...", ac.params.key, uname))
		}
		defer s.Stop()

		return api.ProxyAssignIssue(client, ac.params.key, u, assignee)
	}()
	if err != nil {
		return err
	}

	if uname == "unassigned" {
		cmdutil.Success("User unassigned from the issue %q", ac.params.key)
	} else {
		cmdutil.Success("User %q assigned to issue %q", uname, ac.params.key)
	}
	fmt.Printf("%s\n", cmdutil.GenerateServerBrowseURL(viper.GetString("server"), ac.params.key))
	return nil
}

func assignBulkIssues(client *jira.Client, issueKeys []string, assignee string, project string) error {
	// Use the bulk assign logic from bulk.go
	normalizedKeys := make([]string, 0, len(issueKeys))
	for _, key := range issueKeys {
		normalizedKeys = append(normalizedKeys, cmdutil.GetJiraIssueKey(project, key))
	}
	
	lu := strings.ToLower(assignee)
	var user *jira.User
	var assigneeValue string

	switch {
	case lu == "x" || lu == strings.ToLower(optionNone):
		assigneeValue = jira.AssigneeNone
	case lu == strings.ToLower(optionDefault):
		assigneeValue = jira.AssigneeDefault
	default:
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

		for _, u := range users {
			name := strings.ToLower(getQueryableName(u.Name, u.DisplayName))
			if name == lu || strings.ToLower(u.Email) == lu {
				user = u
				break
			}
		}

		if user == nil {
			user = users[0]
		}
	}

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

type assignParams struct {
	key   string
	user  string
	debug bool
}

func parseArgsAndFlags(flags query.FlagParser, args []string, project string) *assignParams {
	var key, user string

	nargs := len(args)
	if nargs >= 1 {
		key = cmdutil.GetJiraIssueKey(project, args[0])
	}
	if nargs >= 2 {
		user = args[1]
	}

	debug, err := flags.GetBool("debug")
	cmdutil.ExitIfError(err)

	return &assignParams{
		key:   key,
		user:  user,
		debug: debug,
	}
}

type assignCmd struct {
	client *jira.Client
	users  []*jira.User
	params *assignParams
}

func (ac *assignCmd) setIssueKey(project string) error {
	if ac.params.key != "" {
		return nil
	}

	var ans string

	qs := &survey.Question{
		Name:     "key",
		Prompt:   &survey.Input{Message: "Issue key"},
		Validate: survey.Required,
	}
	if err := survey.Ask([]*survey.Question{qs}, &ans); err != nil {
		return err
	}
	ac.params.key = cmdutil.GetJiraIssueKey(project, ans)

	return nil
}

func (ac *assignCmd) setAssignee(project string) error {
	if ac.params.user != "" && len(ac.users) == 1 {
		ac.params.user = getQueryableName(ac.users[0].Name, ac.users[0].DisplayName)
		return nil
	}

	lu := strings.ToLower(ac.params.user)
	if lu == strings.ToLower(optionNone) || lu == strings.ToLower(optionDefault) || lu == "x" {
		return nil
	}

	var (
		ans  string
		last bool
	)
	if ac.params.user != "" && len(ac.users) > 0 {
		last = true
	}

	for {
		qs := &survey.Question{
			Name: "user",
			Prompt: &survey.Select{
				Message: "Assign to user:",
				Help:    "Can't find the user? Select search and look for a keyword or cancel to abort",
				Options: ac.getOptions(last),
			},
			Validate: func(val interface{}) error {
				errInvalidSelection := fmt.Errorf("invalid selection")

				ans, ok := val.(core.OptionAnswer)
				if !ok {
					return errInvalidSelection
				}
				if ans.Value == "" || ans.Value == lineBreak {
					return errInvalidSelection
				}

				return nil
			},
		}

		if err := survey.Ask([]*survey.Question{qs}, &ans); err != nil {
			return err
		}
		if ans == optionCancel {
			return fmt.Errorf("action aborted")
		}
		if ans != optionSearch {
			break
		}
		if err := ac.getSearchKeyword(); err != nil {
			return err
		}
		if err := ac.searchAndAssignUser(project); err != nil {
			return err
		}
		last = true
	}
	ac.params.user = ans

	return nil
}

func (ac *assignCmd) getOptions(last bool) []string {
	var validUsers []string

	for _, t := range ac.users {
		if t.Active {
			name := t.DisplayName
			if t.Name != "" {
				name += fmt.Sprintf(" (%s)", t.Name)
			}
			validUsers = append(validUsers, name)
		}
	}
	always := []string{optionDefault, optionNone, optionCancel}
	options := []string{optionSearch}

	if last {
		options = append(options, validUsers...)
		options = append(options, lineBreak)
		options = append(options, always...)
	} else {
		options = append(options, always...)
		options = append(options, lineBreak)
		options = append(options, validUsers...)
	}

	return options
}

func (ac *assignCmd) getSearchKeyword() error {
	qs := &survey.Question{
		Name: "user",
		Prompt: &survey.Input{
			Message: "Search user:",
			Help:    "Type user email or display name to search for a user",
		},
		Validate: func(val interface{}) error {
			errInvalidKeyword := fmt.Errorf("enter atleast 3 characters to search")

			str, ok := val.(string)
			if !ok {
				return errInvalidKeyword
			}
			if len(str) < 3 {
				return errInvalidKeyword
			}

			return nil
		},
	}
	return survey.Ask([]*survey.Question{qs}, &ac.params.user)
}

func (ac *assignCmd) searchAndAssignUser(project string) error {
	u, err := api.ProxyUserSearch(ac.client, &jira.UserSearchOptions{
		Query:      ac.params.user,
		Project:    project,
		MaxResults: maxResults,
	})
	if err != nil {
		return err
	}
	ac.users = u
	return nil
}

func (ac *assignCmd) setAvailableUsers(project string) error {
	s := cmdutil.Info("Fetching available users. Please wait...")
	defer s.Stop()

	return ac.searchAndAssignUser(project)
}

func (ac *assignCmd) verifyAssignee() (*jira.User, error) {
	assignee := strings.ToLower(ac.params.user)
	if assignee == strings.ToLower(optionDefault) || assignee == strings.ToLower(optionNone) || assignee == "x" {
		return nil, nil
	}

	var user *jira.User

	for _, u := range ac.users {
		name := strings.ToLower(getQueryableName(u.Name, u.DisplayName))
		if name == assignee || strings.ToLower(u.Email) == assignee {
			user = u
		}
		if strings.ToLower(fmt.Sprintf("%s (%s)", u.DisplayName, u.Name)) == assignee {
			user = u
		}
		if user != nil {
			break
		}
	}

	if user == nil {
		return nil, fmt.Errorf("invalid assignee %q", ac.params.user)
	}
	if !user.Active {
		return nil, fmt.Errorf("user %q is not active", getQueryableName(user.Name, user.DisplayName))
	}
	return user, nil
}

func getQueryableName(name, displayName string) string {
	if name != "" {
		return name
	}
	return displayName
}
