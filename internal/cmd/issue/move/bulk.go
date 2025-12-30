package move

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
	bulkHelpText = `Bulk move transitions multiple issues from one state to another.

You can transition up to 50 issues at once. All issues will be transitioned to the same state.`
	bulkExamples = `# Transition multiple issues to "In Progress"
$ jira issue move-bulk PROJ-1 PROJ-2 PROJ-3 "In Progress"

# Transition with comment
$ jira issue move-bulk PROJ-1 PROJ-2 PROJ-3 Done --comment "All completed" -RFixed`
)

// NewCmdMoveBulk is a bulk move command.
func NewCmdMoveBulk() *cobra.Command {
	cmd := cobra.Command{
		Use:     "move-bulk ISSUE-KEY... STATE",
		Short:   "Transition multiple issues to a given state",
		Long:    bulkHelpText,
		Example: bulkExamples,
		Aliases: []string{"move-batch", "transition-bulk"},
		Args:    cobra.MinimumNArgs(2),
		RunE:    moveBulk,
	}

	cmd.Flags().String("comment", "", "Add comment to all issues")
	cmd.Flags().StringP("assignee", "a", "", "Assign all issues to a user")
	cmd.Flags().StringP("resolution", "R", "", "Set resolution for all issues")

	return &cmd
}

func moveBulk(cmd *cobra.Command, args []string) error {
	project := viper.GetString("project.key")

	// Last argument is the state
	state := args[len(args)-1]
	issueKeys := args[:len(args)-1]

	// Normalize issue keys
	normalizedKeys := make([]string, 0, len(issueKeys))
	for _, key := range issueKeys {
		normalizedKeys = append(normalizedKeys, cmdutil.GetJiraIssueKey(project, key))
	}

	comment, _ := cmd.Flags().GetString("comment")
	assignee, _ := cmd.Flags().GetString("assignee")
	resolution, _ := cmd.Flags().GetString("resolution")
	debug, _ := cmd.Flags().GetBool("debug")

	client := api.DefaultClient(debug)

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

