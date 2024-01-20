package schedule

import (
	"context"
	"errors"
	"sync"
	"time"
)

var ErrJobStopped = errors.New("the job has been stopped")

type Job struct {
	interval int
	ticker   *time.Ticker

	lock      *sync.RWMutex
	isRunning bool
	// Due to potential long-running tasks, it's going to be useful to interrupt cancel
	// the context as soon as the job is closed to
	cancelCtx context.CancelCauseFunc
	done      chan interface{}

	id string

	// having a reference to the scheduler will allow us to cal
	// job.Do(func(ctx context.Content) {}) and finally have the job
	// appended to the scheduler
	scheduler *Schedule
	handler   JobHandler
}

type JobHandler func(ctx context.Context)

func (j *Job) Second() *Job {
	return j
}

func (j *Job) Seconds() *Job {
	j.ticker = time.NewTicker(time.Duration(j.interval) * time.Second)
	return j
}

func (j *Job) Do(handler JobHandler) {
	j.handler = handler
	j.scheduler.appendJob(j)
}

func (j *Job) run(ctx context.Context) {
	jobCtx, cancel := context.WithCancelCause(ctx)
	j.cancelCtx = cancel

	j.isRunning = true

	for {
		<-j.ticker.C

		select {
		case <-j.done:
			j.lock.Lock()
			j.isRunning = false
			j.lock.Unlock()
			return
		default:
			j.handler(jobCtx)
		}
	}
}

func (j *Job) stop() {
	j.scheduler.logger.Infof("Stopping the job %s", j.id)

	// Stop further tickers which is rather helpful to prevent the ticker from delaying
	// the job stoppage
	// j.ticker.Stop()

	// send a signal to the job handler to inform it about the stoppage of the job
	if j.cancelCtx != nil {
		j.cancelCtx(ErrJobStopped)
	}

	// done, effectively and immediately
	close(j.done)

	// wait until the job is actually stopped
	for {
		j.lock.RLock()
		if !j.isRunning {
			break
		}
		j.lock.RUnlock()
	}

	j.scheduler.logger.Infof("Job %s has been stopped", j.id)
}

func (j *Job) Minute()    {}
func (j *Job) Minutes()   {}
func (j *Job) Hour()      {}
func (j *Job) Hours()     {}
func (j *Job) Day()       {}
func (j *Job) Days()      {}
func (j *Job) Week()      {}
func (j *Job) Weeks()     {}
func (j *Job) Monday()    {}
func (j *Job) Tuesday()   {}
func (j *Job) Wednesday() {}
func (j *Job) Thursday()  {}
func (j *Job) Friday()    {}
func (j *Job) Saturday()  {}
func (j *Job) Sunday()    {}
func (j *Job) Tag()       {}
func (j *Job) At()        {}
func (j *Job) To()        {}
func (j *Job) Until()     {}
func (j *Job) ShouldRun() {}
