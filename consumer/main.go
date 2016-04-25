package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func clearScreen() {
	fmt.Print("\033[2J\033[1;1H") // clear + move up ;)
}

type Monitor struct {
	tc     *SlidingCounter
	hc     *HitCounter
	ic     *HitCounter
	alerts *AlertMonitor
}

func NewMonitor(reqAlert int, reqInterval time.Duration) *Monitor {
	ret := &Monitor{
		tc: NewSlidingCounter(1*time.Second, 300), // keep raw 5 minutes
	}

	ret.alerts = NewAlertMonitor(func() int {
		return ret.tc.CountEntriesBack("requests", reqInterval)
	}, 1*time.Second, reqAlert)

	ret.clearHits()

	return ret
}

func (m *Monitor) Feed(log *LogEntry) {
	m.tc.Push("requests", 1)
	m.tc.Push("bytes", log.Bytes)

	if log.Status >= 200 && log.Status <= 300 {
		m.tc.Push("ok_status", 1)
	} else {
		m.tc.Push("failed_status", 1)
	}

	switch log.Method {
	case "GET":
		m.tc.Push("get_reqs", 1)
	case "POST":
		m.tc.Push("post_reqs", 1)
	case "PUT":
		m.tc.Push("put_reqs", 1)
	case "PATCH":
		m.tc.Push("patch_reqs", 1)
	case "DELETE":
		m.tc.Push("delete_reqs", 1)
	}

	parts := strings.Split(log.Path, "/")
	m.hc.Push(parts[1])
	m.ic.Push(log.IP.String())
}

func (m *Monitor) Watch() {
	for {
		time.Sleep(10 * time.Second)
		clearScreen()
		m.PrintStats()
		m.clearHits()
	}
}

func (m *Monitor) clearHits() {
	m.hc = NewHitCounter()
	m.ic = NewHitCounter()
}

func (m *Monitor) PrintStats() {
	fmt.Printf("Stats from past minute [%s]\n", time.Now())

	monitorDuration := -1 * time.Minute

	sumBytes := m.tc.SumEntriesBack("bytes", monitorDuration)
	countRequests := m.tc.CountEntriesBack("requests", monitorDuration)

	countGET := m.tc.CountEntriesBack("get_reqs", monitorDuration)
	countPOST := m.tc.CountEntriesBack("post_reqs", monitorDuration)
	countPUT := m.tc.CountEntriesBack("put_reqs", monitorDuration)
	countDELETE := m.tc.CountEntriesBack("delete_reqs", monitorDuration)
	countPATCH := m.tc.CountEntriesBack("patch_reqs", monitorDuration)

	fmt.Printf("Requests(count): \t%d\n", countRequests)
	fmt.Printf("\tGET=%d POST=%d PUT=%d PATCH=%d DELETE=%d\n",
		countGET,
		countPOST,
		countPUT,
		countPATCH,
		countDELETE)
	fmt.Printf("Bytes(sum): \t\t%d\n", sumBytes)

	fmt.Println()
	fmt.Printf("Top Sections from past 10 seconds\n")
	descendingHits := m.hc.TopN(10)
	for _, section := range descendingHits {
		fmt.Printf("\tHits:%d\t\tPath: %s\n", section.Counter, section.Key)
	}

	fmt.Println()
	fmt.Printf("Top IPs from past 10 seconds\n")
	descendingIPs := m.ic.TopN(10)
	for _, section := range descendingIPs {
		fmt.Printf("\tHits:%d\t\tIP: %s\n", section.Counter, section.Key)
	}

	fmt.Println()
	m.alerts.Print()
}

func main() {
	var alertThreshold = flag.Int("threshold", 1000, "Alert threshold")
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)
	monitor := NewMonitor(*alertThreshold, -2*time.Minute)

	clearScreen()
	fmt.Println("stats incoming....")

	go monitor.Watch()

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		log, err := ParseLogEntry(line)
		if err != nil {
			fmt.Println(err)
			fmt.Println(line)
			os.Exit(1)
		}
		monitor.Feed(log)
	}
}
