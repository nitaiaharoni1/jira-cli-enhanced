package estimate

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

const (
	helpText = `Estimate sets or updates the time estimate for an issue.

You can set the original estimate or remaining estimate. The estimate format
follows Jira's time tracking format: "2d 3h 30m", "10m", "1w", etc.`
	examples = `# Set original estimate
$ jira issue estimate PROJ-123 "2d 3h"

# Update remaining estimate
$ jira issue estimate PROJ-123 "1d" --remaining

# Set estimate for multiple issues
$ jira issue estimate PROJ-123 PROJ-456 PROJ-789 "3d"`
)

// NewCmdEstimate is an estimate command.
func NewCmdEstimate() *cobra.Command {
	cmd := cobra.Command{
		Use:     "estimate ISSUE-KEY... ESTIMATE",
		Short:   "Set or update time estimate for issue(s)",
		Long:    helpText,
		Example: examples,
		Aliases: []string{"est"},
		Args:    cobra.MinimumNArgs(2),
		RunE:    estimate,
	}

	cmd.Flags().Bool("remaining", false, "Update remaining estimate instead of original")

	return &cmd
}

func estimate(cmd *cobra.Command, args []string) error {
	project := viper.GetString("project.key")

	// Last argument is the estimate
	estimateValue := args[len(args)-1]
	issueKeys := args[:len(args)-1]

	remaining, _ := cmd.Flags().GetBool("remaining")
	debug, _ := cmd.Flags().GetBool("debug")

	client := api.DefaultClient(debug)

	// Normalize issue keys
	normalizedKeys := make([]string, 0, len(issueKeys))
	for _, key := range issueKeys {
		normalizedKeys = append(normalizedKeys, cmdutil.GetJiraIssueKey(project, key))
	}

	// Update estimates
	var action string
	if remaining {
		action = "remaining estimate"
	} else {
		action = "original estimate"
	}

	s := cmdutil.Info(fmt.Sprintf("Setting %s for %d issues...", action, len(normalizedKeys)))
	defer s.Stop()

	var failed []string
	var succeeded []string

	for _, key := range normalizedKeys {
		var err error
		if remaining {
			// Update remaining estimate via worklog (with 0 time spent)
			err = client.AddIssueWorklog(key, "", "0m", "", estimateValue)
		} else {
			// Update original estimate via edit API
			editReq := &jira.EditRequest{
				OriginalEstimate: estimateValue,
			}
			err = client.Edit(key, editReq)
		}

		if err != nil {
			failed = append(failed, fmt.Sprintf("%s (%v)", key, err))
			continue
		}
		succeeded = append(succeeded, key)
	}

	s.Stop()

	if len(failed) > 0 {
		if len(succeeded) > 0 {
			cmdutil.Warn("Updated %d issues successfully, %d failed", len(succeeded), len(failed))
			fmt.Printf("Failed: %s\n", strings.Join(failed, ", "))
		} else {
			return fmt.Errorf("failed to update all issues: %s", strings.Join(failed, ", "))
		}
	} else {
		cmdutil.Success("Successfully updated %s for %d issues", action, len(succeeded))
	}

	return nil
}

