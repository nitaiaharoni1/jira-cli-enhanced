package remove

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

const (
	helpText = `Remove removes labels from an issue.`
	examples = `# Remove a single label
$ jira issue label remove PROJ-123 "urgent"

# Remove multiple labels
$ jira issue label remove PROJ-123 "urgent" "backend"`
)

// NewCmdRemove is a remove label command.
func NewCmdRemove() *cobra.Command {
	return &cobra.Command{
		Use:     "remove ISSUE-KEY LABEL...",
		Short:   "Remove labels from an issue",
		Long:    helpText,
		Example: examples,
		Args:    cobra.MinimumNArgs(2),
		RunE:    removeLabels,
	}
}

func removeLabels(cmd *cobra.Command, args []string) error {
	project := viper.GetString("project.key")
	issueKey := cmdutil.GetJiraIssueKey(project, args[0])
	labelsToRemove := args[1:]

	debug, _ := cmd.Flags().GetBool("debug")
	client := api.DefaultClient(debug)

	// Get current issue to get existing labels
	issue, err := api.ProxyGetIssue(client, issueKey)
	if err != nil {
		return fmt.Errorf("failed to get issue: %w", err)
	}

	// Remove specified labels
	existingLabels := issue.Fields.Labels
	removeMap := make(map[string]bool)
	for _, l := range labelsToRemove {
		removeMap[l] = true
	}

	newLabels := make([]string, 0)
	for _, l := range existingLabels {
		if !removeMap[l] {
			newLabels = append(newLabels, l)
		}
	}

	// Use minus prefix to remove labels
	labelsWithMinus := make([]string, 0, len(labelsToRemove))
	for _, l := range labelsToRemove {
		labelsWithMinus = append(labelsWithMinus, "-"+l)
	}

	editReq := &jira.EditRequest{
		Labels: labelsWithMinus,
	}

	s := cmdutil.Info(fmt.Sprintf("Removing labels from issue %q...", issueKey))
	defer s.Stop()

	err = client.Edit(issueKey, editReq)
	if err != nil {
		return fmt.Errorf("failed to remove labels: %w", err)
	}

	s.Stop()
	cmdutil.Success("Labels removed from issue %q", issueKey)
	fmt.Printf("%s\n", cmdutil.GenerateServerBrowseURL(viper.GetString("server"), issueKey))

	return nil
}

