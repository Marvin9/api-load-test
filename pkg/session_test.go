package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
	// "github.com/Marvin9/api-load-test/pkg"
)

type testSession struct {
	session          Session
	expectedMetadata SessionMetadata
}

func TestGenerateMetadata(t *testing.T) {
	sessions := []testSession{
		{
			session: Session{Rate: 10, Until: 1},
			expectedMetadata: SessionMetadata{
				ConstantRate:  100,
				Sleep:         time.Duration(time.Millisecond * 100),
				TotalRequests: 10,
			},
		},
		{
			session: Session{Rate: 100, Until: 2},
			expectedMetadata: SessionMetadata{
				ConstantRate:  10,
				Sleep:         time.Duration(time.Millisecond * 10),
				TotalRequests: 200,
			},
		},
		{
			session: Session{Rate: 100, Until: 100},
			expectedMetadata: SessionMetadata{
				ConstantRate:  10,
				Sleep:         time.Duration(time.Millisecond * 10),
				TotalRequests: 10000,
			},
		},
		{
			session: Session{Rate: 12345, Until: 12345},
			expectedMetadata: SessionMetadata{
				ConstantRate:  float64(1000 / 12345),
				Sleep:         time.Duration(time.Millisecond * (1000 / 12345)),
				TotalRequests: 12345 * 12345,
			},
		},
		{
			session: Session{Rate: 11, Until: 1},
			expectedMetadata: SessionMetadata{
				ConstantRate:  float64(1000 / 11),
				Sleep:         time.Duration(time.Millisecond * (1000 / 11)),
				TotalRequests: 11,
			},
		},
	}

	for _, testSession := range sessions {
		metadata := testSession.session.GenerateMetadata()

		if metadata.ConstantRate != testSession.expectedMetadata.ConstantRate {
			t.Errorf("request constant rate expected %vms but got %vms",
				testSession.expectedMetadata.ConstantRate, metadata.ConstantRate)
		}

		if metadata.Sleep != testSession.expectedMetadata.Sleep {
			t.Errorf("sleep (break after any request is sent to stabalize the request rate) expected %vms but got %vms",
				testSession.expectedMetadata.Sleep, metadata.Sleep)
		}

		if metadata.TotalRequests != testSession.expectedMetadata.TotalRequests {
			t.Errorf("total requests expected %v but got %v",
				testSession.expectedMetadata.TotalRequests, metadata.TotalRequests)
		}
	}
}

func TestNonZeroRateSession(t *testing.T) {
	runErrorThrowingTest("TestNonZeroRateSession", func() {
		session := Session{
			Rate:  0,
			Until: 10,
		}

		session.GenerateMetadata()
	},
		"request rate should be non-zero",
		t)
}

func runErrorThrowingTest(test string, f func(), errMessageShouldContain string, t *testing.T) {
	if os.Getenv("CRASHER") == "1" {
		f()
	}
	var errOut strings.Builder
	cmd := exec.Command(os.Args[0], fmt.Sprintf("-test.run=%v", test))
	cmd.Env = append(os.Environ(), "CRASHER=1")
	cmd.Stderr = &errOut
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		errStr := errOut.String()
		if !strings.Contains(errStr, errMessageShouldContain) {
			t.Errorf("Thrown error was => %v\nBut it was expected it to contain => %v",
				errStr, errMessageShouldContain)
		}
		return
	}
	t.Errorf("test %v was expected to fail", test)
}
