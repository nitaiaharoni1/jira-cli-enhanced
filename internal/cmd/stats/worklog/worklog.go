package worklog

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

const (
	helpText = `Worklog displays worklog summary for a user.`
	examples = `# Get worklog summary for current user
$ jira stats worklog

# Get worklog summary for a specific user
$ jira stats worklog --user "john@example.com"

# Get worklog summary for date range
$ jira stats worklog --from "2025-01-01" --to "2025-01-31"`
)

// NewCmdWorklog is a worklog stats command.
func NewCmdWorklog() *cobra.Command {
	cmd := cobra.Command{
		Use:     "worklog",
		Short:   "Display worklog summary",
		Long:    helpText,
		Example: examples,
		Aliases: []string{"wl", "time"},
		RunE:    worklogStats,
	}

	cmd.Flags().String("user", "", "User email or display name (defaults to current user)")
	cmd.Flags().String("from", "", "Start date (YYYY-MM-DD, defaults to start of current month)")
	cmd.Flags().String("to", "", "End date (YYYY-MM-DD, defaults to today)")

	return &cmd
}

func worklogStats(cmd *cobra.Command, args []string) error {
	debug, _ := cmd.Flags().GetBool("debug")
	userFlag, _ := cmd.Flags().GetString("user")
	fromFlag, _ := cmd.Flags().GetString("from")
	toFlag, _ := cmd.Flags().GetString("to")

	client := api.DefaultClient(debug)

	// Get current user if not specified
	user := userFlag
	if user == "" {
		me, err := api.ProxyMe(client)
		if err != nil {
			return fmt.Errorf("failed to get current user: %w", err)
		}
		user = me.Login
		if user == "" {
			user = me.Email
		}
		if user == "" {
			user = me.Name
		}
	}

	// Parse dates
	now := time.Now()
	from := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	to := now

	if fromFlag != "" {
		parsed, err := time.Parse("2006-01-02", fromFlag)
		if err != nil {
			return fmt.Errorf("invalid from date format (use YYYY-MM-DD): %w", err)
		}
		from = parsed
	}

	if toFlag != "" {
		parsed, err := time.Parse("2006-01-02", toFlag)
		if err != nil {
			return fmt.Errorf("invalid to date format (use YYYY-MM-DD): %w", err)
		}
		to = parsed
	}

	s := cmdutil.Info(fmt.Sprintf("Fetching worklog summary for %s...", user))
	defer s.Stop()

	summary, err := client.GetUserWorklogs(user, from, to)
	if err != nil {
		return fmt.Errorf("failed to get worklog summary: %w", err)
	}

	s.Stop()

	// Display summary
	fmt.Printf("\nWorklog Summary: %s\n", summary.User)
	fmt.Println(strings.Repeat("─", 50))
	fmt.Printf("Date Range:      %s\n", summary.DateRange)
	fmt.Printf("Total Hours:     %.2f\n", summary.TotalHours)
	fmt.Printf("Total Days:       %.2f\n", summary.TotalDays)
	fmt.Printf("Entries:         %d\n", summary.EntryCount)
	fmt.Printf("Issues:          %d\n", len(summary.Issues))
	if len(summary.Issues) > 0 {
		fmt.Printf("\nIssue Keys:\n")
		for _, key := range summary.Issues {
			fmt.Printf("  - %s\n", key)
		}
	}
	fmt.Println()

	return nil
}

