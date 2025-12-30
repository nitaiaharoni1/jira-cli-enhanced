package worklog

import (
	"github.com/spf13/cobra"

	"github.com/ankitpokhrel/jira-cli/internal/cmd/issue/worklog/add"
	"github.com/ankitpokhrel/jira-cli/internal/cmd/issue/worklog/delete"
	"github.com/ankitpokhrel/jira-cli/internal/cmd/issue/worklog/list"
	"github.com/ankitpokhrel/jira-cli/internal/cmd/issue/worklog/update"
)

const helpText = `Worklog command helps you manage issue worklogs. See available commands below.`

// NewCmdWorklog is a worklog command.
func NewCmdWorklog() *cobra.Command {
	cmd := cobra.Command{
		Use:     "worklog",
		Short:   "Manage issue worklog",
		Long:    helpText,
		Aliases: []string{"wlg"},
		RunE:    worklog,
	}

	cmd.AddCommand(
		add.NewCmdWorklogAdd(),
		list.NewCmdList(),
		update.NewCmdUpdate(),
		delete.NewCmdDelete(),
	)

	return &cmd
}

func worklog(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}
