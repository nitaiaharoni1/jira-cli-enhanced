package update

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
)

const (
	helpText = `Update updates an existing worklog entry.`
	examples = `# Update worklog
$ jira issue worklog update PROJ-123 WORKLOG-ID "2h" "Updated comment"

# Update with start date
$ jira issue worklog update PROJ-123 WORKLOG-ID "2h" "Comment" --started "2025-01-01T10:00:00"`
)

// NewCmdUpdate is an update command.
func NewCmdUpdate() *cobra.Command {
	cmd := cobra.Command{
		Use:     "update ISSUE-KEY WORKLOG-ID TIME-SPENT COMMENT",
		Short:   "Update a worklog entry",
		Long:    helpText,
		Example: examples,
		Args:    cobra.ExactArgs(4),
		RunE:    update,
	}

	cmd.Flags().String("started", "", "Start date/time (RFC3339 format)")

	return &cmd
}

func update(cmd *cobra.Command, args []string) error {
	project := viper.GetString("project.key")
	issueKey := cmdutil.GetJiraIssueKey(project, args[0])
	worklogID := args[1]
	timeSpent := args[2]
	comment := args[3]

	started, _ := cmd.Flags().GetString("started")
	debug, _ := cmd.Flags().GetBool("debug")

	client := api.DefaultClient(debug)

	s := cmdutil.Info("Updating worklog...")
	err := client.UpdateWorklog(issueKey, worklogID, started, timeSpent, comment)
	s.Stop()

	if err != nil {
		return err
	}

	cmdutil.Success("Worklog updated successfully")
	return nil
}

