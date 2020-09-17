FROM alpine

COPY jira-sd-exporter /usr/bin/
ENTRYPOINT ["/usr/bin/jira-sd-exporter"]