package api_runner

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRequestPacer(t *testing.T) {
	testCases := []struct {
		name     string
		parallel bool
	}{
		{name: "последовательные запросы"},
		{name: "параллельные запросы", parallel: true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var pacer requestPacer
			started := time.Now()
			sent := make(chan time.Time, 3)

			for range cap(sent) {
				if testCase.parallel {
					go func() {
						pacer.wait()
						sent <- time.Now()
					}()
				} else {
					pacer.wait()
					sent <- time.Now()
				}
			}

			for i := range cap(sent) {
				require.GreaterOrEqual(t, (<-sent).Sub(started), time.Duration(i)*requestInterval)
			}
		})
	}
}
