package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	bbl "objectswaterfall.com/BBL"
	"objectswaterfall.com/core/models"
	"objectswaterfall.com/data/repositories"
	"objectswaterfall.com/stores"
)

// TODO: In plans making the list of workers which are saved with all settings in database for reuse

func Start(ctx *gin.Context) {
	var workerSettings models.BackgroundWorkerSettings
	if err := ctx.ShouldBindBodyWithJSON(&workerSettings); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	store := stores.GetWorkerStore()
	duration := time.Now().Add(time.Minute * time.Duration(workerSettings.Timer))
	context, cancel := context.WithDeadline(context.Background(), duration)
	worker := bbl.NewSendWorker(workerSettings, cancel)
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
	store.CancelWork(id)
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

func GetTables(ctx *gin.Context) {
	repo, err := repositories.NewRepository[any]()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tables, err := repo.GetAllTables()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": tables})
}
