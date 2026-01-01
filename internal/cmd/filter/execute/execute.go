package execute

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/internal/view"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
	"github.com/ankitpokhrel/jira-cli/pkg/tui"
)

var limit uint

// NewCmdExecute is an execute command.
func NewCmdExecute() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "execute <filter-id>",
		Short:   "Execute executes a saved filter and shows matching issues",
		Long:    "Execute executes a saved filter and displays the matching issues.",
		Aliases: []string{"run", "search"},
		Args:    cobra.ExactArgs(1),
		Run:     Execute,
	}

	cmd.Flags().UintVarP(&limit, "limit", "l", 50, "Maximum number of issues to return")

	return cmd
}

// Execute executes a filter and shows the results.
func Execute(cmd *cobra.Command, args []string) {
	debug, err := cmd.Flags().GetBool("debug")
	cmdutil.ExitIfError(err)

	filterID := args[0]
	limit, err := cmd.Flags().GetUint("limit")
	cmdutil.ExitIfError(err)

	server := viper.GetString("server")
	project := viper.GetString("project.key")

	s := cmdutil.Info("Executing filter...")
	defer s.Stop()

	result, err := api.DefaultClient(debug).ExecuteFilter(filterID, limit)
	cmdutil.ExitIfError(err)

	if len(result.Issues) == 0 {
		cmdutil.Failed("No issues found matching the filter.")
		return
	}

	// Use the existing issue list view to display results
	display := view.DisplayFormat{
		Plain:        false,
		Delimiter:    "\t",
		CSV:          false,
		NoHeaders:    false,
		NoTruncate:   false,
		FixedColumns: 0,
		Comments:     0,
		TableStyle:   cmdutil.GetTUIStyleConfig(),
		Timezone:     viper.GetString("timezone"),
	}

	issueList := view.IssueList{
		Project: project,
		Server:  server,
		Data:    result.Issues,
		Display: display,
		FooterText: fmt.Sprintf("Showing %d results from filter %s", len(result.Issues), filterID),
	}

	if tui.IsDumbTerminal() || tui.IsNotTTY() {
		display.Plain = true
		issueList.Display = display
	}

	cmdutil.ExitIfError(issueList.Render())
}

