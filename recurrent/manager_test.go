package recurrent_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/flowck/recurrent/recurrent"
)

func TestManager(t *testing.T) {
	testCases := []struct {
		name                           string
		expectedNumberOfCallsToHandler int
		timeout                        time.Duration
		job                            func(s *recurrent.Manager) *recurrent.Job
	}{
		{
			name:                           "every_1_seconds_for_5_seconds",
			expectedNumberOfCallsToHandler: 5,
			timeout:                        time.Second * 5,
			job: func(s *recurrent.Manager) *recurrent.Job {
				return s.Every(time.Second * 1)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tc.timeout)
			defer cancel()

			s := recurrent.New()
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

func TestManager_Clear(t *testing.T) {
	s := recurrent.New()

	// Given x amount of handlers
	for i := 0; i < 5; i++ {
		s.Every(time.Second * 1).Do(func(ctx context.Context) {
			t.Log("Handler ran successfully")
		})
	}
	require.NotEmpty(t, s.GetJobs())

	// Do the following
	s.Clear()

	// Then expect
	assert.Empty(t, s.GetJobs())
}

func TestManager_CancelJob(t *testing.T) {
	s := recurrent.New()

	// Given
	job := s.Every(time.Second * 1).Do(func(ctx context.Context) {
		t.Log("Handler ran successfully")
	})
	require.Contains(t, s.GetJobs(), job)

	// Do
	s.CancelJob(job)

	// Expect
	assert.NotContains(t, s.GetJobs(), job)
}

func TestManager_RunAll(t *testing.T) {
	s := recurrent.New()
	result := make([]int, 3)

	handler := func(idx int) func(ctx context.Context) {
		return func(ctx context.Context) {
			result[idx] = idx
			t.Logf("Job %d ran successfully", idx)
		}
	}

	// Given
	s.Every(time.Second * 1).Do(handler(0))
	s.Every(time.Second * 1).Do(handler(1))
	s.Every(time.Second * 1).Do(handler(2))

	// Do
	s.RunAll(context.Background(), time.Millisecond*10)

	// Expect
	assert.Equal(t, []int{0, 1, 2}, result, "all jobs to have been run and mutated the 'result' slice")
}
