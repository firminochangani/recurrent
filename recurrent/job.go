package recurrent

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

var ErrJobStopped = errors.New("the job has been stopped")

type JobHandler func(ctx context.Context)

type Job struct {
	ticker *time.Ticker

	// isRunning is read and written in different places simultaneously thus causing data race issues
	// hence the atomic type
	isRunning atomic.Bool

	// Due to potential long-running tasks, it's going to be useful to interrupt cancel
	// the context as soon as the job is closed to
	cancelCtx context.CancelCauseFunc
	done      chan interface{}

	id string

	// having a reference to the manager will allow us to cal
	// job.Do(func(ctx context.Content) {}) and finally have the job
	// appended to the manager
	manager *Manager
	handler JobHandler
}

func (j *Job) Do(handler JobHandler) *Job {
	j.handler = handler
	j.manager.appendJob(j)

	return j
}

func (j *Job) run(ctx context.Context) {
	jobCtx, cancel := context.WithCancelCause(ctx)
	j.cancelCtx = cancel

	j.isRunning.Store(true)

	for {
		<-j.ticker.C

		select {
		case <-j.done:
			j.isRunning.Store(false)
			return
		default:
			j.handler(jobCtx)
		}
	}
}

func (j *Job) runHandler(ctx context.Context) {
	j.isRunning.Store(true)
	j.handler(ctx)
	j.isRunning.Store(false)
}

func (j *Job) stop() {
	j.manager.logger.Infof("Stopping the job %s", j.id)

	// send a signal to the job handler to inform it about the stoppage of the job
	if j.cancelCtx != nil {
		j.cancelCtx(ErrJobStopped)
	}

	// done, effectively and immediately
	close(j.done)

	// wait until the job is actually stopped
	for j.isRunning.Load() {
	}

	j.manager.logger.Infof("Job %s has been stopped", j.id)
}

func (j *Job) Until() {}
