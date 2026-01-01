package delete

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
)

const (
	helpText = `Delete deletes a comment from an issue.`
	examples = `# Delete comment
$ jira issue comment delete PROJ-123 COMMENT-ID`
)

// NewCmdDelete is a delete command.
func NewCmdDelete() *cobra.Command {
	return &cobra.Command{
		Use:     "delete ISSUE-KEY COMMENT-ID",
		Short:   "Delete a comment from an issue",
		Long:    helpText,
		Example: examples,
		Args:    cobra.ExactArgs(2),
		RunE:    delete,
	}
}

func delete(cmd *cobra.Command, args []string) error {
	project := viper.GetString("project.key")
	issueKey := cmdutil.GetJiraIssueKey(project, args[0])
	commentID := args[1]

	debug, _ := cmd.Flags().GetBool("debug")
	client := api.DefaultClient(debug)

	s := cmdutil.Info("Deleting comment...")
	err := client.DeleteComment(issueKey, commentID)
	s.Stop()

	if err != nil {
		return err
	}

	cmdutil.Success("Comment deleted successfully")
	return nil
}


