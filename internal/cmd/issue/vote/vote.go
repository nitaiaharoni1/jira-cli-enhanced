package vote

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
)

const (
	helpText = `Vote adds your vote to an issue.`
	examples = `$ jira issue vote PROJ-123`
)

// NewCmdVote is a vote command.
func NewCmdVote() *cobra.Command {
	return &cobra.Command{
		Use:     "vote <issue-key>",
		Short:   "Vote adds your vote to an issue",
		Long:    helpText,
		Example: examples,
		Args:    cobra.ExactArgs(1),
		Run:     Vote,
	}
}

// Vote adds a vote to an issue.
func Vote(cmd *cobra.Command, args []string) {
	debug, err := cmd.Flags().GetBool("debug")
	cmdutil.ExitIfError(err)

	project := viper.GetString("project.key")
	key := cmdutil.GetJiraIssueKey(project, args[0])

	s := cmdutil.Info("Voting on issue...")
	defer s.Stop()

	err = api.DefaultClient(debug).VoteIssue(key)
	cmdutil.ExitIfError(err)

	cmdutil.Success("Vote added to issue %s", key)
}

