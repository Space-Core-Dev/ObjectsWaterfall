package bbl

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"objectswaterfall.com/core/models"
	"objectswaterfall.com/core/services"
	"objectswaterfall.com/data/repositories"
	"objectswaterfall.com/utils/stopwatch"
)

// All logs will be moved to display presenter soon
type SendWorker struct {
	settings     models.BackgroundWorkerSettings
	cancelFunc   context.CancelFunc
	group        sync.WaitGroup
	repo         repositories.Repository[string]
	totalSended  int64
	tokenService TokenService
}

type dataResult struct {
	data []string
	err  error
}

type requestResult struct {
	requestRes models.ResponseResult
	err        error
}

func NewSendWorker(settings models.BackgroundWorkerSettings /*, cancel context.CancelFunc */) services.Worker {
	repo, err := repositories.NewRepository[string]()
	if err != nil {
		panic(err)
	}
	return &SendWorker{
		settings: settings,
		//cancelFunc:   cancel,
		repo:         repo,
		tokenService: TokenService{},
	}
}

func (w *SendWorker) DoWork(ctx context.Context) {
	log.Printf("Worker was started at %v", time.Now())
	var counter int64
	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker was stoped at %v, because of: %s ", time.Now(), ctx.Err().Error())
			return
		default:
			w.work(counter)
		}
	}
}

func (w *SendWorker) SetCancel(cancel context.CancelFunc) {
	w.cancelFunc = cancel
}

func (w *SendWorker) Cancel() {
	w.cancelFunc()
	w.group.Wait()
}

func (w *SendWorker) GetTableName() string {
	return w.settings.WorkerName
}

func (w *SendWorker) work(counter int64) {
	w.group.Add(1)
	sw := stopwatch.NewStopWatch()
	sw.Start()
	w.actualWork()
	requstDuration := sw.Elapsed(time.Second)
	log.Printf("Request %d takes %.2f seconds || Total amount of records have been sent: %d", counter, requstDuration, w.totalSended)

	time.Sleep(time.Duration(w.settings.RequestDelay) * time.Second)

	tableCount, _ := w.repo.Count(w.settings.WorkerName)
	stopWhenEnd := w.settings.StopWhenTableEnds
	random := w.settings.Random

	if tableCount <= w.totalSended {
		switch {
		case !stopWhenEnd && !random:
			w.totalSended = 0
		case stopWhenEnd && !random:
			w.cancelFunc()
		}
	}
}

func (w *SendWorker) actualWork() {
	defer w.group.Done()
	dataCh := make(chan dataResult)
	go w.getData(dataCh)

	dataResult := <-dataCh
	if dataResult.err != nil {
		log.Println(dataResult.err)
		return
	}

	respCh := make(chan requestResult)
	go w.sendRequest(dataResult, respCh)

	respRes := <-respCh
	if respRes.err != nil {
		log.Println(respRes.err)
		return
	}

	log.Println(respRes)
}

func (w *SendWorker) getData(dataCh chan dataResult) {
	defer close(dataCh)
	var skip int64
	if w.settings.Random {
		count, err := w.repo.Count(w.settings.WorkerName)
		if err != nil {
			dataCh <- dataResult{
				data: nil,
				err:  err,
			}
			return
		}
		skip = rand.Int63n(count)
	} else {
		skip = w.totalSended
	}

	data, err := w.repo.GetData(w.settings.WorkerName, w.settings.Random, w.settings.WritesNumberToSend, skip)
	dataCh <- dataResult{
		data: data,
		err:  err,
	}
}

func (w *SendWorker) sendRequest(data dataResult, respCh chan requestResult) {
	defer close(respCh)
	sending := NewSendingService()
	var (
		token   string
		headers = make(map[string]string)
	)
	if w.settings.ConsumerSettings.AuthModel != "" {
		var err error
		token, err = w.tokenService.Token()
		if err != nil {
			respCh <- requestResult{
				requestRes: models.ResponseResult{},
				err:        err,
			}
			return
		}
		headers["Authorization"] = fmt.Sprintf("Bearer %s", token)
	}
	resp, err := sending.SendRequest(w.settings.ConsumerSettings.Host, data.data, headers)
	respCh <- requestResult{
		requestRes: resp,
		err:        err,
	}
	w.totalSended += int64(len(data.data))
}
