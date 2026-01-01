package favorite

import (
	"github.com/spf13/cobra"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/internal/view"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

// NewCmdFavorite is a favorite command.
func NewCmdFavorite() *cobra.Command {
	return &cobra.Command{
		Use:     "favorite",
		Short:   "Favorite lists favorite filters",
		Long:    "Favorite lists all filters marked as favorites.",
		Aliases: []string{"fav", "favorites"},
		Run:     Favorite,
	}
}

// Favorite displays favorite filters.
func Favorite(cmd *cobra.Command, _ []string) {
	debug, err := cmd.Flags().GetBool("debug")
	cmdutil.ExitIfError(err)

	filters, total, err := func() ([]*jira.SavedFilter, int, error) {
		s := cmdutil.Info("Fetching favorite filters...")
		defer s.Stop()

		f, err := api.DefaultClient(debug).GetFilters()
		if err != nil {
			return nil, 0, err
		}
		return f, len(f), nil
	}()

	cmdutil.ExitIfError(err)

	if total == 0 {
		cmdutil.Failed("No favorite filters found.")
		return
	}

	v := view.NewSavedFilter(filters)
	cmdutil.ExitIfError(v.Render())
}

