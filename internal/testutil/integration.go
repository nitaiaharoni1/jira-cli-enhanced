package testutil

import (
	"os"
	"testing"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

// IntegrationTestConfig holds configuration for integration tests.
type IntegrationTestConfig struct {
	Server   string
	Login    string
	APIToken string
	Project  string
}

// GetIntegrationConfig retrieves integration test configuration from environment variables.
// Returns nil if integration tests are not enabled.
func GetIntegrationConfig() *IntegrationTestConfig {
	if os.Getenv("JIRA_INTEGRATION_TEST") != "true" {
		return nil
	}

	server := os.Getenv("JIRA_TEST_SERVER")
	login := os.Getenv("JIRA_TEST_LOGIN")
	token := os.Getenv("JIRA_TEST_API_TOKEN")
	project := os.Getenv("JIRA_TEST_PROJECT")

	if server == "" || login == "" || token == "" {
		return nil
	}

	return &IntegrationTestConfig{
		Server:   server,
		Login:    login,
		APIToken: token,
		Project:  project,
	}
}

// SkipIfNotIntegration skips the test if integration tests are not enabled.
func SkipIfNotIntegration(t *testing.T) {
	if GetIntegrationConfig() == nil {
		t.Skip("Skipping integration test. Set JIRA_INTEGRATION_TEST=true and required env vars to run.")
	}
}

// NewIntegrationClient creates a Jira client for integration testing.
func NewIntegrationClient() *jira.Client {
	config := GetIntegrationConfig()
	if config == nil {
		return nil
	}

	return jira.NewClient(jira.Config{
		Server:   config.Server,
		Login:    config.Login,
		APIToken: config.APIToken,
		Debug:    os.Getenv("JIRA_TEST_DEBUG") == "true",
	})
}

// NewIntegrationAPIClient creates an API client for integration testing.
func NewIntegrationAPIClient() *jira.Client {
	config := GetIntegrationConfig()
	if config == nil {
		return nil
	}

	return api.Client(jira.Config{
		Server:   config.Server,
		Login:    config.Login,
		APIToken: config.APIToken,
		Debug:    os.Getenv("JIRA_TEST_DEBUG") == "true",
	})
}

