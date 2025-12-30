package history

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

const (
	helpText = `History displays the changelog/history of an issue.`
	examples = `# View issue history
$ jira issue history PROJ-123

# Filter by field
$ jira issue history PROJ-123 --field status`
)

// NewCmdHistory is a history command.
func NewCmdHistory() *cobra.Command {
	cmd := cobra.Command{
		Use:     "history ISSUE-KEY",
		Short:   "Display issue changelog/history",
		Long:    helpText,
		Example: examples,
		Aliases: []string{"changelog", "changes"},
		Args:    cobra.ExactArgs(1),
		RunE:    history,
	}

	cmd.Flags().String("field", "", "Filter by field name")

	return &cmd
}

func history(cmd *cobra.Command, args []string) error {
	project := viper.GetString("project.key")
	issueKey := cmdutil.GetJiraIssueKey(project, args[0])

	fieldFilter, _ := cmd.Flags().GetString("field")
	debug, _ := cmd.Flags().GetBool("debug")

	client := api.DefaultClient(debug)

	s := cmdutil.Info("Fetching issue history...")
	historyFlat, err := client.GetIssueHistoryFlat(issueKey)
	s.Stop()

	if err != nil {
		return err
	}

	if len(historyFlat) == 0 {
		fmt.Printf("No history found for issue %s\n", issueKey)
		return nil
	}

	plain, _ := cmd.Flags().GetBool("plain")
	if plain {
		for _, h := range historyFlat {
			if fieldFilter == "" || strings.EqualFold(h.Field, fieldFilter) {
				fmt.Printf("%s\t%s\t%s\t%s\t%s\n",
					h.Created,
					h.Author.DisplayName,
					h.Field,
					h.FromString,
					h.ToString,
				)
			}
		}
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)
	fmt.Fprintf(w, "DATE\tAUTHOR\tFIELD\tFROM\tTO\n")

	for _, h := range historyFlat {
		if fieldFilter == "" || strings.EqualFold(h.Field, fieldFilter) {
			created, _ := time.Parse(jira.RFC3339MilliLayout, h.Created)
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				created.Format("2006-01-02 15:04:05"),
				h.Author.DisplayName,
				h.Field,
				h.FromString,
				h.ToString,
			)
		}
	}

	w.Flush()

	return nil
}

