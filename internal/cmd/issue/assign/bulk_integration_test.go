//go:build integration
// +build integration

package assign

import (
	"testing"

	"github.com/ankitpokhrel/jira-cli/internal/testutil"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

func TestAssignBulkIntegration(t *testing.T) {
	testutil.SkipIfNotIntegration(t)

	config := testutil.GetIntegrationConfig()
	if config == nil {
		t.Fatal("Integration test config not available")
	}

	client := testutil.NewIntegrationClient()
	if client == nil {
		t.Fatal("Failed to create integration client")
	}

	// Create test issues
	var issueKeys []string
	for i := 1; i <= 3; i++ {
		createReq := &jira.CreateRequest{
			Project:   config.Project,
			IssueType: "Task",
			Summary:   "Integration Test - Bulk Assign",
			Body:      "Test issue for bulk assign integration test",
		}

		createResp, err := client.Create(createReq)
		if err != nil {
			t.Fatalf("Failed to create test issue %d: %v", i, err)
		}
		issueKeys = append(issueKeys, createResp.Key)
	}

	// Cleanup
	defer func() {
		for _, key := range issueKeys {
			_ = client.DeleteV2(key)
		}
	}()

	t.Logf("Created test issues: %v", issueKeys)

	// Test: Bulk assign to user
	t.Run("BulkAssignToUser", func(t *testing.T) {
		// Get current user
		me, err := client.Me()
		if err != nil {
			t.Fatalf("Failed to get current user: %v", err)
		}

		assignee := me.Login
		if assignee == "" {
			assignee = me.Name
		}

		// Assign all issues
		for _, key := range issueKeys {
			err := client.AssignIssue(key, assignee)
			if err != nil {
				t.Errorf("Failed to assign issue %s: %v", key, err)
			}
		}
	})

	// Test: Bulk unassign
	t.Run("BulkUnassign", func(t *testing.T) {
		for _, key := range issueKeys {
			err := client.AssignIssue(key, jira.AssigneeNone)
			if err != nil {
				t.Errorf("Failed to unassign issue %s: %v", key, err)
			}
		}
	})
}

