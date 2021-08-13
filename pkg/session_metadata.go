package pkg

import (
	"fmt"
	"time"
)

type SessionMetadata struct {
	// distribute the milliseconds for each request
	constantRate float64
	// break between two requests
	sleep                time.Duration
	totalRequests        int64
	fileDescriptorsLimit int64
}

func (s *SessionMetadata) info() {
	warning := ""
	if s.fileDescriptorsLimit > -1 && s.totalRequests > s.fileDescriptorsLimit {
		mxRequestTime := s.fileDescriptorsLimit * int64(s.constantRate)
		warning = fmt.Sprintf("\n\n[WARNING]: With the rate of 1 request in %vms, If round-trip of request is more than %v-%vms, It might exhaust file descriptors with error as 'socket: too many files open' error",
			s.constantRate, mxRequestTime, mxRequestTime+100)
	}
	fmt.Printf(`Request window: %vms
Total requests: %v%v
`, s.constantRate, s.totalRequests, warning)
}
