package start

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

const (
	helpText = `Start moves an issue to "In Progress" and assigns it to yourself.`
	examples = `$ jira quick start PROJ-123`
)

// NewCmdStart is a start command.
func NewCmdStart() *cobra.Command {
	return &cobra.Command{
		Use:     "start <issue-key>",
		Short:   "Start moves issue to In Progress and assigns to self",
		Long:    helpText,
		Example: examples,
		Args:    cobra.ExactArgs(1),
		Run:     Start,
	}
}

// Start performs the quick start action.
func Start(cmd *cobra.Command, args []string) {
	debug, err := cmd.Flags().GetBool("debug")
	cmdutil.ExitIfError(err)

	project := viper.GetString("project.key")
	key := cmdutil.GetJiraIssueKey(project, args[0])

	client := api.DefaultClient(debug)

	// Get current user
	me, err := client.Me()
	cmdutil.ExitIfError(err)

	// Assign to self
	s := cmdutil.Info("Assigning issue to self...")
	err = api.ProxyAssignIssue(client, key, me, "")
	s.Stop()
	cmdutil.ExitIfError(err)

	// Move to "In Progress"
	s = cmdutil.Info("Moving issue to In Progress...")
	transitions, err := api.ProxyTransitions(client, key)
	cmdutil.ExitIfError(err)

	var inProgressTransition *jira.Transition
	for _, t := range transitions {
		if t.Name == "In Progress" || t.Name == "Start Progress" || t.Name == "Start" {
			inProgressTransition = t
			break
		}
	}

	if inProgressTransition == nil {
		s.Stop()
		cmdutil.Failed("Could not find 'In Progress' transition for issue %s", key)
		return
	}

	trReq := &jira.TransitionRequest{
		Transition: &jira.TransitionRequestData{
			ID:   inProgressTransition.ID.String(),
			Name: inProgressTransition.Name,
		},
	}

	_, err = client.Transition(key, trReq)
	s.Stop()
	cmdutil.ExitIfError(err)

	cmdutil.Success("Issue %s started: assigned to %s and moved to In Progress", key, me.DisplayName)
}

