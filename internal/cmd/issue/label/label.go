package label

import (
	"github.com/spf13/cobra"

	"github.com/ankitpokhrel/jira-cli/internal/cmd/issue/label/add"
	"github.com/ankitpokhrel/jira-cli/internal/cmd/issue/label/bulk"
	"github.com/ankitpokhrel/jira-cli/internal/cmd/issue/label/remove"
)

const helpText = `Label command helps you manage issue labels. See available commands below.`

// NewCmdLabel is a label command.
func NewCmdLabel() *cobra.Command {
	cmd := cobra.Command{
		Use:     "label",
		Short:   "Manage issue labels",
		Long:    helpText,
		Aliases: []string{"labels"},
		RunE:    label,
	}

	cmd.AddCommand(
		add.NewCmdAdd(),
		remove.NewCmdRemove(),
		bulk.NewCmdBulk(),
	)

	return &cmd
}

func label(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}

