package bbl

import (
	"context"
	"log"
	"sync"
	"time"

	"object-shooter.com/core/models"
	"object-shooter.com/core/services"
	"object-shooter.com/data/repositories"
)

// All logs will be moved to display presenter soon
type SendWorker struct {
	settings   models.BackgroundWorkerSettings
	cancelFunc context.CancelFunc
	group      sync.WaitGroup
	repo       repositories.Repository[string]
}

type dataResult struct {
	data []string
	err  error
}

type requestResult struct {
	requestRes models.ResponseResult
	err        error
}

func NewSendWorker(settings models.BackgroundWorkerSettings, cancel context.CancelFunc) services.Worker {
	repo, err := repositories.NewRepository[string]()
	if err != nil {
		panic(err)
	}
	return &SendWorker{
		settings:   settings,
		cancelFunc: cancel,
		repo:       repo,
	}
}

func (w *SendWorker) DoWork(ctx context.Context) {
	i := 0
	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker has been stoped because of: %s", ctx.Err().Error())
			return
		default:
			w.group.Add(1) // Work starts here that's why w.group.Add(1) is here
			i++
			w.actualWork()
			log.Printf("It has worked %d times", i)
			time.Sleep(time.Duration(w.settings.RequestDelay))
		}
	}
}

func (w *SendWorker) Cancel() {
	w.cancelFunc()
	w.group.Wait()
}

func (w *SendWorker) actualWork() {
	defer w.group.Done()
	dataCh := make(chan dataResult)
	go func() {
		data, err := w.repo.GetData(w.settings.TableName, w.settings.Random, w.settings.WritesNumberToSend, 0)
		dataCh <- dataResult{
			data: data,
			err:  err,
		}
		close(dataCh)
	}()

	dataResult := <-dataCh
	if dataResult.err != nil {
		log.Println(dataResult.err)
		return
	}

	respCh := make(chan requestResult)
	go func() {
		sending := NewSendingService()
		resp, err := sending.SendRequest(w.settings.ConsumerSettings.Host, dataResult.data, nil)
		respCh <- requestResult{
			requestRes: resp,
			err:        err,
		}
		close(respCh)
	}()

	respRes := <-respCh
	if respRes.err != nil {
		log.Println(respRes.err)
		return
	}

	log.Println(respRes)
}
