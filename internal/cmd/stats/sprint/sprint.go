package sprint

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
)

const (
	helpText = `Sprint displays statistics for a sprint including completion percentage and issue distribution.`
	examples = `# Get statistics for a sprint
$ jira stats sprint 123

# Get statistics for current sprint (if board is configured)
$ jira stats sprint`
)

// NewCmdSprint is a sprint stats command.
func NewCmdSprint() *cobra.Command {
	cmd := cobra.Command{
		Use:     "sprint [SPRINT-ID]",
		Short:   "Display sprint statistics",
		Long:    helpText,
		Example: examples,
		Aliases: []string{"s"},
		Args:    cobra.MaximumNArgs(1),
		RunE:    sprintStats,
	}

	return &cmd
}

func sprintStats(cmd *cobra.Command, args []string) error {
	debug, _ := cmd.Flags().GetBool("debug")
	client := api.DefaultClient(debug)

	var sprintID int
	var err error

	if len(args) > 0 {
		sprintID, err = strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid sprint ID: %s", args[0])
		}
	} else {
		// Try to get current sprint from board
		boardID := viper.GetInt("board.id")
		if boardID == 0 {
			return fmt.Errorf("sprint ID required or configure board.id in config")
		}

		// Get active sprints
		sprints, err := client.Sprints(boardID, "state=active", 0, 10)
		if err != nil {
			return fmt.Errorf("failed to get sprints: %w", err)
		}

		if len(sprints.Sprints) == 0 {
			return fmt.Errorf("no active sprint found. Please specify sprint ID")
		}

		sprintID = sprints.Sprints[0].ID
	}

	s := cmdutil.Info(fmt.Sprintf("Fetching sprint statistics for sprint %d...", sprintID))
	defer s.Stop()

	stats, err := client.GetSprintStatistics(sprintID)
	if err != nil {
		return fmt.Errorf("failed to get sprint statistics: %w", err)
	}

	s.Stop()

	// Display statistics
	fmt.Printf("\nSprint: %s\n", stats.SprintName)
	fmt.Println(strings.Repeat("─", 40))
	fmt.Printf("Total Issues:     %d\n", stats.TotalIssues)
	fmt.Printf("Completed:        %d (%.0f%%)\n", stats.Completed, stats.CompletionPct)
	fmt.Printf("In Progress:      %d\n", stats.InProgress)
	fmt.Printf("To Do:            %d\n", stats.ToDo)
	if stats.StoryPoints > 0 {
		fmt.Printf("Story Points:     %d/%d (%.0f%%)\n", stats.CompletedSP, stats.StoryPoints, stats.VelocityPct)
	}
	fmt.Println(strings.Repeat("─", 40))
	fmt.Println()

	return nil
}

