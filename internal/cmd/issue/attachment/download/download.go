package download

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
)

const (
	helpText = `Download downloads an attachment from an issue.`
	examples = `# Download attachment
$ jira issue attachment download PROJ-123 ATTACHMENT-ID

# Download to specific file
$ jira issue attachment download PROJ-123 ATTACHMENT-ID -o downloaded.pdf`
)

// NewCmdDownload is a download command.
func NewCmdDownload() *cobra.Command {
	cmd := cobra.Command{
		Use:     "download ISSUE-KEY ATTACHMENT-ID",
		Short:   "Download an attachment from an issue",
		Long:    helpText,
		Example: examples,
		Args:    cobra.ExactArgs(2),
		RunE:    download,
	}

	cmd.Flags().StringP("output", "o", "", "Output file path (defaults to attachment filename)")

	return &cmd
}

func download(cmd *cobra.Command, args []string) error {
	project := viper.GetString("project.key")
	issueKey := cmdutil.GetJiraIssueKey(project, args[0])
	attachmentID := args[1]

	output, _ := cmd.Flags().GetString("output")
	if output == "" {
		// Get attachment info to get filename
		debug, _ := cmd.Flags().GetBool("debug")
		client := api.DefaultClient(debug)
		attachments, err := client.GetAttachments(issueKey)
		if err != nil {
			return fmt.Errorf("failed to get attachment info: %w", err)
		}

		for _, att := range attachments {
			if att.ID == attachmentID {
				output = att.Filename
				break
			}
		}

		if output == "" {
			return fmt.Errorf("attachment %s not found", attachmentID)
		}
	}

	debug, _ := cmd.Flags().GetBool("debug")
	client := api.DefaultClient(debug)

	s := cmdutil.Info(fmt.Sprintf("Downloading attachment %s...", attachmentID))
	err := client.DownloadAttachment(attachmentID, output)
	s.Stop()

	if err != nil {
		return err
	}

	cmdutil.Success("Downloaded attachment to %s", output)
	return nil
}

