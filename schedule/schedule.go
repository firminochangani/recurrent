package schedule

import (
	"context"
	"sync"

	"github.com/oklog/ulid/v2"
)

type Schedule struct {
	logger Logger
	lock   *sync.RWMutex
	jobs   map[string]*Job
}

func New() *Schedule {
	return &Schedule{
		lock:   &sync.RWMutex{},
		logger: NewDefaultLogger(),
		jobs:   make(map[string]*Job),
	}
}

func (s *Schedule) Run(ctx context.Context) {
	s.lock.Lock()
	for _, job := range s.jobs {
		go func(ctx context.Context, job *Job) {
			job.run(ctx)
		}(ctx, job)
	}
	s.lock.Unlock()

	<-ctx.Done()

	s.cancelAllJobs()
}

func (s *Schedule) Every(interval int) *Job {
	return &Job{
		scheduler: s,
		interval:  interval,
		id:        ulid.Make().String(),
		done:      make(chan interface{}),
	}
}

func (s *Schedule) appendJob(job *Job) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.jobs[job.id] = job
	s.logger.Infof("Job with id %s has been appended to the scheduler", job.id)
}

func (s *Schedule) RunAll(delayInSeconds int) {}

// GetJobs Return all jobs
func (s *Schedule) GetJobs() []*Job {
	s.lock.RLock()
	defer s.lock.RUnlock()

	index := 0
	result := make([]*Job, len(s.jobs))

	for _, job := range s.jobs {
		result[index] = job
	}

	return result
}

// Clear Stop all jobs and then delete them
func (s *Schedule) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()

	for key, job := range s.jobs {
		job.stop()
		delete(s.jobs, key)
	}
}

func (s *Schedule) CancelJob(job *Job) {
}

func (s *Schedule) cancelAllJobs() {
	s.logger.Info("Cancelling all jobs")

	s.lock.RLock()
	for _, job := range s.jobs {
		job.stop()
	}
	s.lock.RUnlock()

	s.logger.Info("All jobs have been cancelled")
}
