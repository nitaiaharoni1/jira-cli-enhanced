package velocity

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
)

const (
	helpText = `Velocity displays velocity trends across multiple sprints.`
	examples = `# Show velocity for last 5 sprints
$ jira stats velocity --sprints 5

# Show velocity for specific board
$ jira stats velocity --sprints 10 --board 123`
)

// NewCmdVelocity is a velocity stats command.
func NewCmdVelocity() *cobra.Command {
	cmd := cobra.Command{
		Use:     "velocity",
		Short:   "Display velocity trends",
		Long:    helpText,
		Example: examples,
		Aliases: []string{"vel", "v"},
		RunE:    velocityStats,
	}

	cmd.Flags().Int("sprints", 5, "Number of sprints to analyze")
	cmd.Flags().Int("board", 0, "Board ID (defaults to configured board)")

	return &cmd
}

func velocityStats(cmd *cobra.Command, args []string) error {
	debug, _ := cmd.Flags().GetBool("debug")
	numSprints, _ := cmd.Flags().GetInt("sprints")
	boardID, _ := cmd.Flags().GetInt("board")

	if boardID == 0 {
		boardID = viper.GetInt("board.id")
		if boardID == 0 {
			return fmt.Errorf("board ID required or configure board.id in config")
		}
	}

	client := api.DefaultClient(debug)

	s := cmdutil.Info(fmt.Sprintf("Fetching velocity data for %d sprints...", numSprints))
	defer s.Stop()

	// Get closed sprints
	sprints, err := client.Sprints(boardID, "state=closed", 0, numSprints)
	if err != nil {
		return fmt.Errorf("failed to get sprints: %w", err)
	}

	if len(sprints.Sprints) == 0 {
		return fmt.Errorf("no closed sprints found")
	}

	// Reverse to show oldest first
	sprintList := sprints.Sprints
	for i, j := 0, len(sprintList)-1; i < j; i, j = i+1, j-1 {
		sprintList[i], sprintList[j] = sprintList[j], sprintList[i]
	}

	s.Stop()

	// Display velocity table
	fmt.Println("\nSprint          | Points | Completed | Velocity")
	fmt.Println(strings.Repeat("─", 60))

	totalPoints := 0
	totalCompleted := 0

	for _, sprint := range sprintList {
		if len(sprintList) > numSprints {
			break
		}

		stats, err := client.GetSprintStatistics(sprint.ID)
		if err != nil {
			continue
		}

		velocity := 0.0
		if stats.StoryPoints > 0 {
			velocity = stats.VelocityPct
		} else if stats.TotalIssues > 0 {
			velocity = stats.CompletionPct
		}

		totalPoints += stats.StoryPoints
		totalCompleted += stats.CompletedSP

		sprintName := sprint.Name
		if len(sprintName) > 14 {
			sprintName = sprintName[:11] + "..."
		}

		if stats.StoryPoints > 0 {
			fmt.Printf("%-15s | %6d | %9d | %.0f%%\n", sprintName, stats.StoryPoints, stats.CompletedSP, velocity)
		} else {
			fmt.Printf("%-15s | %6d | %9d | %.0f%%\n", sprintName, stats.TotalIssues, stats.Completed, velocity)
		}
	}

	fmt.Println(strings.Repeat("─", 60))
	avgVelocity := 0.0
	if totalPoints > 0 {
		avgVelocity = (float64(totalCompleted) / float64(totalPoints)) * 100
		fmt.Printf("%-15s | %6d | %9d | %.0f%%\n", "Average", totalPoints/len(sprintList), totalCompleted/len(sprintList), avgVelocity)
	}
	fmt.Println()

	return nil
}

