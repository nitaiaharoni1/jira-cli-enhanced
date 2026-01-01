package filter

import (
	"github.com/spf13/cobra"

	"github.com/ankitpokhrel/jira-cli/internal/cmd/filter/create"
	"github.com/ankitpokhrel/jira-cli/internal/cmd/filter/execute"
	"github.com/ankitpokhrel/jira-cli/internal/cmd/filter/favorite"
	"github.com/ankitpokhrel/jira-cli/internal/cmd/filter/list"
)

const helpText = `Filter manages saved Jira filters. See available commands below.`

// NewCmdFilter is a filter command.
func NewCmdFilter() *cobra.Command {
	cmd := cobra.Command{
		Use:         "filter",
		Short:       "Filter manages saved Jira filters",
		Long:        helpText,
		Aliases:     []string{"filters"},
		Annotations: map[string]string{"cmd:main": "true"},
		RunE:        filter,
	}

	lc := list.NewCmdList()
	cc := create.NewCmdCreate()
	ec := execute.NewCmdExecute()
	fc := favorite.NewCmdFavorite()

	cmd.AddCommand(lc, cc, ec, fc)

	list.SetFlags(lc)

	return &cmd
}

func filter(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}

