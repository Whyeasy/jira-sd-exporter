![build](https://github.com/Whyeasy/jira-sd-exporter/workflows/build/badge.svg)
![status-badge](https://goreportcard.com/badge/github.com/Whyeasy/jira-sd-exporter)
![Github go.mod Go version](https://img.shields.io/github/go-mod/go-version/Whyeasy/jira-sd-exporter)

# jira-sd-exporter

A Prometheus Exporter for Jira Service Desk Cloud.

Currently this exporter retrieves the following metrics:

- Service Desk project Info within Jira (Key, Name and ID) `jira_sd_project_info`
- Total amount of tickets(Project key, Done) `jira_sd_ticket`
- Total amount of breached SLAs(Project key, Done) `jira_sd_sla`

## Requirements

Provide your Jira cloud URI; `--jiraURI <string>` or as env variable `JIRA_URI`

Provide a Jira API Key; `--jiraAPIKey` or as env variable `JIRA_API_KEY`

Provide the Jira user who created the API key; `--jiraAPIUser` or as env variable `JIRA_API_USER`

### Optional

Change listening port of the exporter; `--listenAddress <string>` or as env variable `LISTEN_ADDRESS`. Default = `8080`

Change listening path of the exporter; `--listenPath <string>` or as env variable `LISTEN_PATH`. Default = `/metrics`

Change the interval of retrieving data in the background; `--interval <string>` or as env variable `JIRA_INTERVAL`. Default is `60`

To include or exclude projects from the Bugs metrics, please provide a comma separated string with the project keys. Please only provide 1.

Either with `--jiraKeyExclude <string>` or `--jiraKeyInclude <string>`. You can also provide it via env variables `JIRA_KEY_EXCL` or `JIRA_KEY_INCL`.

If you use a custom field to assign issues/tickets to a custom team field. You can provide this custom field; `--jiraTeamCustomField <string>` or as env variable `JIRA_TEAM_CUSTOM_FIELD`. Please provide the full field, `customfield_<number>`

If you want to monitor if your SLAs are breached or not. You can provide these SLAs with a comma separated string to retrieve these; `--jiraSLACustomFields <string,string>` or as env variables `JIRA_SLA_CUSTOM_FIELD`. Please provide the full field, `customfield_<number>,customfield_<number>`
