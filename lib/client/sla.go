package client

import (
	"fmt"
	"strings"

	"github.com/whyeasy/jira-exporter/lib/jira"
)

//SLAStats is the struct that holds the data we want from Jira SD
type SLAStats struct {
	ProjectKey string
	Name       string
	Breached   bool
	Total      float64
	Team       string
}

func getSLAs(c *ExporterClient) (*[]SLAStats, error) {

	var jql string

	switch {
	case c.jiraKeyExclude != "":
		jql = fmt.Sprintf("issuetype = 'IT Help' AND project NOT IN (%s)", c.jiraKeyExclude)
	case c.jiraKeyInclude != "":
		jql = fmt.Sprintf("issuetype = 'IT Help' AND project IN (%s)", c.jiraKeyInclude)
	default:
		jql = "issuetype = 'IT Help'"
	}
	var max int
	var expression string
	var apiResults []*jira.ExpressionResult

	for _, customSLA := range c.slaCustomFields {

		if c.teamCustomField != "" {
			max = 500
			expression = fmt.Sprintf("issues.map(issue => ({status: issue.status.name, teamName: issue.%s ? issue.%s.value : 'Unassigned', slaName: issue.%s.name, slaBreached: issue.%s.completedCycles[0] ? issue.%s.completedCycles[0].breached : true , project: issue.project.key })).reduce((result, issue) => result.set(issue.project + ':' + issue.teamName + ':' + issue.slaName,(result[issue.project + ':' + issue.teamName + ':' + issue.slaName] || {}).set(issue.slaBreached == true ? 'BREACHED' : 'NOT_BREACHED',((result[issue.project + ':' + issue.teamName + ':' + issue.slaName] || {})[issue.slaBreached == true ? 'BREACHED' : 'NOT_BREACHED'] || 0) + 1)),new Map())", c.teamCustomField, c.teamCustomField, customSLA, customSLA, customSLA)
		} else {
			max = 600
			expression = fmt.Sprintf("issues.map(issue => ({status: issue.status.name, slaName: issue.%s.name, slaBreached: issue.%s.completedCycles[0] ? issue.%s.completedCycles[0].breached : true , project: issue.project.key })).reduce((result, issue) => result.set(issue.project + ':' + issue.slaName ,(result[issue.project + ':' + issue.slaName] || {}).set(issue.slaBreached == true ? 'BREACHED' : 'NOT_BREACHED',((result[issue.project + ':' + issue.slaName] || {})[issue.slaBreached == true ? 'BREACHED' : 'NOT_BREACHED'] || 0) + 1)),new Map())", customSLA, customSLA, customSLA)
		}

		apiResult, err := c.jc.DoExpression(
			max,
			expression,
			jql)
		if err != nil {
			return nil, err
		}
		apiResults = append(apiResults, apiResult...)
	}

	type breachedCounter struct {
		doneCounter    float64
		notDoneCounter float64
	}
	breachedByProjects := make(map[string]breachedCounter)

	for _, apiResult := range apiResults {

		for project, issues := range apiResult.Value {

			projectIssues, ok := breachedByProjects[project]

			if !ok {
				projectIssues = breachedCounter{}
				breachedByProjects[project] = projectIssues
			}

			projectIssues.doneCounter = projectIssues.doneCounter + issues["BREACHED"]
			projectIssues.notDoneCounter = projectIssues.notDoneCounter + issues["NOT_BREACHED"]

			breachedByProjects[project] = projectIssues
		}
	}

	var results []SLAStats

	for project, counter := range breachedByProjects {
		if c.teamCustomField == "" {

			splits := strings.Split(project, ":")

			results = append(results, SLAStats{
				ProjectKey: splits[0],
				Total:      counter.doneCounter,
				Breached:   true,
				Name:       splits[1],
			})
			results = append(results, SLAStats{
				ProjectKey: splits[0],
				Total:      counter.notDoneCounter,
				Breached:   false,
				Name:       splits[1],
			})
		} else {

			splits := strings.Split(project, ":")

			results = append(results, SLAStats{
				ProjectKey: splits[0],
				Total:      counter.doneCounter,
				Breached:   true,
				Team:       splits[1],
				Name:       splits[2],
			})
			results = append(results, SLAStats{
				ProjectKey: splits[0],
				Total:      counter.notDoneCounter,
				Breached:   false,
				Team:       splits[1],
				Name:       splits[2],
			})
		}
	}

	return &results, nil
}
