package stores

import (
	"errors"

	"objectswaterfall.com/core/services"
)

type workerStore struct {
	workers map[int]*services.Worker
}

var store services.WorkerStore

func GetWorkerStore() services.WorkerStore {
	if store != nil {
		return store
	}
	store = &workerStore{
		workers: map[int]*services.Worker{},
	}
	return store
}

func (w *workerStore) Add(worker *services.Worker) int {
	var last int
	for k := range w.workers {
		last = k
	}
	workerId := last + 1
	w.workers[workerId] = worker

	return workerId
}

func (w *workerStore) CancelWork(workerId int) {
	(*w.workers[workerId]).Cancel()
}

func (w *workerStore) Remove(workerId int) error {
	if _, ok := w.workers[workerId]; !ok {
		return errors.New("wrong worker identifire")
	}
	(*w.workers[workerId]).Cancel()
	delete(w.workers, workerId)
	return nil
}
