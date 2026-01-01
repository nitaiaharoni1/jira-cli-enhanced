package bulk

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
	bulkHelpText = `Bulk label adds or removes labels from multiple issues.

You can label up to 50 issues at once with the same labels.`
	bulkExamples = `# Add labels to multiple issues
$ jira issue label bulk PROJ-1 PROJ-2 PROJ-3 "urgent" "backend"

# Remove labels from multiple issues
$ jira issue label bulk PROJ-1 PROJ-2 --remove "urgent"

# Add labels from stdin
$ jira issue list --keys-only | jira issue label bulk "urgent" "backend"

# Add labels to issues from JQL
$ jira issue label bulk --jql "status = 'To Do'" "ready-for-review"`
)

// NewCmdBulk is a bulk label command.
func NewCmdBulk() *cobra.Command {
	cmd := cobra.Command{
		Use:     "bulk ISSUE-KEY... LABEL...",
		Short:   "Add or remove labels from multiple issues",
		Long:    bulkHelpText,
		Example: bulkExamples,
		Aliases: []string{"batch"},
		Args:    cobra.MinimumNArgs(1),
		RunE:    labelBulk,
	}

	cmd.Flags().Bool("stdin", false, "Read issue keys from stdin (one per line)")
	cmd.Flags().String("jql", "", "Apply to all issues matching JQL query")
	cmd.Flags().Bool("remove", false, "Remove labels instead of adding")

	return &cmd
}

func labelBulk(cmd *cobra.Command, args []string) error {
	project := viper.GetString("project.key")
	debug, _ := cmd.Flags().GetBool("debug")
	stdin, _ := cmd.Flags().GetBool("stdin")
	jql, _ := cmd.Flags().GetString("jql")
	remove, _ := cmd.Flags().GetBool("remove")
	client := api.DefaultClient(debug)

	var issueKeys []string
	var labels []string

	// Parse args - labels are at the end
	// If stdin or jql, all args are labels
	// Otherwise, last args are labels, first args are issue keys
	if stdin || jql != "" {
		labels = args
	} else {
		// Find where labels start (first arg that doesn't look like an issue key)
		labelStart := 0
		for i, arg := range args {
			if !strings.Contains(arg, "-") || len(strings.Split(arg, "-")) != 2 {
				labelStart = i
				break
			}
		}
		if labelStart == 0 {
			return fmt.Errorf("no labels provided")
		}
		issueKeys = args[:labelStart]
		labels = args[labelStart:]
	}

	if len(labels) == 0 {
		return fmt.Errorf("no labels provided")
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
	}

	if len(issueKeys) == 0 {
		return fmt.Errorf("no issues found")
	}

	// Prepare labels
	labelsToApply := labels
	if remove {
		labelsToApply = make([]string, 0, len(labels))
		for _, l := range labels {
			labelsToApply = append(labelsToApply, "-"+l)
		}
	}

	action := "Adding"
	if remove {
		action = "Removing"
	}

	s := cmdutil.Info(fmt.Sprintf("%s labels to %d issues...", action, len(issueKeys)))
	defer s.Stop()

	var failed []string
	var succeeded []string

	for _, key := range issueKeys {
		// Get current issue to preserve existing labels
		issue, err := api.ProxyGetIssue(client, key)
		if err != nil {
			failed = append(failed, key)
			continue
		}

		var newLabels []string
		if remove {
			// Remove labels
			existingLabels := issue.Fields.Labels
			removeMap := make(map[string]bool)
			for _, l := range labels {
				removeMap[l] = true
			}
			newLabels = make([]string, 0)
			for _, l := range existingLabels {
				if !removeMap[l] {
					newLabels = append(newLabels, l)
				}
			}
			// Use minus prefix for removal
			labelsToApply = make([]string, 0, len(labels))
			for _, l := range labels {
				labelsToApply = append(labelsToApply, "-"+l)
			}
		} else {
			// Add labels
			existingLabels := issue.Fields.Labels
			labelMap := make(map[string]bool)
			for _, l := range existingLabels {
				labelMap[l] = true
			}
			for _, l := range labels {
				labelMap[l] = true
			}
			newLabels = make([]string, 0, len(labelMap))
			for l := range labelMap {
				newLabels = append(newLabels, l)
			}
		}

		editReq := &jira.EditRequest{
			Labels: labelsToApply,
		}

		err = client.Edit(key, editReq)
		if err != nil {
			failed = append(failed, key)
			continue
		}
		succeeded = append(succeeded, key)
	}

	s.Stop()

	if len(failed) > 0 {
		if len(succeeded) > 0 {
			cmdutil.Warn("%s labels to %d issues successfully, %d failed", action, len(succeeded), len(failed))
			fmt.Printf("Failed: %s\n", strings.Join(failed, ", "))
		} else {
			return fmt.Errorf("failed to %s labels from all issues", strings.ToLower(action))
		}
	} else {
		cmdutil.Success("Successfully %s labels to %d issues", strings.ToLower(action), len(succeeded))
	}

	return nil
}

