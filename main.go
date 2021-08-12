package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptrace"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

func simulateRequestAndGetServerProcessingTime(url string, uniqueId int64) (int64, error) {
	defer wg.Done()
	var start time.Time
	var connection int64

	start = time.Now()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return -1, err
	}
	req.Close = true
	req.Header.Set("Cache-Control", "no-store")
	req.Header.Set("Connection", "close")

	trace := &httptrace.ClientTrace{
		GotConn: func(gci httptrace.GotConnInfo) { connection = time.Since(start).Milliseconds() },
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	req.Header.Set("X-Http-Load-Test", fmt.Sprint(uniqueId))
	res, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		log.Fatalf("[Error]: You have to increase the \"ulimit -a\"\nError: %v", err)
		return -1, err
	}
	res.Body.Close()

	performanceTime := time.Since(start).Milliseconds() - connection
	return performanceTime, nil
}

func getUlimit() (int64, error) {
	cmd := exec.Command("ulimit", "-n")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {
		return -1, err
	}

	ulimit, err := strconv.Atoi(strings.TrimSuffix(out.String(), "\n"))
	if err != nil {
		return -1, err
	}

	return int64(ulimit), nil
}

func main() {
	fmt.Print("\033[H\033[2J")
	url := "http://127.0.0.1:8000"
	rate := 1000 // requests/second
	until := 10  // seconds

	millisecond := 1000
	constantRate := float64(float64(millisecond) / float64(rate))
	constantRateSleep := time.Millisecond * time.Duration(constantRate)
	loopsAtConstantRate := int64(until * millisecond / int(constantRate))

	ulimit, err := getUlimit()
	if err != nil {
		log.Fatal("Error getting ulimit => ", err)
	}

	fmt.Printf(`Number of Requests to process: %v
Each request distributed %vms
ulimit %v

NOTE: Proceed only if you are sure that request time won't exceed %vms, otherwise you have to increase "ulimit -a"

`,
		loopsAtConstantRate, constantRate, ulimit, constantRate*float64(ulimit))

	time.Sleep(time.Second)
	for i := int64(1); i <= loopsAtConstantRate; i++ {
		wg.Add(1)
		go simulateRequestAndGetServerProcessingTime(url, i)
		time.Sleep(constantRateSleep)
	}
	wg.Wait()
}
