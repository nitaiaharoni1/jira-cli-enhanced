package update

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

const (
	helpText = `Update updates an existing project version.`
	examples = `$ jira release update VERSION-ID --name "v1.0.1"
$ jira release update VERSION-ID --released
$ jira release update VERSION-ID --archived`
)

var (
	name        string
	description string
	released    *bool
	archived    *bool
	releaseDate string
	startDate   string
)

// NewCmdUpdate is an update command.
func NewCmdUpdate() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <version-id>",
		Short:   "Update updates an existing project version",
		Long:    helpText,
		Example: examples,
		Args:    cobra.ExactArgs(1),
		Run:     Update,
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Version name")
	cmd.Flags().StringVarP(&description, "description", "d", "", "Version description")
	cmd.Flags().BoolVar(&released, "released", false, "Mark version as released")
	cmd.Flags().BoolVar(&archived, "archived", false, "Mark version as archived")
	cmd.Flags().StringVar(&releaseDate, "release-date", "", "Release date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&startDate, "start-date", "", "Start date (YYYY-MM-DD)")

	return cmd
}

// Update updates a version.
func Update(cmd *cobra.Command, args []string) {
	debug, err := cmd.Flags().GetBool("debug")
	cmdutil.ExitIfError(err)

	versionID := args[0]

	name, _ := cmd.Flags().GetString("name")
	description, _ := cmd.Flags().GetString("description")
	releasedFlag, _ := cmd.Flags().GetBool("released")
	archivedFlag, _ := cmd.Flags().GetBool("archived")
	releaseDate, _ := cmd.Flags().GetString("release-date")
	startDate, _ := cmd.Flags().GetString("start-date")

	if name == "" && description == "" && !releasedFlag && !archivedFlag && releaseDate == "" && startDate == "" {
		cmdutil.Failed("At least one field must be provided to update")
		return
	}

	req := &jira.UpdateVersionRequest{
		Name:        name,
		Description: description,
		ReleaseDate: formatDate(releaseDate),
		StartDate:   formatDate(startDate),
	}

	if cmd.Flags().Changed("released") {
		req.Released = &releasedFlag
	}
	if cmd.Flags().Changed("archived") {
		req.Archived = &archivedFlag
	}

	s := cmdutil.Info("Updating version...")
	defer s.Stop()

	version, err := api.DefaultClient(debug).UpdateVersion(versionID, req)
	cmdutil.ExitIfError(err)

	cmdutil.Success("Version updated successfully: %s (ID: %s)", version.Name, version.ID)
}

func formatDate(dateStr string) string {
	if dateStr == "" {
		return ""
	}

	formats := []string{
		"2006-01-02",
		"2006/01/02",
		time.RFC3339,
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t.Format("2006-01-02")
		}
	}

	return dateStr
}

