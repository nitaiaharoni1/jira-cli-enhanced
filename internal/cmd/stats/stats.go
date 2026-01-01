package stats

import (
	"github.com/spf13/cobra"

	"github.com/ankitpokhrel/jira-cli/internal/cmd/stats/assigned"
	"github.com/ankitpokhrel/jira-cli/internal/cmd/stats/sprint"
	"github.com/ankitpokhrel/jira-cli/internal/cmd/stats/velocity"
	"github.com/ankitpokhrel/jira-cli/internal/cmd/stats/worklog"
)

const helpText = `Stats provides statistics and reporting for sprints, worklogs, and issues.`

// NewCmdStats is a stats command.
func NewCmdStats() *cobra.Command {
	cmd := cobra.Command{
		Use:         "stats",
		Short:       "Statistics and reporting",
		Long:        helpText,
		Aliases:     []string{"stat", "report"},
		Annotations: map[string]string{"cmd:main": "true"},
		RunE:        stats,
	}

	cmd.AddCommand(
		sprint.NewCmdSprint(),
		velocity.NewCmdVelocity(),
		worklog.NewCmdWorklog(),
		assigned.NewCmdAssigned(),
	)

	return &cmd
}

func stats(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}

