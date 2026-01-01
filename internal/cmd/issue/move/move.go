package move

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/internal/query"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

const (
	helpText = `Move transitions an issue from one state to another.`
	examples = `$ jira issue move ISSUE-1 "In Progress"
$ jira issue move ISSUE-1 Done`

	optionCancel = "Cancel"
)

// NewCmdMove is a move command.
func NewCmdMove() *cobra.Command {
	cmd := cobra.Command{
		Use:     "move ISSUE-KEY... STATE",
		Short:   "Transition an issue to a given state",
		Long:    helpText,
		Example: examples,
		Aliases: []string{"transition", "mv"},
		Annotations: map[string]string{
			"help:args": `ISSUE-KEY	Issue key(s), eg: ISSUE-1
STATE		State you want to transition the issue to`,
		},
		RunE: move,
	}

	cmd.Flags().SortFlags = false

	cmd.Flags().String("comment", "", "Add comment to the issue")
	cmd.Flags().StringP("assignee", "a", "", "Assign issue to a user")
	cmd.Flags().StringP("resolution", "R", "", "Set resolution")
	cmd.Flags().Bool("web", false, "Open issue in web browser after successful transition")
	cmd.Flags().Bool("stdin", false, "Read issue keys from stdin (one per line)")
	cmd.Flags().String("jql", "", "Apply to all issues matching JQL query")

	return &cmd
}

func move(cmd *cobra.Command, args []string) error {
	project := viper.GetString("project.key")
	installation := viper.GetString("installation")
	
	// Check for stdin or JQL flags
	stdin, _ := cmd.Flags().GetBool("stdin")
	jql, _ := cmd.Flags().GetString("jql")
	
	if stdin || jql != "" {
		return moveBulk(cmd, args)
	}
	
	params := parseArgsAndFlags(cmd.Flags(), args, project)
	client := api.DefaultClient(params.debug)
	mc := moveCmd{
		client:      client,
		transitions: nil,
		params:      params,
	}

	if err := mc.setIssueKey(project); err != nil {
		return err
	}
	if err := mc.setAvailableTransitions(); err != nil {
		return err
	}
	if err := mc.setDesiredState(installation); err != nil {
		return err
	}

	if mc.params.state == optionCancel {
		cmdutil.Fail("Action aborted")
		return fmt.Errorf("action aborted")
	}

	tr, err := mc.verifyTransition(installation)
	if err != nil {
		fmt.Println()
		return fmt.Errorf("error: %s", err.Error())
	}

	err = func() error {
		s := cmdutil.Info(fmt.Sprintf("Transitioning issue to %q...", tr.Name))
		defer s.Stop()

		trFieldsReq := jira.TransitionRequestFields{}
		trUpdateReq := jira.TransitionRequestUpdate{}

		if mc.params.assignee != "" {
			trFieldsReq.Assignee = &struct {
				Name string `json:"name"`
			}{Name: mc.params.assignee}
		}
		if mc.params.resolution != "" {
			trFieldsReq.Resolution = &struct {
				Name string `json:"name"`
			}{Name: mc.params.resolution}
		}
		if mc.params.comment != "" {
			trUpdateReq.Comment = []struct {
				Add struct {
					Body string `json:"body"`
				} `json:"add"`
			}{
				{Add: struct {
					Body string `json:"body"`
				}{Body: mc.params.comment}},
			}
		}

		_, err := client.Transition(mc.params.key, &jira.TransitionRequest{
			Fields: &trFieldsReq,
			Update: &trUpdateReq,
			Transition: &jira.TransitionRequestData{
				ID:   tr.ID.String(),
				Name: tr.Name,
			},
		})
		return err
	}()
	if err != nil {
		return err
	}

	server := viper.GetString("server")

	cmdutil.Success("Issue transitioned to state %q", tr.Name)
	fmt.Printf("%s\n", cmdutil.GenerateServerBrowseURL(server, mc.params.key))

	if web, _ := cmd.Flags().GetBool("web"); web {
		if err := cmdutil.Navigate(server, mc.params.key); err != nil {
			return err
		}
	}
	return nil
}

