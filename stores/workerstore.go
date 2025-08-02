package stores

import (
	"errors"
	"fmt"

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

func (w *workerStore) Get(workerId int) (*services.Worker, error) {
	if worker, ok := w.workers[workerId]; ok {
		return worker, nil
	}

	return nil, errors.New("wrong worker identifire")
}

func (w *workerStore) Exists(name string) bool {
	for _, v := range w.workers {
		if (*v).GetTableName() == name {
			return true
		}
	}

	return false
}

func (w *workerStore) CancelWork(workerId int) error {
	if _, ok := (*w).workers[workerId]; !ok {
		return fmt.Errorf("there is no worker with id %d", workerId)
	}
	(*w.workers[workerId]).Cancel()
	return nil
}

func (w *workerStore) Remove(workerId int) error {
	if _, ok := w.workers[workerId]; !ok {
		return errors.New("wrong worker identifire")
	}
	(*w.workers[workerId]).Cancel()
	delete(w.workers, workerId)
	return nil
}
