package main

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 127.0.0.1 user-identifier frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326

var logFormat = regexp.MustCompile(`(?P<ip>\d+\.\d+\.\d+\.\d+) user-identifier (?P<user>\w+) \[(?P<time>.+)] "(?P<method>\w+) (?P<path>[^\s]+) HTTP/1\.\d" (?P<status>\d+) (?P<bytes>\d+)`)
var timeFormat = "02/Jan/2006:15:04:05 -0700"
var errFormat = errors.New("Invalid Format")

type LogEntry struct {
	IP     net.IP
	User   string
	Time   time.Time
	Method string
	Path   string
	Status int
	Bytes  int
}

func ParseLogEntry(line string) (*LogEntry, error) {
	matches := logFormat.FindStringSubmatch(line)
	reIdx := logFormat.SubexpNames()
	if len(matches) == 0 {
		return nil, errFormat
	}

	match := func(item string) string {
		for n, name := range reIdx {
			if name == item {
				return matches[n]
			}
		}
		return ""
	}

	status, err := strconv.Atoi(match("status"))
	if err != nil {
		return nil, errFormat
	}

	bytes, err := strconv.Atoi(match("bytes"))
	if err != nil {
		return nil, errFormat
	}

	parsedTime, err := time.Parse(timeFormat, match("time"))
	if err != nil {
		fmt.Println("time")
		return nil, errFormat
	}

	ret := &LogEntry{
		IP:     net.ParseIP(match("ip")),
		User:   match("user"),
		Time:   parsedTime,
		Method: strings.ToUpper(match("method")),
		Path:   match("path"),
		Status: status,
		Bytes:  bytes,
	}

	return ret, nil
}
