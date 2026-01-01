package view

import (
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
$ jira issue view ISSUE-1 --output json`

	flagOutput   = "output"
	flagDebug    = "debug"
	flagComments = "comments"
	flagPlain    = "plain"

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

	return &cmd
}

func view(cmd *cobra.Command, args []string) error {
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
