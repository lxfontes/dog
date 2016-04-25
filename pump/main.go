package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var methods = []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
var statuses = []string{"200", "404", "500", "302", "204"}
var sections = []string{}
var users = []string{"bob", "jebus", "roxie", "luna"}

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func randEntry(options []string) string {
	return options[rand.Intn(len(options))]
}

func randIP() string {
	rndBytes := make([]byte, 4)
	rand.Read(rndBytes)
	return net.IP(rndBytes).String()
}

func printLine() {
	// 127.0.0.1 user-identifier frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326
	rndIP := randIP()
	rndUser := randEntry(users)
	rndMethod := randEntry(methods)
	rndStatus := randEntry(statuses)
	nSections := rand.Intn(3) + 1
	rndSections := []string{""}
	for i := 0; i < nSections; i++ {
		rndSections = append(rndSections, randEntry(sections))
	}
	rndFullPath := strings.Join(rndSections, "/")
	rndBytes := rand.Intn(5000)
	now := time.Now().Format("02/Jan/2006:15:04:05 -0700")

	fmt.Printf("%s user-identifier %s [%s] \"%s %s HTTP/1.1\" %s %d\n", rndIP, rndUser, now, rndMethod, rndFullPath, rndStatus, rndBytes)
}

func main() {
	var waitMs = flag.Int("wait", 100, "Random ms wait between logs")
	var runFor = flag.Duration("run", 30*time.Second, "quit after this duration")

	flag.Parse()

	for i := 0; i < 10; i++ {
		sections = append(sections, randSeq(6))
	}

	quit := make(chan bool)
	go func() {
		time.Sleep(*runFor)
		close(quit)
	}()

	running := true
	for running {
		rndWait := time.Duration(rand.Intn(*waitMs)) * time.Millisecond
		select {
		case <-time.After(rndWait):
			printLine()
		case <-quit:
			running = false
		}
	}
}
