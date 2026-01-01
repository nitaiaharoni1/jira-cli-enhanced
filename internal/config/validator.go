package config

import (
	"fmt"
	"net/url"
	"os"

	"github.com/spf13/viper"
	"github.com/zalando/go-keyring"

	"github.com/ankitpokhrel/jira-cli/pkg/jira"
	"github.com/ankitpokhrel/jira-cli/pkg/netrc"
)

// ValidateConfig validates the configuration and returns an error if invalid.
func ValidateConfig() error {
	server := viper.GetString("server")
	if server == "" {
		return &jira.ErrValidation{
			Field:   "server",
			Message: "server URL is required",
		}
	}

	if err := validateServerURL(server); err != nil {
		return err
	}

	login := viper.GetString("login")
	if login == "" {
		return &jira.ErrValidation{
			Field:   "login",
			Message: "login is required",
		}
	}

	// Check if API token exists (via env, netrc, or keyring)
	if !hasAPIToken(server, login) {
		return &jira.ErrAuthentication{
			Reason: "API token not found. Set JIRA_API_TOKEN environment variable or configure via netrc/keyring",
		}
	}

	return nil
}

// validateServerURL validates the server URL format.
func validateServerURL(serverURL string) error {
	u, err := url.Parse(serverURL)
	if err != nil {
		return &jira.ErrValidation{
			Field:   "server",
			Message: fmt.Sprintf("invalid URL: %s", err),
		}
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return &jira.ErrValidation{
			Field:   "server",
			Message: "URL must use http:// or https://",
		}
	}

	return nil
}

// hasAPIToken checks if API token is available from any source.
func hasAPIToken(server, login string) bool {
	// Check environment variable
	if os.Getenv("JIRA_API_TOKEN") != "" {
		return true
	}

	// Check netrc
	netrcConfig, _ := netrc.Read(server, login)
	if netrcConfig != nil && netrcConfig.Password != "" {
		return true
	}

	// Check keyring
	secret, _ := keyring.Get("jira-cli", login)
	if secret != "" {
		return true
	}

	// Check config file
	if viper.GetString("api_token") != "" {
		return true
	}

	return false
}


