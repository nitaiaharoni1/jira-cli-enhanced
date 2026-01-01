package create

import (
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
)

const (
	helpText = `Create creates a new sprint in a board.`
	examples = `$ jira sprint create "Sprint 1" --board 123 --start "2025-01-01" --end "2025-01-14"
$ jira sprint create "Sprint 2" --board 123 --start "2025-01-15" --end "2025-01-28" --goal "Complete feature X"`
)

var (
	name      string
	boardID   int
	startDate string
	endDate   string
	goal      string
)

// NewCmdCreate is a create command.
func NewCmdCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <name>",
		Short:   "Create creates a new sprint",
		Long:    helpText,
		Example: examples,
		Args:    cobra.RangeArgs(0, 1),
		Run:     Create,
	}

	cmd.Flags().IntVarP(&boardID, "board", "b", 0, "Board ID (required)")
	cmd.Flags().StringVarP(&startDate, "start", "s", "", "Start date (YYYY-MM-DD)")
	cmd.Flags().StringVarP(&endDate, "end", "e", "", "End date (YYYY-MM-DD)")
	cmd.Flags().StringVarP(&goal, "goal", "g", "", "Sprint goal")

	_ = cmd.MarkFlagRequired("board")

	return cmd
}

// Create creates a new sprint.
func Create(cmd *cobra.Command, args []string) {
	debug, err := cmd.Flags().GetBool("debug")
	cmdutil.ExitIfError(err)

	boardID, err := cmd.Flags().GetInt("board")
	cmdutil.ExitIfError(err)

	var name string
	if len(args) > 0 {
		name = args[0]
	}

	startDate, _ := cmd.Flags().GetString("start")
	endDate, _ := cmd.Flags().GetString("end")
	goal, _ := cmd.Flags().GetString("goal")

	// Interactive mode if name not provided
	if name == "" {
		prompt := &survey.Input{
			Message: "Sprint name",
		}
		err := survey.AskOne(prompt, &name, survey.WithValidator(survey.Required))
		cmdutil.ExitIfError(err)
	}

	// Format dates if provided
	if startDate != "" {
		startDate = formatDate(startDate)
	}
	if endDate != "" {
		endDate = formatDate(endDate)
	}

	s := cmdutil.Info("Creating sprint...")
	defer s.Stop()

	sprint, err := api.DefaultClient(debug).CreateSprint(boardID, name, startDate, endDate, goal)
	cmdutil.ExitIfError(err)

	cmdutil.Success("Sprint created successfully: %s (ID: %d)", sprint.Name, sprint.ID)
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

