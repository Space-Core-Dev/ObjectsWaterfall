package stores

import (
	"errors"

	bbl "object-shooter.com/BBL"
)

type WorkerStore interface {
	Add(worker *bbl.Worker) int
	CancelWork(id int)
	Remove(id int) error
}

type workerStore struct {
	workers map[int]*bbl.Worker
}

var store WorkerStore

func GetWorkerStore() WorkerStore {
	if store != nil {
		return store
	}
	store = &workerStore{}
	return store
}

func (w *workerStore) Add(worker *bbl.Worker) int {
	var last int
	for k := range w.workers {
		last = k
	}
	workerId := last + 1
	w.workers[workerId] = worker

	return workerId
}

func (w *workerStore) CancelWork(workerId int) {
	w.workers[workerId].Cancel()
}

func (w *workerStore) Remove(workerId int) error {
	if _, ok := w.workers[workerId]; !ok {
		return errors.New("wrong worker identifire")
	}
	w.workers[workerId].Cancel()
	delete(w.workers, workerId)
	return nil
}
