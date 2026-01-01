package comment

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

const (
	bulkHelpText = `Bulk comment adds a comment to multiple issues.

You can comment on up to 50 issues at once with the same comment text.`
	bulkExamples = `# Add comment to multiple issues
$ jira issue comment-bulk PROJ-1 PROJ-2 PROJ-3 "Fixed in latest release"

# Add comment from stdin
$ jira issue list --keys-only | jira issue comment-bulk "Deployed to production"

# Add comment to issues from JQL
$ jira issue comment-bulk --jql "status = 'Done'" "Ready for review"

# Add internal comment
$ jira issue comment-bulk PROJ-1 PROJ-2 --internal "Internal note"`
)

// NewCmdCommentBulk is a bulk comment command.
func NewCmdCommentBulk() *cobra.Command {
	cmd := cobra.Command{
		Use:     "comment-bulk ISSUE-KEY... COMMENT",
		Short:   "Add comment to multiple issues",
		Long:    bulkHelpText,
		Example: bulkExamples,
		Aliases: []string{"comment-batch"},
		Args:    cobra.MinimumNArgs(1),
		RunE:    commentBulk,
	}

	cmd.Flags().Bool("stdin", false, "Read issue keys from stdin (one per line)")
	cmd.Flags().String("jql", "", "Apply to all issues matching JQL query")
	cmd.Flags().Bool("internal", false, "Add as internal comment")

	return &cmd
}

func commentBulk(cmd *cobra.Command, args []string) error {
	project := viper.GetString("project.key")
	debug, _ := cmd.Flags().GetBool("debug")
	stdin, _ := cmd.Flags().GetBool("stdin")
	jql, _ := cmd.Flags().GetString("jql")
	internal, _ := cmd.Flags().GetBool("internal")
	client := api.DefaultClient(debug)

	var issueKeys []string
	var comment string

	// Get comment from args (last argument)
	if len(args) > 0 {
		comment = args[len(args)-1]
		args = args[:len(args)-1]
	}

	if comment == "" {
		return fmt.Errorf("comment text required")
	}

	// Get issue keys
	if stdin {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			key := strings.TrimSpace(scanner.Text())
			if key != "" {
				issueKeys = append(issueKeys, cmdutil.GetJiraIssueKey(project, key))
			}
		}
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("failed to read from stdin: %w", err)
		}
	} else if jql != "" {
		result, err := api.ProxySearch(client, jql, 0, 1000)
		if err != nil {
			return fmt.Errorf("failed to search issues: %w", err)
		}
		for _, issue := range result.Issues {
			issueKeys = append(issueKeys, issue.Key)
		}
	} else {
		if len(args) == 0 {
			return fmt.Errorf("no issue keys provided")
		}
		for _, key := range args {
			issueKeys = append(issueKeys, cmdutil.GetJiraIssueKey(project, key))
		}
	}

	if len(issueKeys) == 0 {
		return fmt.Errorf("no issues found")
	}

	s := cmdutil.Info(fmt.Sprintf("Adding comment to %d issues...", len(issueKeys)))
	defer s.Stop()

	var failed []string
	var succeeded []string

	for _, key := range issueKeys {
		err := client.AddIssueComment(key, comment, internal)
		if err != nil {
			failed = append(failed, key)
			continue
		}
		succeeded = append(succeeded, key)
	}

	s.Stop()

	if len(failed) > 0 {
		if len(succeeded) > 0 {
			cmdutil.Warn("Added comment to %d issues successfully, %d failed", len(succeeded), len(failed))
			fmt.Printf("Failed: %s\n", strings.Join(failed, ", "))
		} else {
			return fmt.Errorf("failed to add comment to all issues")
		}
	} else {
		cmdutil.Success("Successfully added comment to %d issues", len(succeeded))
	}

	return nil
}

