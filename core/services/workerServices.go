package services

import (
	"context"
)

type WorkerStore interface {
	Add(worker *Worker) int
	CancelWork(id int)
	Remove(id int) error
}

type Worker interface {
	DoWork(ctx context.Context)
	Cancel()
}
