package jira

import (
	"fmt"
	"strings"
	"time"
)

// SprintStatistics holds sprint metrics.
type SprintStatistics struct {
	SprintID      int
	SprintName    string
	TotalIssues   int
	Completed     int
	InProgress    int
	ToDo          int
	StoryPoints   int
	CompletedSP   int
	CompletionPct float64
	VelocityPct   float64
}

// IssueDistribution holds issue distribution by status.
type IssueDistribution struct {
	Status string
	Count  int
}

// WorklogSummary holds worklog statistics.
type WorklogSummary struct {
	User        string
	TotalHours  float64
	TotalDays   float64
	EntryCount  int
	Issues      []string
	DateRange   string
}

// GetSprintStatistics calculates statistics for a sprint.
func (c *Client) GetSprintStatistics(sprintID int) (*SprintStatistics, error) {
	sprint, err := c.GetSprint(sprintID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sprint: %w", err)
	}

	issues, err := c.SprintIssues(sprintID, "", 0, 1000)
	if err != nil {
		return nil, fmt.Errorf("failed to get sprint issues: %w", err)
	}

	stats := &SprintStatistics{
		SprintID:   sprintID,
		SprintName: sprint.Name,
		TotalIssues: len(issues.Issues),
	}

	completed := 0
	inProgress := 0
	toDo := 0
	totalSP := 0
	completedSP := 0

	for _, issue := range issues.Issues {
		status := strings.ToLower(issue.Fields.Status.Name)
		
		if status == "done" || status == "closed" || status == "resolved" {
			completed++
		} else if status == "in progress" || status == "in review" {
			inProgress++
		} else {
			toDo++
		}

		// Try to get story points (custom field)
		if sp := getStoryPoints(issue); sp > 0 {
			totalSP += sp
			if status == "done" || status == "closed" || status == "resolved" {
				completedSP += sp
			}
		}
	}

	stats.Completed = completed
	stats.InProgress = inProgress
	stats.ToDo = toDo
	stats.StoryPoints = totalSP
	stats.CompletedSP = completedSP

	if stats.TotalIssues > 0 {
		stats.CompletionPct = (float64(completed) / float64(stats.TotalIssues)) * 100
	}
	if totalSP > 0 {
		stats.VelocityPct = (float64(completedSP) / float64(totalSP)) * 100
	}

	return stats, nil
}

// GetIssueDistribution groups issues by status.
func (c *Client) GetIssueDistribution(jql string) ([]IssueDistribution, error) {
	result, err := c.SearchV2(jql, 0, 1000)
	if err != nil {
		return nil, fmt.Errorf("failed to search issues: %w", err)
	}

	statusMap := make(map[string]int)
	for _, issue := range result.Issues {
		status := issue.Fields.Status.Name
		statusMap[status]++
	}

	dist := make([]IssueDistribution, 0, len(statusMap))
	for status, count := range statusMap {
		dist = append(dist, IssueDistribution{
			Status: status,
			Count:  count,
		})
	}

	return dist, nil
}

// GetUserWorklogs retrieves worklogs for a user across issues.
func (c *Client) GetUserWorklogs(user string, from, to time.Time) (*WorklogSummary, error) {
	// Search for issues with worklogs by this user
	jql := fmt.Sprintf("worklogAuthor = %q AND worklogDate >= %s AND worklogDate <= %s",
		user, from.Format("2006-01-02"), to.Format("2006-01-02"))
	
	result, err := c.SearchV2(jql, 0, 1000)
	if err != nil {
		return nil, fmt.Errorf("failed to search issues: %w", err)
	}

	summary := &WorklogSummary{
		User:      user,
		EntryCount: 0,
		Issues:    make([]string, 0),
		DateRange: fmt.Sprintf("%s to %s", from.Format("2006-01-02"), to.Format("2006-01-02")),
	}

	issueSet := make(map[string]bool)
	totalSeconds := 0

	for _, issue := range result.Issues {
		worklogs, err := c.GetWorklogs(issue.Key)
		if err != nil {
			continue
		}

		for _, wl := range worklogs {
			if strings.ToLower(wl.Author.Name) == strings.ToLower(user) ||
				strings.ToLower(wl.Author.DisplayName) == strings.ToLower(user) {
				summary.EntryCount++
				totalSeconds += wl.TimeSpentSeconds
				if !issueSet[issue.Key] {
					summary.Issues = append(summary.Issues, issue.Key)
					issueSet[issue.Key] = true
				}
			}
		}
	}

	summary.TotalHours = float64(totalSeconds) / 3600.0
	summary.TotalDays = summary.TotalHours / 8.0

	return summary, nil
}

// getStoryPoints extracts story points from issue custom fields.
// Note: This requires custom field configuration. Returns 0 if not found.
func getStoryPoints(issue *Issue) int {
	// Story points are stored in custom fields which vary by Jira instance
	// This is a placeholder - actual implementation would need to check
	// configured custom field IDs from the metadata API
	return 0
}

