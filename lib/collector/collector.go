package collector

import (
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/whyeasy/jira-sd-exporter/lib/client"
)

//Collector struct for holding Prometheus Desc and Exporter Client
type Collector struct {
	up     *prometheus.Desc
	client *client.ExporterClient

	projectInfo *prometheus.Desc
	ticketInfo  *prometheus.Desc
	slaInfo     *prometheus.Desc
}

//New creates a new Collecotor with Prometheus descriptors
func New(c *client.ExporterClient) *Collector {
	log.Info("Creating collector")
	return &Collector{
		up:     prometheus.NewDesc("jira_sd_up", "Whether Jira SD scrap was successful", nil, nil),
		client: c,

		projectInfo: prometheus.NewDesc("jira_sd_project_info", "General information about projects", []string{"project_key", "project_id", "project_name"}, nil),

		ticketInfo: prometheus.NewDesc("jira_sd_ticket", "Information of tickets if they are done or not", []string{"project_key", "done", "team"}, nil),
		slaInfo:    prometheus.NewDesc("jira_sd_sla", "Information of tickets if they are breached or not", []string{"project_key", "breached", "sla_name", "team"}, nil),
	}
}

//Describe the metrics that are collected.
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up

	ch <- c.projectInfo
	ch <- c.ticketInfo
	ch <- c.slaInfo
}

//Collect gathers the metrics that are exported.
func (c *Collector) Collect(ch chan<- prometheus.Metric) {

	log.Info("Running scrape")

	if stats, err := c.client.GetStats(); err != nil {
		log.Error(err)
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 0)
	} else {
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 1)

		collectProjectInfo(c, ch, stats)

		collectTickets(c, ch, stats)

		collectSlas(c, ch, stats)

		log.Info("Scrape Complete")
	}
}

func collectProjectInfo(c *Collector, ch chan<- prometheus.Metric, stats *client.Stats) {
	for _, project := range *stats.Projects {
		ch <- prometheus.MustNewConstMetric(c.projectInfo, prometheus.GaugeValue, 1, project.Key, project.ID, project.Name)
	}
}

func collectTickets(c *Collector, ch chan<- prometheus.Metric, stats *client.Stats) {
	for _, tickets := range *stats.Tickets {
		ch <- prometheus.MustNewConstMetric(c.ticketInfo, prometheus.GaugeValue, tickets.Total, tickets.ProjectKey, strconv.FormatBool(tickets.Done), strings.ToLower(tickets.Team))
	}
}

func collectSlas(c *Collector, ch chan<- prometheus.Metric, stats *client.Stats) {
	for _, slas := range *stats.SLAs {
		ch <- prometheus.MustNewConstMetric(c.slaInfo, prometheus.GaugeValue, slas.Total, slas.ProjectKey, strconv.FormatBool(slas.Breached), slas.Name, strings.ToLower(slas.Team))
	}
}
