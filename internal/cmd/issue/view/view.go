package view

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	tuiView "github.com/ankitpokhrel/jira-cli/internal/view"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
	"github.com/ankitpokhrel/jira-cli/pkg/jira/filter/issue"
)

const (
	helpText = `View displays contents of an issue.`
	examples = `$ jira issue view ISSUE-1

# Show 5 recent comments when viewing the issue
$ jira issue view ISSUE-1 --comments 5

# Get the raw JSON data
$ jira issue view ISSUE-1 --output json

# Extract sprint IDs from an issue
$ jira issue view ISSUE-1 --sprint-ids`

	flagOutput    = "output"
	flagDebug     = "debug"
	flagComments  = "comments"
	flagPlain     = "plain"
	flagSprintIDs = "sprint-ids"

	configProject = "project.key"
	configServer  = "server"

	messageFetchingData = "Fetching issue details..."
)

// NewCmdView is a view command.
func NewCmdView() *cobra.Command {
	cmd := cobra.Command{
		Use:     "view ISSUE-KEY",
		Short:   "View displays contents of an issue",
		Long:    helpText,
		Example: examples,
		Aliases: []string{"show"},
		Annotations: map[string]string{
			"help:args": "ISSUE-KEY\tIssue key, eg: ISSUE-1",
		},
		Args: cobra.MinimumNArgs(1),
		RunE: view,
	}

	cmd.Flags().Uint(flagComments, 1, "Show N comments")
	cmd.Flags().Bool(flagPlain, false, "Display output in plain mode")
	cmd.Flags().String(flagOutput, "", "Output format: json (default: formatted)")
	cmd.Flags().Bool(flagSprintIDs, false, "Extract and display sprint IDs only")

	return &cmd
}

func view(cmd *cobra.Command, args []string) error {
	sprintIDs, err := cmd.Flags().GetBool(flagSprintIDs)
	if err != nil {
		return err
	}

	if sprintIDs {
		return viewSprintIDs(cmd, args)
	}

	outputFormat, err := cmd.Flags().GetString(flagOutput)
	if err != nil {
		return err
	}

	if outputFormat == "json" {
		return viewRaw(cmd, args)
	}
	return viewPretty(cmd, args)
}

func viewRaw(cmd *cobra.Command, args []string) error {
	debug, err := cmd.Flags().GetBool(flagDebug)
	if err != nil {
		return err
	}

	key := cmdutil.GetJiraIssueKey(viper.GetString(configProject), args[0])

	apiResp, err := func() (string, error) {
		s := cmdutil.Info(messageFetchingData)
		defer s.Stop()

		client := api.DefaultClient(debug)
		return api.ProxyGetIssueRaw(client, key)
	}()
	if err != nil {
		return err
	}

	fmt.Println(apiResp)
	return nil
}

func viewSprintIDs(cmd *cobra.Command, args []string) error {
	debug, err := cmd.Flags().GetBool(flagDebug)
	if err != nil {
		return err
	}

	key := cmdutil.GetJiraIssueKey(viper.GetString(configProject), args[0])

	apiResp, err := func() (string, error) {
		s := cmdutil.Info(messageFetchingData)
		defer s.Stop()

		client := api.DefaultClient(debug)
		return api.ProxyGetIssueRaw(client, key)
	}()
	if err != nil {
		return err
	}

	// Parse JSON to extract sprint IDs from customfield_10020
	var issueData map[string]interface{}
	if err := json.Unmarshal([]byte(apiResp), &issueData); err != nil {
		return fmt.Errorf("failed to parse issue JSON: %w", err)
	}

	fields, ok := issueData["fields"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid issue structure: fields not found")
	}

	// Try common sprint custom field names
	sprintFields := []string{"customfield_10020", "customfield_10021"}
	var sprints []map[string]interface{}

	for _, fieldName := range sprintFields {
		if sprintData, exists := fields[fieldName]; exists {
			if sprintArray, ok := sprintData.([]interface{}); ok {
				for _, sprint := range sprintArray {
					if sprintMap, ok := sprint.(map[string]interface{}); ok {
						sprints = append(sprints, sprintMap)
					}
				}
			} else if sprintMap, ok := sprintData.(map[string]interface{}); ok {
				sprints = append(sprints, sprintMap)
			}
		}
	}

	if len(sprints) == 0 {
		fmt.Println("No sprints found")
		return nil
	}

	// Display sprint IDs and names
	for _, sprint := range sprints {
		if id, ok := sprint["id"].(float64); ok {
			name := ""
			if n, ok := sprint["name"].(string); ok {
				name = n
			}
			if name != "" {
				fmt.Printf("%.0f - %s\n", id, name)
			} else {
				fmt.Printf("%.0f\n", id)
			}
		}
	}

	return nil
}

func viewPretty(cmd *cobra.Command, args []string) error {
	debug, err := cmd.Flags().GetBool(flagDebug)
	if err != nil {
		return err
	}

	var comments uint
	if cmd.Flags().Changed(flagComments) {
		comments, err = cmd.Flags().GetUint(flagComments)
		if err != nil {
			return err
		}
	} else {
		numComments := viper.GetUint("num_comments")
		comments = max(numComments, 1)
	}

	key := cmdutil.GetJiraIssueKey(viper.GetString(configProject), args[0])
	iss, err := func() (*jira.Issue, error) {
		s := cmdutil.Info(messageFetchingData)
		defer s.Stop()

		client := api.DefaultClient(debug)
		return api.ProxyGetIssue(client, key, issue.NewNumCommentsFilter(comments))
	}()
	if err != nil {
		return err
	}

	plain, err := cmd.Flags().GetBool(flagPlain)
	if err != nil {
		return err
	}

	v := tuiView.Issue{
		Server:  viper.GetString(configServer),
		Data:    iss,
		Display: tuiView.DisplayFormat{Plain: plain},
		Options: tuiView.IssueOption{NumComments: comments},
	}
	return v.Render()
}
