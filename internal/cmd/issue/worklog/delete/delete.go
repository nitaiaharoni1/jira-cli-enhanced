package delete

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
)

const (
	helpText = `Delete deletes a worklog entry from an issue.`
	examples = `# Delete worklog
$ jira issue worklog delete PROJ-123 WORKLOG-ID`
)

// NewCmdDelete is a delete command.
func NewCmdDelete() *cobra.Command {
	return &cobra.Command{
		Use:     "delete ISSUE-KEY WORKLOG-ID",
		Short:   "Delete a worklog entry",
		Long:    helpText,
		Example: examples,
		Args:    cobra.ExactArgs(2),
		RunE:    delete,
	}
}

func delete(cmd *cobra.Command, args []string) error {
	project := viper.GetString("project.key")
	issueKey := cmdutil.GetJiraIssueKey(project, args[0])
	worklogID := args[1]

	debug, _ := cmd.Flags().GetBool("debug")
	client := api.DefaultClient(debug)

	s := cmdutil.Info("Deleting worklog...")
	err := client.DeleteWorklog(issueKey, worklogID)
	s.Stop()

	if err != nil {
		return err
	}

	cmdutil.Success("Worklog deleted successfully")
	return nil
}


