package block

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
)

const (
	helpText = `Block links an issue as "blocks" and optionally moves it to blocked status.`
	examples = `$ jira quick block PROJ-123 PROJ-456`
)

// NewCmdBlock is a block command.
func NewCmdBlock() *cobra.Command {
	return &cobra.Command{
		Use:     "block <issue-key> <blocked-by-key>",
		Short:   "Block links issue as blocks",
		Long:    helpText,
		Example: examples,
		Args:    cobra.ExactArgs(2),
		Run:     Block,
	}
}

// Block performs the quick block action.
func Block(cmd *cobra.Command, args []string) {
	debug, err := cmd.Flags().GetBool("debug")
	cmdutil.ExitIfError(err)

	project := viper.GetString("project.key")
	issueKey := cmdutil.GetJiraIssueKey(project, args[0])
	blockedByKey := cmdutil.GetJiraIssueKey(project, args[1])

	client := api.DefaultClient(debug)

	// Link as "blocks"
	s := cmdutil.Info("Linking issues...")
	err = client.LinkIssue(issueKey, blockedByKey, "Blocks")
	s.Stop()
	cmdutil.ExitIfError(err)

	cmdutil.Success("Issue %s now blocks %s", issueKey, blockedByKey)
}

