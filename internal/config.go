package internal

//Config struct for holding config for exporter and Gitlab
type Config struct {
	ListenAddress   string
	ListenPath      string
	JiraURI         string
	JiraAPIKey      string
	JiraAPIUser     string
	JiraKeyExclude  string
	JiraKeyInclude  string
	TeamCustomField string
	SLACustomFields string
	Interval        string
}
