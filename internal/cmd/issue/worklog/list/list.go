package list

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

const (
	helpText = `List lists all worklogs for an issue.`
	examples = `# List worklogs
$ jira issue worklog list PROJ-123`
)

// NewCmdList is a list command.
func NewCmdList() *cobra.Command {
	return &cobra.Command{
		Use:     "list ISSUE-KEY",
		Short:   "List worklogs for an issue",
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

	s := cmdutil.Info("Fetching worklogs...")
	worklogs, err := client.GetWorklogs(issueKey)
	s.Stop()

	if err != nil {
		return err
	}

	if len(worklogs) == 0 {
		fmt.Printf("No worklogs found for issue %s\n", issueKey)
		return nil
	}

	plain, _ := cmd.Flags().GetBool("plain")
	if plain {
		for _, wl := range worklogs {
			fmt.Printf("%s\t%s\t%s\t%s\t%s\n",
				wl.ID,
				wl.Author.DisplayName,
				wl.Started,
				wl.TimeSpent,
				wl.Comment,
			)
		}
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)
	fmt.Fprintf(w, "ID\tAUTHOR\tSTARTED\tTIME SPENT\tCOMMENT\n")

	for _, wl := range worklogs {
		started, _ := time.Parse(jira.RFC3339MilliLayout, wl.Started)
		comment := wl.Comment
		if len(comment) > 40 {
			comment = comment[:37] + "..."
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			wl.ID,
			wl.Author.DisplayName,
			started.Format("2006-01-02 15:04:05"),
			wl.TimeSpent,
			comment,
		)
	}

	w.Flush()

	return nil
}

