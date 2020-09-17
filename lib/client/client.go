package client

import (
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/whyeasy/jira-exporter/lib/jira"
	"github.com/whyeasy/jira-sd-exporter/internal"
)

//Stats struct is the list of expected results to export
type Stats struct {
	Projects *[]ProjectStats
	Tickets  *[]TicketStats
	SLAs     *[]SLAStats
}

//ExporterClient contains Jira information for connecting
type ExporterClient struct {
	jc              *jira.Client
	jiraKeyInclude  string
	jiraKeyExclude  string
	teamCustomField string
	slaCustomFields []string
	interval        time.Duration
}

//New returns a new Client for connecting to Jira
func New(c internal.Config) *ExporterClient {

	convertedTime, _ := strconv.ParseInt(c.Interval, 10, 64)

	exporter := &ExporterClient{
		jc:              jira.NewClient(c.JiraAPIKey, c.JiraAPIUser, c.JiraURI),
		jiraKeyExclude:  c.JiraKeyExclude,
		jiraKeyInclude:  c.JiraKeyInclude,
		teamCustomField: c.TeamCustomField,
		slaCustomFields: strings.Split(c.SLACustomFields, ","),
		interval:        time.Duration(convertedTime),
	}

	exporter.startFetchData()

	return exporter
}

// CachedStats is to store scraped data for caching purposes.
var CachedStats *Stats = &Stats{
	Projects: &[]ProjectStats{},
	SLAs:     &[]SLAStats{},
	Tickets:  &[]TicketStats{},
}

//GetStats retrieves data from API to create metrics from.
func (c *ExporterClient) GetStats() (*Stats, error) {

	return CachedStats, nil
}

func (c *ExporterClient) getData() error {

	projects, err := getProjects(c)
	if err != nil {
		return err
	}

	tickets, err := getTickets(c)
	if err != nil {
		return err
	}

	var slas *[]SLAStats
	if len(c.slaCustomFields) > 0 {
		slas, err = getSLAs(c)
		if err != nil {
			return err
		}
	}

	CachedStats = &Stats{
		Projects: projects,
		Tickets:  tickets,
		SLAs:     slas,
	}

	log.Info("New data retrieved.")

	return nil
}

func (c *ExporterClient) startFetchData() {

	// Do initial call to have data from the start.
	go func() {
		err := c.getData()
		if err != nil {
			log.Error("Scraping failed.")
		}
	}()

	ticker := time.NewTicker(c.interval * time.Second)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				err := c.getData()
				if err != nil {
					log.Error("Scraping failed.")
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
