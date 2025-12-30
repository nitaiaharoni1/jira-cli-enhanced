package edit

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
)

const (
	helpText = `Edit edits an existing comment on an issue.`
	examples = `# Edit comment
$ jira issue comment edit PROJ-123 COMMENT-ID "Updated comment text"

# Edit as internal comment
$ jira issue comment edit PROJ-123 COMMENT-ID "Updated comment" --internal`
)

// NewCmdEdit is an edit command.
func NewCmdEdit() *cobra.Command {
	cmd := cobra.Command{
		Use:     "edit ISSUE-KEY COMMENT-ID BODY",
		Short:   "Edit a comment on an issue",
		Long:    helpText,
		Example: examples,
		Args:    cobra.ExactArgs(3),
		RunE:    edit,
	}

	cmd.Flags().Bool("internal", false, "Make comment internal (visible to admins only)")

	return &cmd
}

func edit(cmd *cobra.Command, args []string) error {
	project := viper.GetString("project.key")
	issueKey := cmdutil.GetJiraIssueKey(project, args[0])
	commentID := args[1]
	body := args[2]

	internal, _ := cmd.Flags().GetBool("internal")
	debug, _ := cmd.Flags().GetBool("debug")

	client := api.DefaultClient(debug)

	s := cmdutil.Info("Updating comment...")
	err := client.UpdateComment(issueKey, commentID, body, internal)
	s.Stop()

	if err != nil {
		return err
	}

	cmdutil.Success("Comment updated successfully")
	return nil
}

