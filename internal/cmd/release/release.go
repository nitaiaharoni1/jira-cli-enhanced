package release

import (
	"github.com/spf13/cobra"

	"github.com/ankitpokhrel/jira-cli/internal/cmd/release/create"
	"github.com/ankitpokhrel/jira-cli/internal/cmd/release/delete"
	"github.com/ankitpokhrel/jira-cli/internal/cmd/release/list"
	"github.com/ankitpokhrel/jira-cli/internal/cmd/release/update"
)

const helpText = `Release manages Jira Project versions. See available commands below.`

// NewCmdRelease is a release command.
func NewCmdRelease() *cobra.Command {
	cmd := cobra.Command{
		Use:         "release",
		Short:       "Release manages Jira Project versions",
		Long:        helpText,
		Annotations: map[string]string{"cmd:main": "true"},
		Aliases:     []string{"releases"},
		RunE:        releases,
	}

	cmd.AddCommand(
		list.NewCmdList(),
		create.NewCmdCreate(),
		update.NewCmdUpdate(),
		delete.NewCmdDelete(),
	)

	return &cmd
}

func releases(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}
