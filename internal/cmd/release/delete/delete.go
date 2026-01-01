package delete

import (
	"github.com/spf13/cobra"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
)

const (
	helpText = `Delete deletes a project version.`
	examples = `$ jira release delete VERSION-ID
$ jira release delete VERSION-ID --move-fix-issues-to VERSION-ID-2`
)

var (
	moveFixIssuesTo     string
	moveAffectedIssuesTo string
)

// NewCmdDelete is a delete command.
func NewCmdDelete() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <version-id>",
		Short:   "Delete deletes a project version",
		Long:    helpText,
		Example: examples,
		Args:    cobra.ExactArgs(1),
		Run:     Delete,
	}

	cmd.Flags().StringVar(&moveFixIssuesTo, "move-fix-issues-to", "", "Move fix issues to this version ID")
	cmd.Flags().StringVar(&moveAffectedIssuesTo, "move-affected-issues-to", "", "Move affected issues to this version ID")

	return cmd
}

// Delete deletes a version.
func Delete(cmd *cobra.Command, args []string) {
	debug, err := cmd.Flags().GetBool("debug")
	cmdutil.ExitIfError(err)

	versionID := args[0]

	moveFixIssuesTo, _ := cmd.Flags().GetString("move-fix-issues-to")
	moveAffectedIssuesTo, _ := cmd.Flags().GetString("move-affected-issues-to")

	s := cmdutil.Info("Deleting version...")
	defer s.Stop()

	err = api.DefaultClient(debug).DeleteVersion(versionID, moveFixIssuesTo, moveAffectedIssuesTo)
	cmdutil.ExitIfError(err)

	cmdutil.Success("Version %s deleted successfully", versionID)
}

