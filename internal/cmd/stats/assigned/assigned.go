package assigned

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

const (
	helpText = `Assigned displays issue distribution by status for assigned issues.`
	examples = `# Get issue distribution for current user
$ jira stats assigned

# Get issue distribution for a specific user
$ jira stats assigned --user "john@example.com"

# Get issue distribution with custom JQL
$ jira stats assigned --jql "project = PROJ AND assignee = currentUser()"`
)

// NewCmdAssigned is an assigned stats command.
func NewCmdAssigned() *cobra.Command {
	cmd := cobra.Command{
		Use:     "assigned",
		Short:   "Display issue distribution by status",
		Long:    helpText,
		Example: examples,
		Aliases: []string{"dist", "d"},
		RunE:    assignedStats,
	}

	cmd.Flags().String("user", "", "User email or display name (defaults to current user)")
	cmd.Flags().String("jql", "", "Custom JQL query")

	return &cmd
}

func assignedStats(cmd *cobra.Command, args []string) error {
	debug, _ := cmd.Flags().GetBool("debug")
	userFlag, _ := cmd.Flags().GetString("user")
	jqlFlag, _ := cmd.Flags().GetString("jql")

	client := api.DefaultClient(debug)

	var jql string
	if jqlFlag != "" {
		jql = jqlFlag
	} else {
		// Build JQL for assigned issues
		if userFlag != "" {
			jql = fmt.Sprintf("assignee = %q", userFlag)
		} else {
			me, err := api.ProxyMe(client)
			if err != nil {
				return fmt.Errorf("failed to get current user: %w", err)
			}
			assignee := me.Login
			if assignee == "" {
				assignee = me.Email
			}
			if assignee == "" {
				assignee = me.Name
			}
			jql = fmt.Sprintf("assignee = %q", assignee)
		}
	}

	s := cmdutil.Info("Fetching issue distribution...")
	defer s.Stop()

	dist, err := client.GetIssueDistribution(jql)
	if err != nil {
		return fmt.Errorf("failed to get issue distribution: %w", err)
	}

	s.Stop()

	if len(dist) == 0 {
		fmt.Println("\nNo issues found.")
		return nil
	}

	// Sort by count descending
	sort.Slice(dist, func(i, j int) bool {
		return dist[i].Count > dist[j].Count
	})

	// Display distribution
	fmt.Println("\nIssue Distribution by Status")
	fmt.Println(strings.Repeat("─", 40))
	total := 0
	for _, d := range dist {
		fmt.Printf("%-20s %d\n", d.Status+":", d.Count)
		total += d.Count
	}
	fmt.Println(strings.Repeat("─", 40))
	fmt.Printf("%-20s %d\n", "Total:", total)
	fmt.Println()

	return nil
}

