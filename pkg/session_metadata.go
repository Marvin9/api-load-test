package pkg

import (
	"fmt"
	"time"
)

type SessionMetadata struct {
	// distribute the milliseconds for each request
	ConstantRate float64
	// break between two requests
	Sleep                time.Duration
	TotalRequests        int64
	FileDescriptorsLimit int64
}

func (s *SessionMetadata) info() {
	warning := ""
	if s.FileDescriptorsLimit > -1 && s.TotalRequests > s.FileDescriptorsLimit {
		mxRequestTime := s.FileDescriptorsLimit * int64(s.ConstantRate)
		warning = fmt.Sprintf("\n\n[WARNING]: With the rate of 1 request in %vms, If round-trip of request is more than %v-%vms, It might exhaust file descriptors with error as 'socket: too many files open'",
			s.ConstantRate, mxRequestTime, mxRequestTime+100)
	}
	fmt.Printf(`Request window: %vms
Total requests: %v%v
`, s.ConstantRate, s.TotalRequests, warning)
}
