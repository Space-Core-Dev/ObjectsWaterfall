package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	bbl "objectswaterfall.com/BBL"
	"objectswaterfall.com/application/dtos"
	"objectswaterfall.com/core/mappers"
	"objectswaterfall.com/core/models"
	"objectswaterfall.com/data/repositories"
	"objectswaterfall.com/stores"
)

func Add(ctx *gin.Context) {
	var workerSettingsDto dtos.BackgroundWorkerSettingsDto
	if err := ctx.ShouldBindBodyWithJSON(&workerSettingsDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	repo, err := repositories.NewRepository[any]()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workerSettings := mappers.FromDtoToWorkerSettings(workerSettingsDto)
	if err := repo.AddSettings(workerSettings); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

func Start(ctx *gin.Context) {
	name := ctx.Query("name")

	store := stores.GetWorkerStore()
	if store.Exists(name) {
		ctx.JSON(http.StatusConflict, gin.H{"Error": fmt.Sprintf("The worker %s is running alredy", name)})
	}

	var consumerSettings models.ConsumerSettings
	if err := ctx.ShouldBindBodyWithJSON(&consumerSettings); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	repo, err := repositories.NewRepository[any]()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workerSettings, err := repo.GetWorkerSettings(name)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	workerSettings.ConsumerSettings = consumerSettings

	duration := time.Now().Add(time.Minute * time.Duration(workerSettings.Timer))
	context, cancel := context.WithDeadline(context.Background(), duration)
	worker := bbl.NewSendWorker(*workerSettings)
	worker.SetCancel(cancel)
	workerId := store.Add(&worker)

	go worker.DoWork(context)

	ctx.JSON(http.StatusOK, gin.H{"workerId": workerId})
}

func Stop(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if id == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.New("id shouldn't be 0")})
		return
	}

	store := stores.GetWorkerStore()
	err = store.CancelWork(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = store.Remove(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": "Ok"})
}

func Seed(ctx *gin.Context) {
	var seedProc bbl.SeedProcessor
	err := ctx.ShouldBindBodyWithJSON(&seedProc)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if repo, err := repositories.NewRepository[any](); err == nil {
		if exists, err := repo.Exists(seedProc.WorkerName); err == nil && !exists {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Errorf("there is no worker named %s", seedProc.WorkerName).Error()})
			return
		}
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	errCh := make(chan error)

	go func() {
		errCh <- seedProc.ProcessJson(false, 0)
	}()

	if err = <-errCh; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": "Ok"})
}

func GetWorkers(ctx *gin.Context) {
	repo, err := repositories.NewRepository[any]()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tables, err := repo.GetAllWorkers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": tables})
}
