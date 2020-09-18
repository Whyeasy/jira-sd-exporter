package client

import (
	"fmt"
	"strings"
)

//TicketStats is the struct that holds the data we want from Jira SD
type TicketStats struct {
	ProjectKey string
	Total      float64
	Done       bool
	Team       string
}

func getTickets(c *ExporterClient) (*[]TicketStats, error) {

	var jql string

	switch {
	case c.jiraKeyExclude != "":
		jql = fmt.Sprintf("issuetype = 'IT Help' AND project NOT IN (%s)", c.jiraKeyExclude)
	case c.jiraKeyInclude != "":
		jql = fmt.Sprintf("issuetype = 'IT Help' AND project IN (%s)", c.jiraKeyInclude)
	default:
		jql = fmt.Sprintf("issuetype = 'IT Help'")
	}

	var expression string
	var max int

	if c.teamCustomField != "" {
		max = 700
		expression = fmt.Sprintf("issues.map(issue => ({status: issue.status.name, teamName: issue.%s ? issue.%s.value : 'Unassigned', project: issue.project.key })).reduce((result, issue) => result.set(issue.project + ':' + issue.teamName, (result[issue.project + ':' + issue.teamName] || {}).set(issue.status == 'Resolved' || issue.status == 'Canceled' || issue.status == 'Closed' ? 'DONE' : 'NOT_DONE', ((result[issue.project + ':' + issue.teamName] || {})[issue.status == 'Resolved' || issue.status == 'Canceled' || issue.status == 'Closed' ? 'DONE' : 'NOT_DONE'] || 0) + 1)), new Map())", c.teamCustomField, c.teamCustomField)
	} else {
		max = 1000
		expression = "issues.reduce((result, issue) => result.set(issue.project.key, (result[issue.project.key] || {}).set(issue.status.name == 'Resolved' || issue.status.name == 'Canceled' || issue.status.name == 'Closed' ? 'DONE' : 'NOT_DONE', ((result[issue.project.key] || {})[issue.status.name == 'Resolved' || issue.status.name == 'Canceled' || issue.status.name == 'Closed' ? 'DONE' : 'NOT_DONE'] || 0) + 1)), new Map())"
	}

	apiResults, err := c.jc.DoExpression(
		max,
		expression,
		jql)
	if err != nil {
		return nil, err
	}

	type issueTypeCounter struct {
		doneCounter    float64
		notDoneCounter float64
	}
	issueTypesByProjects := make(map[string]issueTypeCounter)

	for _, apiResult := range apiResults {

		for project, issues := range apiResult.Value {

			projectIssues, ok := issueTypesByProjects[project]
			if !ok {
				projectIssues = issueTypeCounter{}
				issueTypesByProjects[project] = projectIssues
			}
			projectIssues.doneCounter = projectIssues.doneCounter + issues["DONE"]
			projectIssues.notDoneCounter = projectIssues.notDoneCounter + issues["NOT_DONE"]

			issueTypesByProjects[project] = projectIssues
		}
	}

	var results []TicketStats

	for project, counter := range issueTypesByProjects {
		if c.teamCustomField == "" {
			results = append(results, TicketStats{
				ProjectKey: project,
				Total:      counter.doneCounter,
				Done:       true,
			})
			results = append(results, TicketStats{
				ProjectKey: project,
				Total:      counter.notDoneCounter,
				Done:       false,
			})
		} else {

			splits := strings.Split(project, ":")

			results = append(results, TicketStats{
				ProjectKey: splits[0],
				Total:      counter.doneCounter,
				Done:       true,
				Team:       splits[1],
			})
			results = append(results, TicketStats{
				ProjectKey: splits[0],
				Total:      counter.notDoneCounter,
				Done:       false,
				Team:       splits[1],
			})
		}
	}

	return &results, nil
}
