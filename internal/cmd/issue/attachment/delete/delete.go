package delete

import (
	"github.com/spf13/cobra"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
)

const (
	helpText = `Delete deletes an attachment from an issue.`
	examples = `# Delete attachment
$ jira issue attachment delete ATTACHMENT-ID`
)

// NewCmdDelete is a delete command.
func NewCmdDelete() *cobra.Command {
	return &cobra.Command{
		Use:     "delete ATTACHMENT-ID",
		Short:   "Delete an attachment",
		Long:    helpText,
		Example: examples,
		Args:    cobra.ExactArgs(1),
		RunE:    delete,
	}
}

func delete(cmd *cobra.Command, args []string) error {
	attachmentID := args[0]

	debug, _ := cmd.Flags().GetBool("debug")
	client := api.DefaultClient(debug)

	s := cmdutil.Info("Deleting attachment...")
	err := client.DeleteAttachment(attachmentID)
	s.Stop()

	if err != nil {
		return err
	}

	cmdutil.Success("Attachment deleted successfully")
	return nil
}

