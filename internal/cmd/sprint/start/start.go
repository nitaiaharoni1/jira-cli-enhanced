package start

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
)

const (
	helpText = `Start starts a sprint by changing its status to active.`
	examples = `$ jira sprint start 123`
)

// NewCmdStart is a start command.
func NewCmdStart() *cobra.Command {
	return &cobra.Command{
		Use:     "start <sprint-id>",
		Short:   "Start starts a sprint",
		Long:    helpText,
		Example: examples,
		Args:    cobra.ExactArgs(1),
		Run:     Start,
	}
}

// Start starts a sprint.
func Start(cmd *cobra.Command, args []string) {
	debug, err := cmd.Flags().GetBool("debug")
	cmdutil.ExitIfError(err)

	sprintID, err := strconv.Atoi(args[0])
	cmdutil.ExitIfError(err)

	s := cmdutil.Info("Starting sprint...")
	defer s.Stop()

	err = api.DefaultClient(debug).StartSprint(sprintID)
	cmdutil.ExitIfError(err)

	cmdutil.Success("Sprint %d started successfully", sprintID)
}

