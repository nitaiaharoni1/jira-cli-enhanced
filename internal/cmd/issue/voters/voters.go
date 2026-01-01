package voters

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
)

const (
	helpText = `Voters lists all voters for an issue.`
	examples = `$ jira issue voters PROJ-123`
)

// NewCmdVoters is a voters command.
func NewCmdVoters() *cobra.Command {
	return &cobra.Command{
		Use:     "voters <issue-key>",
		Short:   "Voters lists all voters for an issue",
		Long:    helpText,
		Example: examples,
		Args:    cobra.ExactArgs(1),
		Run:     Voters,
	}
}

// Voters lists voters for an issue.
func Voters(cmd *cobra.Command, args []string) {
	debug, err := cmd.Flags().GetBool("debug")
	cmdutil.ExitIfError(err)

	project := viper.GetString("project.key")
	key := cmdutil.GetJiraIssueKey(project, args[0])

	s := cmdutil.Info("Fetching voters...")
	defer s.Stop()

	voters, err := api.DefaultClient(debug).GetVoters(key)
	cmdutil.ExitIfError(err)

	if voters.Votes == 0 {
		cmdutil.Failed("No votes for issue %s", key)
		return
	}

	fmt.Printf("\nTotal votes: %d\n", voters.Votes)
	if voters.HasVoted {
		fmt.Println("You have voted on this issue")
	}

	if len(voters.Voters) > 0 {
		fmt.Println("\nVoters:")
		for _, voter := range voters.Voters {
			email := voter.Email
			if email == "" {
				email = voter.Name
			}
			fmt.Printf("  - %s (%s)\n", voter.DisplayName, email)
		}
	}
}

