package pkg

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptrace"
	"strings"
	"time"
)

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

func (s *Session) simulateRequest(uniqueId int64) {
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
	start = time.Now()
	res, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		if strings.Contains(err.Error(), "socket: too many open files") {
			fmt.Printf("\n\033[31m[NOTE]: You have to increase the \"ulimit -a\"\n")
		}
		log.Fatal(err)
	}
	res.Body.Close()

	s.Data[uniqueId] = RequestData{
		RequestSequence: uniqueId,
		Performance:     time.Since(start).Milliseconds() - msUntilConnection,
	}
}

func (s *Session) LoadTest(metadata SessionMetadata) {
	s.Data = make([]RequestData, metadata.TotalRequests+1)
	metadata.info()

	fmt.Print("\n\033[34mLoad Testing...\033[0m")
	for i := int64(1); i <= metadata.TotalRequests; i++ {
		s.wg.Add(1)
		go s.simulateRequest(i)
		time.Sleep(metadata.Sleep)
	}
	s.wg.Wait()
}
