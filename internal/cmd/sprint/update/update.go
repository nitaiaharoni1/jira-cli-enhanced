package update

import (
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
)

const (
	helpText = `Update updates an existing sprint.`
	examples = `$ jira sprint update 123 --name "Updated Sprint Name"
$ jira sprint update 123 --goal "New sprint goal"
$ jira sprint update 123 --start "2025-01-01" --end "2025-01-14"`
)

var (
	name      string
	startDate string
	endDate   string
	goal      string
)

// NewCmdUpdate is an update command.
func NewCmdUpdate() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <sprint-id>",
		Short:   "Update updates an existing sprint",
		Long:    helpText,
		Example: examples,
		Args:    cobra.ExactArgs(1),
		Run:     Update,
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Sprint name")
	cmd.Flags().StringVarP(&startDate, "start", "s", "", "Start date (YYYY-MM-DD)")
	cmd.Flags().StringVarP(&endDate, "end", "e", "", "End date (YYYY-MM-DD)")
	cmd.Flags().StringVarP(&goal, "goal", "g", "", "Sprint goal")

	return cmd
}

// Update updates a sprint.
func Update(cmd *cobra.Command, args []string) {
	debug, err := cmd.Flags().GetBool("debug")
	cmdutil.ExitIfError(err)

	sprintID, err := strconv.Atoi(args[0])
	cmdutil.ExitIfError(err)

	name, _ := cmd.Flags().GetString("name")
	startDate, _ := cmd.Flags().GetString("start")
	endDate, _ := cmd.Flags().GetString("end")
	goal, _ := cmd.Flags().GetString("goal")

	if name == "" && startDate == "" && endDate == "" && goal == "" {
		cmdutil.Failed("At least one field must be provided to update (--name, --start, --end, or --goal)")
		return
	}

	// Format dates if provided
	if startDate != "" {
		startDate = formatDate(startDate)
	}
	if endDate != "" {
		endDate = formatDate(endDate)
	}

	s := cmdutil.Info("Updating sprint...")
	defer s.Stop()

	sprint, err := api.DefaultClient(debug).UpdateSprint(sprintID, name, startDate, endDate, goal)
	cmdutil.ExitIfError(err)

	cmdutil.Success("Sprint updated successfully: %s (ID: %d)", sprint.Name, sprint.ID)
}

func formatDate(dateStr string) string {
	// Try parsing common date formats
	formats := []string{
		"2006-01-02",
		"2006/01/02",
		"01/02/2006",
		time.RFC3339,
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t.Format(time.RFC3339)
		}
	}

	// If parsing fails, return as-is (might already be in correct format)
	return dateStr
}

