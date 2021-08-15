package pkg

import (
	"bytes"
	"fmt"
	"log"
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
	Data  Report
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
	if s.Rate <= 0 {
		log.Fatal("request rate should be non-zero")
	}
	var metadata SessionMetadata
	second := time.Second
	metadata.ConstantRate = float64((second / time.Millisecond) / time.Duration(s.Rate))
	metadata.Sleep = time.Millisecond * time.Duration(metadata.ConstantRate)
	metadata.TotalRequests = int64(s.Rate * s.Until)
	// skip the error for now [no side effects], descriptors data is just used for information
	mxDescriptors, _ := s.getMaxDescriptors()
	metadata.FileDescriptorsLimit = mxDescriptors
	return metadata
}

func (s *Session) Success() {
	fmt.Printf("\n\033[32mSuccess.\033[0m\n")
}
