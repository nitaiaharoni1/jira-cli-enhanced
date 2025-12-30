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
	"github.com/ankitpokhrel/jira-cli/pkg/adf"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

const (
	helpText = `List lists all comments for an issue.`
	examples = `# List comments
$ jira issue comment list PROJ-123`
)

// NewCmdList is a list command.
func NewCmdList() *cobra.Command {
	return &cobra.Command{
		Use:     "list ISSUE-KEY",
		Short:   "List comments for an issue",
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

	s := cmdutil.Info("Fetching comments...")
	comments, err := client.GetComments(issueKey)
	s.Stop()

	if err != nil {
		return err
	}

	if len(comments) == 0 {
		fmt.Printf("No comments found for issue %s\n", issueKey)
		return nil
	}

	plain, _ := cmd.Flags().GetBool("plain")
	if plain {
		for _, comment := range comments {
			body := getCommentBody(comment.Body)
			fmt.Printf("%s\t%s\t%s\t%s\n", comment.ID, comment.Author.DisplayName, comment.Created, body)
		}
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)
	fmt.Fprintf(w, "ID\tAUTHOR\tCREATED\tCOMMENT\n")

	for _, comment := range comments {
		created, _ := time.Parse(jira.RFC3339MilliLayout, comment.Created)
		body := getCommentBody(comment.Body)
		// Truncate long comments
		if len(body) > 50 {
			body = body[:47] + "..."
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			comment.ID,
			comment.Author.DisplayName,
			created.Format("2006-01-02 15:04:05"),
			body,
		)
	}

	w.Flush()

	return nil
}

func getCommentBody(body interface{}) string {
	switch v := body.(type) {
	case string:
		return v
	case *adf.ADF:
		if v == nil {
			return ""
		}
		// Convert ADF to plain text (simplified)
		return "[ADF Content]"
	default:
		return fmt.Sprintf("%v", v)
	}
}

