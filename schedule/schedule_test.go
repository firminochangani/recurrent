package schedule_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/flowck/schedule/schedule"
)

func TestSchedule(t *testing.T) {
	testCases := []struct {
		name                           string
		expectedNumberOfCallsToHandler int
		timeout                        time.Duration
		job                            func(s *schedule.Schedule) *schedule.Job
	}{
		{
			name:                           "schedule_every_1_seconds_for_5_seconds",
			expectedNumberOfCallsToHandler: 5,
			timeout:                        time.Second * 5,
			job: func(s *schedule.Schedule) *schedule.Job {
				return s.Every(1).Seconds()
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tc.timeout)
			defer cancel()

			s := schedule.New()
			counter := 0

			tc.job(s).Do(func(ctx context.Context) {
				counter++
				t.Logf("counter incremented to --> %d\n", counter)
			})

			s.Run(ctx)

			assert.Equal(t, tc.expectedNumberOfCallsToHandler, counter, "expect counter to have been incremented by 5")
		})
	}
}

func TestSchedule_Clear(t *testing.T) {
	s := schedule.New()

	// Given x amount of handlers
	for i := 0; i < 5; i++ {
		s.Every(1).Seconds().Do(func(ctx context.Context) {
			t.Log("Handler ran successfully")
		})
	}
	require.NotEmpty(t, s.GetJobs())

	// Do the following
	s.Clear()

	// Then expect
	assert.Empty(t, s.GetJobs())
}
