package bbl

import (
	"context"
	"sync"
)

type work func() int

type Worker struct {
	ctx    context.Context
	cancel context.CancelFunc
	work   func() int
	group  sync.WaitGroup
}

func NewWorker(ctx context.Context, CancelFunc context.CancelFunc, work work) *Worker {
	return &Worker{
		ctx:    ctx,
		cancel: CancelFunc,
		work:   work,
		group:  sync.WaitGroup{},
	}
}

func (w *Worker) Cancel() {
	w.group.Wait()
	w.cancel()
}

func (w *Worker) StartWork() {
	var err error
	for {
		select {
		case <-w.ctx.Done():
			err = w.ctx.Err()
			w.Cancel()
		default:
			w.group.Add(1)
			w.work() // socket connection will return logs and work return value
			w.group.Done()
		}
		if err != nil {
			break
		}
	}
	// socket connection will return err value
}
