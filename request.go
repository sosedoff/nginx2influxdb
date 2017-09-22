package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// Regular expression to match standard nginx access log outout
var logRegexp = regexp.MustCompile(`(\S+) - (\S+) \[([^\]]+)\] "([^"]+)" (\S+) (\S+) "([^"]*?)" "([^"]*?)"( "([^"]*?)")?`)

type Request struct {
	Ip        string    // Remote IP address of the client
	Proto     string    // HTTP protocol
	Method    string    // Request method (GET, POST, etc)
	Host      string    // Requested hostname
	Path      string    // Requested path
	Status    string    // Responses status code (200, 400, etc)
	Referer   string    // Referer (usually is set to "-")
	Agent     string    // User agent string
	Timestamp time.Time // Request timestamp (UTC)
}

// Parse nginx request data
// Example: "GET http://foobar.com/ HTTP/1.1"
func parseRequest(str string, req *Request) error {
	chunks := strings.Split(str, " ")
	if len(chunks) != 3 {
		return fmt.Errorf("invalid request format")
	}

	req.Method = chunks[0]
	req.Proto = chunks[2]

	if uri, err := url.Parse(chunks[1]); err == nil {
		req.Host = uri.Host
		req.Path = uri.Path
	}

	return nil
}

// Parse nginx log timestamp
// Example: 21/Mar/2016:02:33:29 +0000
func parseTimestamp(str string, req *Request) error {
	ts, err := time.Parse("02/Jan/2006:15:04:05 -0700", str)
	if err == nil {
		req.Timestamp = ts
	}
	return err
}

// Produce wire-formatted string for ingestion into influxdb
func (r *Request) InfluxString() string {
	return fmt.Sprintf(
		"requests,ip=%s,method=%s,path=%s,status=%s ip=\"%s\" %v",
		r.Ip,
		r.Method,
		r.Path,
		r.Status,
		r.Ip,
		r.Timestamp.UnixNano(),
	)
}

// Initialize a new request from the input string
func NewRequest(str string) (Request, error) {
	allmatches := logRegexp.FindAllStringSubmatch(str, -1)
	if len(allmatches) == 0 {
		return Request{}, fmt.Errorf("no matches for line: %s", str)
	}
	matches := allmatches[0]

	req := Request{
		Ip:      matches[1],
		Status:  matches[5],
		Referer: matches[7],
		Agent:   matches[8],
	}

	parseTimestamp(matches[3], &req)
	parseRequest(matches[4], &req)

	return req, nil
}
