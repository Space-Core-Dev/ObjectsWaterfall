package services

import (
	"context"
)

type WorkerStore interface {
	Add(worker *Worker) int
	Get(workerId int) (*Worker, error)
	CancelWork(id int) error
	Remove(id int) error
	Exists(name string) bool
}

type Worker interface {
	DoWork(ctx context.Context)
	SetCancel(context.CancelFunc)
	Cancel()
	GetTableName() string
}
