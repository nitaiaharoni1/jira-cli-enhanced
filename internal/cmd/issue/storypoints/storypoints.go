package storypoints

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdcommon"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

const (
	helpText = `Story-points sets or updates story points for issue(s).

Story points is a custom field that must be configured in your Jira instance.
The field name may vary (e.g., "Story Points", "Story point estimate", etc.).

You can configure the custom field name in your config file.`
	examples = `# Set story points for an issue
$ jira issue story-points PROJ-123 5

# Set story points for multiple issues
$ jira issue story-points PROJ-123 PROJ-456 PROJ-789 8

# Remove story points (set to 0)
$ jira issue story-points PROJ-123 0`
)

// NewCmdStoryPoints is a story points command.
func NewCmdStoryPoints() *cobra.Command {
	cmd := cobra.Command{
		Use:     "story-points ISSUE-KEY... POINTS",
		Short:   "Set or update story points for issue(s)",
		Long:    helpText,
		Example: examples,
		Aliases: []string{"sp", "points"},
		Args:    cobra.MinimumNArgs(2),
		RunE:    storyPoints,
	}

	cmd.Flags().String("field", "", "Custom field name for story points (overrides config)")

	return &cmd
}

func storyPoints(cmd *cobra.Command, args []string) error {
	project := viper.GetString("project.key")

	// Last argument is the story points value
	pointsStr := args[len(args)-1]
	issueKeys := args[:len(args)-1]

	_, err := strconv.ParseFloat(pointsStr, 64)
	if err != nil {
		return fmt.Errorf("invalid story points value: %q (must be a number)", pointsStr)
	}

	fieldName, _ := cmd.Flags().GetString("field")
	debug, _ := cmd.Flags().GetBool("debug")

	client := api.DefaultClient(debug)

	// Get configured custom fields
	configuredFields, err := cmdcommon.GetConfiguredCustomFields()
	if err != nil {
		return fmt.Errorf("failed to get configured custom fields: %w", err)
	}

	// Find story points field
	var storyPointsField *jira.IssueTypeField
	if fieldName != "" {
		// Use provided field name
		for _, f := range configuredFields {
			if strings.EqualFold(f.Name, fieldName) {
				storyPointsField = &f
				break
			}
		}
		if storyPointsField == nil {
			return fmt.Errorf("custom field %q not found in configuration", fieldName)
		}
	} else {
		// Try to find story points field automatically
		storyPointsKeywords := []string{"story point", "storypoint", "story-point"}
		for _, f := range configuredFields {
			nameLower := strings.ToLower(f.Name)
			for _, keyword := range storyPointsKeywords {
				if strings.Contains(nameLower, keyword) {
					storyPointsField = &f
					break
				}
			}
			if storyPointsField != nil {
				break
			}
		}
	}

	if storyPointsField == nil {
		return fmt.Errorf("story points field not found. Configure it in your config file or use --field flag")
	}

	// Normalize issue keys
	normalizedKeys := make([]string, 0, len(issueKeys))
	for _, key := range issueKeys {
		normalizedKeys = append(normalizedKeys, cmdutil.GetJiraIssueKey(project, key))
	}

	// Update story points
	s := cmdutil.Info(fmt.Sprintf("Setting story points to %s for %d issues...", pointsStr, len(normalizedKeys)))
	defer s.Stop()

	var failed []string
	var succeeded []string

	for _, key := range normalizedKeys {
		customFields := map[string]string{
			storyPointsField.Key: pointsStr,
		}

		editReq := &jira.EditRequest{
			CustomFields: customFields,
		}
		editReq.WithCustomFields(configuredFields)

		if err := cmdcommon.ValidateCustomFields(customFields, configuredFields); err != nil {
			failed = append(failed, key)
			continue
		}

		err := client.Edit(key, editReq)
		if err != nil {
			failed = append(failed, key)
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
			return fmt.Errorf("failed to update all issues")
		}
	} else {
		cmdutil.Success("Successfully set story points to %s for %d issues", pointsStr, len(succeeded))
	}

	return nil
}

