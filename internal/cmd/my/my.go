package my

import (
	"github.com/spf13/cobra"

	"github.com/ankitpokhrel/jira-cli/internal/cmd/issue/list"
)

const helpText = `My displays issues assigned to the current user.`

// NewCmdMy is a my command.
func NewCmdMy() *cobra.Command {
	cmd := cobra.Command{
		Use:         "my",
		Short:       "Display issues assigned to you",
		Long:        helpText,
		Aliases:     []string{"mine"},
		Annotations: map[string]string{"cmd:main": "true"},
		Run:         myIssues,
	}

	// Add all flags from issue list
	listCmd := list.NewCmdList()
	list.SetFlags(listCmd)
	cmd.Flags().AddFlagSet(listCmd.Flags())

	// Add shortcuts
	cmd.Flags().Bool("todo", false, "Show only 'To Do' issues")
	cmd.Flags().Bool("done", false, "Show only 'Done' issues")
	cmd.Flags().Bool("blocked", false, "Show only blocked issues")

	return &cmd
}

func myIssues(cmd *cobra.Command, args []string) {
	// Build JQL for assigned issues
	jql := "assignee = currentUser()"
	
	// Add status filters
	if todo, _ := cmd.Flags().GetBool("todo"); todo {
		jql += " AND status = \"To Do\""
	} else if done, _ := cmd.Flags().GetBool("done"); done {
		jql += " AND status = \"Done\""
	}
	
	if blocked, _ := cmd.Flags().GetBool("blocked"); blocked {
		jql += " AND status = \"Blocked\""
	}
	
	// Set JQL flag
	cmd.Flags().Set("jql", jql)
	
	// Call list's LoadList function directly
	list.LoadList(cmd, args)
}

