package done

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

const (
	helpText = `Done moves an issue to "Done" status.`
	examples = `$ jira quick done PROJ-123`
)

// NewCmdDone is a done command.
func NewCmdDone() *cobra.Command {
	return &cobra.Command{
		Use:     "done <issue-key>",
		Short:   "Done moves issue to Done status",
		Long:    helpText,
		Example: examples,
		Args:    cobra.ExactArgs(1),
		Run:     Done,
	}
}

// Done performs the quick done action.
func Done(cmd *cobra.Command, args []string) {
	debug, err := cmd.Flags().GetBool("debug")
	cmdutil.ExitIfError(err)

	project := viper.GetString("project.key")
	key := cmdutil.GetJiraIssueKey(project, args[0])

	client := api.DefaultClient(debug)

	// Move to "Done"
	s := cmdutil.Info("Moving issue to Done...")
	transitions, err := api.ProxyTransitions(client, key)
	cmdutil.ExitIfError(err)

	var doneTransition *jira.Transition
	for _, t := range transitions {
		if t.Name == "Done" || t.Name == "Close" || t.Name == "Resolve" {
			doneTransition = t
			break
		}
	}

	if doneTransition == nil {
		s.Stop()
		cmdutil.Failed("Could not find 'Done' transition for issue %s", key)
		return
	}

	trReq := &jira.TransitionRequest{
		Transition: &jira.TransitionRequestData{
			ID:   doneTransition.ID.String(),
			Name: doneTransition.Name,
		},
	}

	_, err = client.Transition(key, trReq)
	s.Stop()
	cmdutil.ExitIfError(err)

	cmdutil.Success("Issue %s marked as done", key)
}

