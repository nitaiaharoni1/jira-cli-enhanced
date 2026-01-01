package create

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

const (
	helpText = `Create creates a new project version/release.`
	examples = `$ jira release create "v1.0.0" --project PROJ --release-date "2025-01-15"
$ jira release create "v1.0.0" --project PROJ --released`
)

var (
	name        string
	description string
	project     string
	released    bool
	archived    bool
	releaseDate string
	startDate   string
)

// NewCmdCreate is a create command.
func NewCmdCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <name>",
		Short:   "Create creates a new project version",
		Long:    helpText,
		Example: examples,
		Args:    cobra.RangeArgs(0, 1),
		Run:     Create,
	}

	cmd.Flags().StringVarP(&project, "project", "p", "", "Project key or ID (required)")
	cmd.Flags().StringVarP(&description, "description", "d", "", "Version description")
	cmd.Flags().BoolVar(&released, "released", false, "Mark version as released")
	cmd.Flags().BoolVar(&archived, "archived", false, "Mark version as archived")
	cmd.Flags().StringVar(&releaseDate, "release-date", "", "Release date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&startDate, "start-date", "", "Start date (YYYY-MM-DD)")

	_ = cmd.MarkFlagRequired("project")

	return cmd
}

// Create creates a new version.
func Create(cmd *cobra.Command, args []string) {
	debug, err := cmd.Flags().GetBool("debug")
	cmdutil.ExitIfError(err)

	var name string
	if len(args) > 0 {
		name = args[0]
	}

	project, err := cmd.Flags().GetString("project")
	cmdutil.ExitIfError(err)

	description, _ := cmd.Flags().GetString("description")
	released, _ := cmd.Flags().GetBool("released")
	archived, _ := cmd.Flags().GetBool("archived")
	releaseDate, _ := cmd.Flags().GetString("release-date")
	startDate, _ := cmd.Flags().GetString("start-date")

	if name == "" {
		cmdutil.Failed("Version name is required")
		return
	}

	// Use project from config if not provided
	if project == "" {
		project = viper.GetString("project.key")
	}

	// Format dates
	if releaseDate != "" {
		releaseDate = formatDate(releaseDate)
	}
	if startDate != "" {
		startDate = formatDate(startDate)
	}

	s := cmdutil.Info("Creating version...")
	defer s.Stop()

	req := &jira.CreateVersionRequest{
		Name:        name,
		Description: description,
		Archived:    archived,
		Released:    released,
		ReleaseDate: releaseDate,
		StartDate:   startDate,
		Project:     project,
	}

	version, err := api.DefaultClient(debug).CreateVersion(req)
	cmdutil.ExitIfError(err)

	cmdutil.Success("Version created successfully: %s (ID: %s)", version.Name, version.ID)
}

func formatDate(dateStr string) string {
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

