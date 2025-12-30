package main

import (
	"github.com/ankitpokhrel/jira-cli/internal/cmd/root"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
)

func main() {
	rootCmd := root.NewCmdRoot()
	if _, err := rootCmd.ExecuteC(); err != nil {
		// Use enhanced error handling that provides suggestions
		cmdutil.ExitIfError(err)
	}
}