func moveBulkIssues(cmd *cobra.Command, client *jira.Client, issueKeys []string, state string, project string, installation string) error {
	// Use the bulk move logic from bulk.go
	normalizedKeys := make([]string, 0, len(issueKeys))
	for _, key := range issueKeys {
		normalizedKeys = append(normalizedKeys, cmdutil.GetJiraIssueKey(project, key))
	}

	comment, _ := cmd.Flags().GetString("comment")
	assignee, _ := cmd.Flags().GetString("assignee")
	resolution, _ := cmd.Flags().GetString("resolution")

	// Get transitions for first issue to validate state
	transitions, err := api.ProxyTransitions(client, normalizedKeys[0])
	if err != nil {
		return fmt.Errorf("failed to fetch transitions: %w", err)
	}

	var targetTransition *jira.Transition
	stateLower := strings.ToLower(state)
	for _, t := range transitions {
		if strings.ToLower(t.Name) == stateLower {
			targetTransition = t
			break
		}
	}

	if targetTransition == nil {
		available := make([]string, 0, len(transitions))
		for _, t := range transitions {
			available = append(available, fmt.Sprintf("'%s'", t.Name))
		}
		return fmt.Errorf("invalid transition state %q\nAvailable states: %s", state, strings.Join(available, ", "))
	}

	// Prepare transition request
	trFieldsReq := jira.TransitionRequestFields{}
	trUpdateReq := jira.TransitionRequestUpdate{}

	if assignee != "" {
		trFieldsReq.Assignee = &struct {
			Name string `json:"name"`
		}{Name: assignee}
	}
	if resolution != "" {
		trFieldsReq.Resolution = &struct {
			Name string `json:"name"`
		}{Name: resolution}
	}
	if comment != "" {
		trUpdateReq.Comment = []struct {
			Add struct {
				Body string `json:"body"`
			} `json:"add"`
		}{
			{Add: struct {
				Body string `json:"body"`
			}{Body: comment}},
		}
	}

	transitionReq := &jira.TransitionRequest{
		Fields: &trFieldsReq,
		Update: &trUpdateReq,
		Transition: &jira.TransitionRequestData{
			ID:   targetTransition.ID.String(),
			Name: targetTransition.Name,
		},
	}

	// Transition all issues
	s := cmdutil.Info(fmt.Sprintf("Transitioning %d issues to %q...", len(normalizedKeys), state))
	defer s.Stop()

	var failed []string
	var succeeded []string

	for _, key := range normalizedKeys {
		_, err := client.Transition(key, transitionReq)
		if err != nil {
			failed = append(failed, key)
			continue
		}
		succeeded = append(succeeded, key)
	}

	s.Stop()

	if len(failed) > 0 {
		if len(succeeded) > 0 {
			cmdutil.Warn("Transitioned %d issues successfully, %d failed", len(succeeded), len(failed))
			fmt.Printf("Failed: %s\n", strings.Join(failed, ", "))
		} else {
			return fmt.Errorf("failed to transition all issues")
		}
	} else {
		cmdutil.Success("Successfully transitioned %d issues to state %q", len(succeeded), state)
	}

	return nil
}

type moveParams struct {
	key        string
	state      string
	comment    string
	assignee   string
	resolution string
	debug      bool
}

func parseArgsAndFlags(flags query.FlagParser, args []string, project string) *moveParams {
	var key, state string

	nargs := len(args)
	if nargs >= 1 {
		key = cmdutil.GetJiraIssueKey(project, args[0])
	}
	if nargs >= 2 {
		state = args[1]
	}

	comment, err := flags.GetString("comment")
	cmdutil.ExitIfError(err)

	assignee, err := flags.GetString("assignee")
	cmdutil.ExitIfError(err)

	resolution, err := flags.GetString("resolution")
	cmdutil.ExitIfError(err)

	debug, err := flags.GetBool("debug")
	cmdutil.ExitIfError(err)

	return &moveParams{
		key:        key,
		state:      state,
		comment:    comment,
		assignee:   assignee,
		resolution: resolution,
		debug:      debug,
	}
}

type moveCmd struct {
	client      *jira.Client
	transitions []*jira.Transition
	params      *moveParams
}

func (mc *moveCmd) setIssueKey(project string) error {
	if mc.params.key != "" {
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
	mc.params.key = cmdutil.GetJiraIssueKey(project, ans)

	return nil
}

func (mc *moveCmd) setDesiredState(it string) error {
	if mc.params.state != "" {
		return nil
	}

	var (
		options = make([]string, 0, len(mc.transitions))
		ans     string
	)

	for _, t := range mc.transitions {
		if it == jira.InstallationTypeCloud && !t.IsAvailable {
			continue
		}
		options = append(options, t.Name)
	}
	options = append(options, optionCancel)

	qs := &survey.Question{
		Name: "state",
		Prompt: &survey.Select{
			Message: "Desired state:",
			Options: options,
		},
		Validate: survey.Required,
	}
	if err := survey.Ask([]*survey.Question{qs}, &ans); err != nil {
		return err
	}
	mc.params.state = ans

	return nil
}

func (mc *moveCmd) setAvailableTransitions() error {
	s := cmdutil.Info("Fetching available transitions. Please wait...")
	defer s.Stop()

	t, err := api.ProxyTransitions(mc.client, mc.params.key)
	if err != nil {
		return err
	}
	mc.transitions = t

	return nil
}

func (mc *moveCmd) verifyTransition(it string) (*jira.Transition, error) {
	var tr *jira.Transition

	st := strings.ToLower(mc.params.state)
	all := make([]string, 0, len(mc.transitions))
	for _, t := range mc.transitions {
		if strings.ToLower(t.Name) == st {
			tr = t
		}
		all = append(all, fmt.Sprintf("'%s'", t.Name))
	}

	if tr == nil {
		return nil, fmt.Errorf(
			"invalid transition state %q\nAvailable states for issue %s: %s",
			mc.params.state, mc.params.key, strings.Join(all, ", "),
		)
	}

	// Jira API v2 doesn't seem to return "isAvailable" field even if the documentation says it does.
	// So, we will only verify if the transition is available for the cloud installation.
	if it == jira.InstallationTypeCloud && !tr.IsAvailable {
		return nil, fmt.Errorf(
			"transition state %q for issue %q is not available",
			mc.params.state, mc.params.key,
		)
	}
	return tr, nil
}
