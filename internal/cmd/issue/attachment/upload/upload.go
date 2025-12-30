package upload

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
)

const (
	helpText = `Upload uploads a file as an attachment to an issue.`
	examples = `# Upload a file
$ jira issue attachment upload PROJ-123 file.pdf

# Upload multiple files
$ jira issue attachment upload PROJ-123 file1.pdf file2.pdf file3.pdf`
)

// NewCmdUpload is an upload command.
func NewCmdUpload() *cobra.Command {
	return &cobra.Command{
		Use:     "upload ISSUE-KEY FILE...",
		Short:   "Upload file(s) as attachment(s) to an issue",
		Long:    helpText,
		Example: examples,
		Args:    cobra.MinimumNArgs(2),
		RunE:    upload,
	}
}

func upload(cmd *cobra.Command, args []string) error {
	project := viper.GetString("project.key")
	issueKey := cmdutil.GetJiraIssueKey(project, args[0])
	files := args[1:]

	debug, _ := cmd.Flags().GetBool("debug")
	client := api.DefaultClient(debug)

	var uploaded []string
	var failed []string

	for _, filePath := range files {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			failed = append(failed, fmt.Sprintf("%s (file not found)", filePath))
			continue
		}

		s := cmdutil.Info(fmt.Sprintf("Uploading %s...", filePath))
		attachments, err := client.UploadAttachment(issueKey, filePath)
		s.Stop()

		if err != nil {
			failed = append(failed, fmt.Sprintf("%s (%v)", filePath, err))
			continue
		}

		for _, att := range attachments {
			uploaded = append(uploaded, att.Filename)
		}
	}

	if len(failed) > 0 {
		if len(uploaded) > 0 {
			cmdutil.Warn("Uploaded %d file(s) successfully, %d failed", len(uploaded), len(failed))
			fmt.Printf("Failed: %s\n", failed)
		} else {
			return fmt.Errorf("failed to upload all files: %v", failed)
		}
	} else {
		cmdutil.Success("Successfully uploaded %d file(s) to issue %s", len(uploaded), issueKey)
	}

	return nil
}

