package custom

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdcommon"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

const (
	helpText = `Custom updates custom fields for issue(s).

You can update custom fields for multiple issues at once. Custom fields must be
configured in your config file. Use the format: FIELD-NAME=VALUE

Multiple fields can be set: FIELD1=VALUE1,FIELD2=VALUE2`
	examples = `# Set a custom field for an issue
$ jira issue custom PROJ-123 story-points=5

# Set multiple custom fields
$ jira issue custom PROJ-123 story-points=5,epic-link=EPIC-1

# Set custom fields for multiple issues
$ jira issue custom PROJ-123 PROJ-456 PROJ-789 story-points=8`
)

// NewCmdCustom is a custom field command.
func NewCmdCustom() *cobra.Command {
	return &cobra.Command{
		Use:     "custom ISSUE-KEY... FIELD=VALUE [FIELD=VALUE...]",
		Short:   "Update custom fields for issue(s)",
		Long:    helpText,
		Example: examples,
		Aliases: []string{"cf"},
		Args:    cobra.MinimumNArgs(2),
		RunE:    custom,
	}
}

func custom(cmd *cobra.Command, args []string) error {
	project := viper.GetString("project.key")

	// Parse arguments: last N args are FIELD=VALUE pairs, rest are issue keys
	var issueKeys []string
	var fieldPairs []string

	for _, arg := range args {
		if strings.Contains(arg, "=") {
			fieldPairs = append(fieldPairs, arg)
		} else {
			issueKeys = append(issueKeys, arg)
		}
	}

	if len(issueKeys) == 0 {
		return fmt.Errorf("at least one issue key is required")
	}
	if len(fieldPairs) == 0 {
		return fmt.Errorf("at least one field=value pair is required")
	}

	// Parse field=value pairs
	customFields := make(map[string]string)
	for _, pair := range fieldPairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid field format: %q (expected FIELD=VALUE)", pair)
		}
		fieldName := strings.TrimSpace(parts[0])
		fieldValue := strings.TrimSpace(parts[1])
		customFields[fieldName] = fieldValue
	}

	debug, _ := cmd.Flags().GetBool("debug")
	client := api.DefaultClient(debug)

	// Get configured custom fields
	configuredFields, err := cmdcommon.GetConfiguredCustomFields()
	if err != nil {
		return fmt.Errorf("failed to get configured custom fields: %w", err)
	}

	// Validate custom fields
	if err := cmdcommon.ValidateCustomFields(customFields, configuredFields); err != nil {
		return err
	}

	// Normalize issue keys
	normalizedKeys := make([]string, 0, len(issueKeys))
	for _, key := range issueKeys {
		normalizedKeys = append(normalizedKeys, cmdutil.GetJiraIssueKey(project, key))
	}

	// Update custom fields
	fieldNames := make([]string, 0, len(customFields))
	for k := range customFields {
		fieldNames = append(fieldNames, k)
	}
	s := cmdutil.Info(fmt.Sprintf("Updating %s for %d issues...", strings.Join(fieldNames, ", "), len(normalizedKeys)))
	defer s.Stop()

	var failed []string
	var succeeded []string

	for _, key := range normalizedKeys {
		editReq := &jira.EditRequest{
			CustomFields: customFields,
		}
		editReq.WithCustomFields(configuredFields)

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
		cmdutil.Success("Successfully updated custom fields for %d issues", len(succeeded))
	}

	return nil
}

