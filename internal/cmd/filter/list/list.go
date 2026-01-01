package list

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/internal/view"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

var (
	favoriteOnly bool
	plain        bool
)

// NewCmdList is a list command.
func NewCmdList() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List lists saved Jira filters",
		Long:    "List lists saved Jira filters that you have access to.",
		Aliases: []string{"lists", "ls"},
		Run:     List,
	}
}

// SetFlags sets flags for the list command.
func SetFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&favoriteOnly, "favorite", "f", false, "Show only favorite filters")
	cmd.Flags().BoolVar(&plain, "plain", false, "Display output in plain mode")
}

// List displays a list view.
func List(cmd *cobra.Command, _ []string) {
	debug, err := cmd.Flags().GetBool("debug")
	cmdutil.ExitIfError(err)

	favoriteOnly, err := cmd.Flags().GetBool("favorite")
	cmdutil.ExitIfError(err)

	plain, err := cmd.Flags().GetBool("plain")
	cmdutil.ExitIfError(err)

	var filters []*jira.SavedFilter
	var total int

	if favoriteOnly {
		filters, total, err = func() ([]*jira.SavedFilter, int, error) {
			s := cmdutil.Info("Fetching favorite filters...")
			defer s.Stop()

			f, err := api.DefaultClient(debug).GetFilters()
			if err != nil {
				return nil, 0, err
			}
			return f, len(f), nil
		}()
	} else {
		filters, total, err = func() ([]*jira.SavedFilter, int, error) {
			s := cmdutil.Info("Fetching filters...")
			defer s.Stop()

			maxResults := viper.GetInt("filter.max_results")
			if maxResults == 0 {
				maxResults = 50
			}

			resp, err := api.DefaultClient(debug).GetAllFilters(0, maxResults)
			if err != nil {
				return nil, 0, err
			}
			return resp.Values, resp.Total, nil
		}()
	}

	cmdutil.ExitIfError(err)

	if total == 0 {
		cmdutil.Failed("No filters found.")
		return
	}

	if plain {
		for _, f := range filters {
			cmdutil.Success("%s\t%s\t%s", f.ID, f.Name, f.JQL)
		}
		return
	}

	v := view.NewSavedFilter(filters)
	cmdutil.ExitIfError(v.Render())
}

