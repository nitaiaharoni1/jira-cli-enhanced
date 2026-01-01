package add

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

const (
	helpText = `Add adds labels to an issue.`
	examples = `# Add a single label
$ jira issue label add PROJ-123 "urgent"

# Add multiple labels
$ jira issue label add PROJ-123 "urgent" "backend" "high-priority"`
)

// NewCmdAdd is an add label command.
func NewCmdAdd() *cobra.Command {
	return &cobra.Command{
		Use:     "add ISSUE-KEY LABEL...",
		Short:   "Add labels to an issue",
		Long:    helpText,
		Example: examples,
		Args:    cobra.MinimumNArgs(2),
		RunE:    addLabels,
	}
}

func addLabels(cmd *cobra.Command, args []string) error {
	project := viper.GetString("project.key")
	issueKey := cmdutil.GetJiraIssueKey(project, args[0])
	labels := args[1:]

	debug, _ := cmd.Flags().GetBool("debug")
	client := api.DefaultClient(debug)

	// Get current issue to preserve existing labels
	issue, err := api.ProxyGetIssue(client, issueKey)
	if err != nil {
		return fmt.Errorf("failed to get issue: %w", err)
	}

	// Merge with existing labels
	existingLabels := issue.Fields.Labels
	labelMap := make(map[string]bool)
	for _, l := range existingLabels {
		labelMap[l] = true
	}
	for _, l := range labels {
		labelMap[l] = true
	}

	newLabels := make([]string, 0, len(labelMap))
	for l := range labelMap {
		newLabels = append(newLabels, l)
	}

	editReq := &jira.EditRequest{
		Labels: newLabels,
	}

	s := cmdutil.Info(fmt.Sprintf("Adding labels to issue %q...", issueKey))
	defer s.Stop()

	err = client.Edit(issueKey, editReq)
	if err != nil {
		return fmt.Errorf("failed to add labels: %w", err)
	}

	s.Stop()
	cmdutil.Success("Labels added to issue %q", issueKey)
	fmt.Printf("%s\n", cmdutil.GenerateServerBrowseURL(viper.GetString("server"), issueKey))

	return nil
}

