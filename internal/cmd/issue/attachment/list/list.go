package list

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
)

const (
	helpText = `List lists all attachments for an issue.`
	examples = `# List attachments
$ jira issue attachment list PROJ-123`
)

// NewCmdList is a list command.
func NewCmdList() *cobra.Command {
	return &cobra.Command{
		Use:     "list ISSUE-KEY",
		Short:   "List attachments for an issue",
		Long:    helpText,
		Example: examples,
		Args:    cobra.ExactArgs(1),
		RunE:    list,
	}
}

func list(cmd *cobra.Command, args []string) error {
	project := viper.GetString("project.key")
	issueKey := cmdutil.GetJiraIssueKey(project, args[0])

	debug, _ := cmd.Flags().GetBool("debug")
	client := api.DefaultClient(debug)

	s := cmdutil.Info("Fetching attachments...")
	attachments, err := client.GetAttachments(issueKey)
	s.Stop()

	if err != nil {
		return err
	}

	if len(attachments) == 0 {
		fmt.Printf("No attachments found for issue %s\n", issueKey)
		return nil
	}

	plain, _ := cmd.Flags().GetBool("plain")
	if plain {
		for _, att := range attachments {
			fmt.Printf("%s\t%s\t%d\t%s\n", att.ID, att.Filename, att.Size, att.Created)
		}
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)
	fmt.Fprintf(w, "ID\tFILENAME\tSIZE\tCREATED\tAUTHOR\n")
	for _, att := range attachments {
		fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%s\n",
			att.ID,
			att.Filename,
			att.Size,
			att.Created,
			att.Author.DisplayName,
		)
	}
	w.Flush()

	return nil
}

