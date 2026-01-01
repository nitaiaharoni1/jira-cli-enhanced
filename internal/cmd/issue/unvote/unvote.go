package unvote

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
)

const (
	helpText = `Unvote removes your vote from an issue.`
	examples = `$ jira issue unvote PROJ-123`
)

// NewCmdUnvote is an unvote command.
func NewCmdUnvote() *cobra.Command {
	return &cobra.Command{
		Use:     "unvote <issue-key>",
		Short:   "Unvote removes your vote from an issue",
		Long:    helpText,
		Example: examples,
		Args:    cobra.ExactArgs(1),
		Run:     Unvote,
	}
}

// Unvote removes a vote from an issue.
func Unvote(cmd *cobra.Command, args []string) {
	debug, err := cmd.Flags().GetBool("debug")
	cmdutil.ExitIfError(err)

	project := viper.GetString("project.key")
	key := cmdutil.GetJiraIssueKey(project, args[0])

	s := cmdutil.Info("Removing vote from issue...")
	defer s.Stop()

	err = api.DefaultClient(debug).UnvoteIssue(key)
	cmdutil.ExitIfError(err)

	cmdutil.Success("Vote removed from issue %s", key)
}

