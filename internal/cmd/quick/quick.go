package quick

import (
	"github.com/spf13/cobra"

	"github.com/ankitpokhrel/jira-cli/internal/cmd/quick/block"
	"github.com/ankitpokhrel/jira-cli/internal/cmd/quick/done"
	"github.com/ankitpokhrel/jira-cli/internal/cmd/quick/start"
)

const helpText = `Quick provides quick actions for common workflows. See available commands below.`

// NewCmdQuick is a quick command.
func NewCmdQuick() *cobra.Command {
	cmd := cobra.Command{
		Use:         "quick",
		Short:       "Quick provides quick actions for common workflows",
		Long:        helpText,
		Aliases:     []string{"q"},
		Annotations: map[string]string{"cmd:main": "true"},
		RunE:        quick,
	}

	cmd.AddCommand(
		start.NewCmdStart(),
		done.NewCmdDone(),
		block.NewCmdBlock(),
	)

	return &cmd
}

func quick(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}

