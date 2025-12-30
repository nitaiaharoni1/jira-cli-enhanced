package cmdutil

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"golang.org/x/term"

	"github.com/ankitpokhrel/jira-cli/pkg/browser"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
	"github.com/ankitpokhrel/jira-cli/pkg/tui"
)

// ExitIfError exists with error message if err is not nil.
// It provides actionable suggestions based on error type.
func ExitIfError(err error) {
	if err == nil {
		return
	}

	var msg string
	var suggestion string

	switch e := err.(type) {
	case *jira.ErrAuthentication:
		msg = fmt.Sprintf("Authentication failed: %s", e.Reason)
		suggestion = "Run 'jira init' to reconfigure your credentials or check your JIRA_API_TOKEN environment variable"

	case *jira.ErrNotFound:
		msg = e.Error()
		suggestion = "Verify the ID is correct and you have access to this resource"

	case *jira.ErrValidation:
		msg = e.Error()
		suggestion = "Check the command syntax and required parameters"

	case *jira.ErrRateLimit:
		msg = e.Error()
		if e.RetryAfter > 0 {
			suggestion = fmt.Sprintf("Wait %d seconds before retrying", e.RetryAfter)
		} else {
			suggestion = "Wait a moment and try again"
		}

	case *jira.ErrNetwork:
		msg = e.Error()
		suggestion = "Check your internet connection and try again"

	case *jira.ErrUnexpectedResponse:
		dm := fmt.Sprintf(
			"\njira: Received unexpected response '%s'.\nPlease check the parameters you supplied and try again.",
			e.Status,
		)
		bd := e.Error()

		msg = dm
		if len(bd) > 0 {
			msg = fmt.Sprintf("%s%s", bd, dm)
		}

		// Provide specific suggestions based on status code
		switch {
		case e.StatusCode == 401:
			suggestion = "Authentication failed. Run 'jira init' to reconfigure credentials"
		case e.StatusCode == 403:
			suggestion = "You don't have permission to perform this operation. Check your access rights"
		case e.StatusCode == 404:
			suggestion = "Resource not found. Verify the ID is correct"
		case e.StatusCode >= 500:
			suggestion = "Server error. This may be temporary - try again in a moment"
		default:
			suggestion = "Check your parameters and try again"
		}

	case *jira.ErrMultipleFailed:
		msg = fmt.Sprintf("\n%s%s", "SOME REQUESTS REPORTED ERROR:", e.Error())
		suggestion = "Some operations failed. Review the errors above and retry failed items individually"

	default:
		switch err {
		case jira.ErrEmptyResponse:
			msg = "jira: Received empty response.\nPlease try again."
			suggestion = "The server returned an empty response. This may be temporary - try again"
		case jira.ErrNoResult:
			msg = "jira: No results found."
			suggestion = "Try adjusting your search criteria or filters"
		default:
			msg = fmt.Sprintf("Error: %s", err.Error())
		}
	}

	fmt.Fprintf(os.Stderr, "%s\n", msg)
	if suggestion != "" {
		fmt.Fprintf(os.Stderr, "\n💡 %s\n", suggestion)
	}
	os.Exit(1)
}

// Info displays spinner.
func Info(msg string) *spinner.Spinner {
	const refreshRate = 100 * time.Millisecond

	s := spinner.New(
		spinner.CharSets[14],
		refreshRate,
		spinner.WithSuffix(" "+msg),
		spinner.WithHiddenCursor(true),
		spinner.WithWriter(color.Error),
	)
	s.Start()

	return s
}

// Success prints success message in stdout.
func Success(msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stdout, fmt.Sprintf("\n\u001B[0;32m✓\u001B[0m %s\n", msg), args...)
}

// Warn prints warning message in stderr.
func Warn(msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, fmt.Sprintf("\u001B[0;33m%s\u001B[0m\n", msg), args...)
}

// Fail prints failure message in stderr.
func Fail(msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, fmt.Sprintf("\u001B[0;31m✗\u001B[0m %s\n", msg), args...)
}

// Failed prints failure message in stderr and exits.
func Failed(msg string, args ...interface{}) {
	Fail(msg, args...)
	os.Exit(1)
}

// Navigate navigates to jira issue.
func Navigate(server, path string) error {
	url := GenerateServerBrowseURL(server, path)
	return browser.Browse(url)
}

// GenerateServerBrowseURL will return the `browse` URL for a given key.
// The server section can be overridden via `browse_server` in config.
// This is useful if your API endpoint is separate from the web client endpoint.
func GenerateServerBrowseURL(server, key string) string {
	if viper.GetString("browse_server") != "" {
		server = viper.GetString("browse_server")
	}
	return fmt.Sprintf("%s/browse/%s", server, key)
}

// FormatDateTimeHuman formats date time in human readable format.
func FormatDateTimeHuman(dt, format string) string {
	t, err := time.Parse(format, dt)
	if err != nil {
		return dt
	}
	return t.Format("Mon, 02 Jan 06")
}

// GetConfigHome returns the config home directory.
func GetConfigHome() (string, error) {
	home := os.Getenv("XDG_CONFIG_HOME")
	if home != "" {
		return home, nil
	}
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return home + "/.config", nil
}

// StdinHasData checks if standard input has any data to be processed.
func StdinHasData() bool {
	return !term.IsTerminal(int(os.Stdin.Fd()))
}

// ReadFile reads contents of the given file.
func ReadFile(filePath string) ([]byte, error) {
	if filePath != "-" && filePath != "" {
		return os.ReadFile(filePath)
	}
	if filePath == "-" || StdinHasData() {
		b, err := io.ReadAll(os.Stdin)
		_ = os.Stdin.Close()
		return b, err
	}
	return []byte(""), nil
}

// GetJiraIssueKey constructs actual issue key based on given key.
func GetJiraIssueKey(project, key string) string {
	if project == "" {
		return key
	}
	if _, err := strconv.Atoi(key); err != nil {
		return strings.ToUpper(key)
	}
	return fmt.Sprintf("%s-%s", project, key)
}

// NormalizeJiraError normalizes error message we receive from jira.
func NormalizeJiraError(msg string) string {
	msg = strings.TrimSpace(strings.Replace(msg, "Error:\n", "", 1))
	msg = strings.Replace(msg, "- ", "", 1)

	return msg
}

// GetSubtaskHandle fetches actual subtask handle.
// This value can either be handle or name based
// on the used jira version.
func GetSubtaskHandle(issueType string, issueTypes []*jira.IssueType) string {
	get := func(it *jira.IssueType) string {
		if it.Handle != "" {
			return it.Handle
		}
		return it.Name
	}

	var fallback string

	for _, it := range issueTypes {
		if it.Subtask {
			// Exact matches return immediately.
			if strings.EqualFold(issueType, it.Name) {
				return get(it)
			}

			// Store the first subtask type as backup.
			if fallback == "" {
				fallback = get(it)
			}
		}
	}

	// Set default for fallback if none found
	if strings.EqualFold(issueType, jira.IssueTypeSubTask) && fallback == "" {
		fallback = jira.IssueTypeSubTask
	}

	return fallback
}

// GetTUIStyleConfig returns the custom style configured by the user.
func GetTUIStyleConfig() tui.TableStyle {
	var bold bool

	if !viper.IsSet("tui.selection.bold") {
		bold = true
	} else {
		bold = viper.GetBool("tui.selection.bold")
	}

	return tui.TableStyle{
		SelectionBackground: viper.GetString("tui.selection.background"),
		SelectionForeground: viper.GetString("tui.selection.foreground"),
		SelectionTextIsBold: bold,
	}
}
