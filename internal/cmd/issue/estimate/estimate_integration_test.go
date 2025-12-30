//go:build integration
// +build integration

package estimate

import (
	"testing"

	"github.com/ankitpokhrel/jira-cli/internal/testutil"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

func TestEstimateIntegration(t *testing.T) {
	testutil.SkipIfNotIntegration(t)

	config := testutil.GetIntegrationConfig()
	if config == nil {
		t.Fatal("Integration test config not available")
	}

	client := testutil.NewIntegrationClient()
	if client == nil {
		t.Fatal("Failed to create integration client")
	}

	// Create a test issue
	createReq := &jira.CreateRequest{
		Project:   config.Project,
		IssueType: "Task",
		Summary:   "Integration Test - Estimate",
		Body:      "Test issue for estimate integration test",
	}

	createResp, err := client.Create(createReq)
	if err != nil {
		t.Fatalf("Failed to create test issue: %v", err)
	}

	testIssueKey := createResp.Key
	t.Logf("Created test issue: %s", testIssueKey)

	// Cleanup
	defer func() {
		_ = client.DeleteV2(testIssueKey)
		t.Logf("Cleaned up test issue: %s", testIssueKey)
	}()

	// Test: Set original estimate
	t.Run("SetOriginalEstimate", func(t *testing.T) {
		editReq := &jira.EditRequest{
			OriginalEstimate: "2h",
		}
		err := client.Edit(testIssueKey, editReq)
		if err != nil {
			t.Errorf("Failed to set original estimate: %v", err)
		}
	})

	// Test: Update remaining estimate via worklog
	t.Run("UpdateRemainingEstimate", func(t *testing.T) {
		err := client.AddIssueWorklog(testIssueKey, "", "0m", "", "1h")
		if err != nil {
			t.Errorf("Failed to update remaining estimate: %v", err)
		}
	})
}

