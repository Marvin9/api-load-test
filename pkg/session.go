package pkg

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

type Session struct {
	TargetEndpoint string
	Method         string
	// Rate = x requests/second
	Rate int
	// Until = x, stop after x seconds
	Until int
	wg    sync.WaitGroup
}

func (s *Session) getMaxDescriptors() (int64, error) {
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

// generates helpful numbers to control loop, time and prediction(resource exhaust), for given rate & limit in session
func (s *Session) GenerateMetadata() SessionMetadata {
	var metadata SessionMetadata
	second := time.Second
	metadata.constantRate = float64((second / time.Millisecond) / time.Duration(s.Rate))
	metadata.sleep = time.Millisecond * time.Duration(metadata.constantRate)
	metadata.totalRequests = int64(s.Rate * s.Until)
	// skip the error for now [no side effects], descriptors data is just used for information
	mxDescriptors, _ := s.getMaxDescriptors()
	metadata.fileDescriptorsLimit = mxDescriptors
	return metadata
}

func (s *Session) prepareBasicRequest() (*http.Request, error) {
	req, err := http.NewRequest(s.Method, s.TargetEndpoint, nil)
	if err != nil {
		return req, err
	}
	req.Close = true
	req.Header.Set("Cache-Control", "no-store")
	req.Header.Set("Connection", "close")
	return req, nil
}

func (s *Session) simulateRequest(uniqueId int64) int64 {
	defer s.wg.Done()
	var start time.Time
	var msUntilConnection int64

	req, err := s.prepareBasicRequest()
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("X-Http-Load-Test", fmt.Sprint(uniqueId))

	trace := &httptrace.ClientTrace{
		GotConn: func(gci httptrace.GotConnInfo) { msUntilConnection = time.Since(start).Milliseconds() },
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	res, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		if strings.Contains(err.Error(), "socket: too many open files") {
			fmt.Printf("\n\033[31m[NOTE]: You have to increase the \"ulimit -a\"\n")
		}
		log.Fatal(err)
	}
	res.Body.Close()

	return time.Since(start).Milliseconds() - msUntilConnection
}

func (s *Session) LoadTest(metadata SessionMetadata) {
	metadata.info()

	for i := int64(1); i <= metadata.totalRequests; i++ {
		s.wg.Add(1)
		go s.simulateRequest(i)
		time.Sleep(metadata.sleep)
	}
	s.wg.Wait()
}

func (s *Session) Success() {
	fmt.Printf("\n\033[32mSuccess.\n")
}
