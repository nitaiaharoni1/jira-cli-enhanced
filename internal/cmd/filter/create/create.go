package create

import (
	"github.com/spf13/cobra"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

var (
	name        string
	description string
	jql         string
	favorite    bool
)

// NewCmdCreate is a create command.
func NewCmdCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create creates a new saved filter",
		Long:    "Create creates a new saved filter with the specified name and JQL query.",
		Aliases: []string{"new"},
		Run:     Create,
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Filter name (required)")
	cmd.Flags().StringVarP(&description, "description", "d", "", "Filter description")
	cmd.Flags().StringVarP(&jql, "jql", "j", "", "JQL query (required)")
	cmd.Flags().BoolVarP(&favorite, "favorite", "f", false, "Mark filter as favorite")

	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("jql")

	return cmd
}

// Create creates a new filter.
func Create(cmd *cobra.Command, _ []string) {
	debug, err := cmd.Flags().GetBool("debug")
	cmdutil.ExitIfError(err)

	name, err := cmd.Flags().GetString("name")
	cmdutil.ExitIfError(err)

	jql, err := cmd.Flags().GetString("jql")
	cmdutil.ExitIfError(err)

	description, _ := cmd.Flags().GetString("description")
	favorite, _ := cmd.Flags().GetBool("favorite")

	s := cmdutil.Info("Creating filter...")
	defer s.Stop()

	req := &jira.CreateFilterRequest{
		Name:      name,
		JQL:       jql,
		Favourite: favorite,
	}

	if description != "" {
		req.Description = description
	}

	filter, err := api.DefaultClient(debug).CreateFilter(req)
	cmdutil.ExitIfError(err)

	cmdutil.Success("Filter created successfully: %s (ID: %s)", filter.Name, filter.ID)
}

