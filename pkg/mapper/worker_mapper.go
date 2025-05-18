package mapper

import (
	"errors"
	"github.com/go-viper/mapstructure/v2"
	"go-takemikazuchi-api/internal/model"
	"go-takemikazuchi-api/internal/worker/dto"
	workerWalletDto "go-takemikazuchi-api/internal/worker_wallet/dto"
	"go-takemikazuchi-api/pkg/exception"
	"go-takemikazuchi-api/pkg/helper"
	"net/http"
)

func MapCreateWorkerDtoIntoWorkerModel(workerModel *model.Worker, createWorkerDto *dto.CreateWorkerDto) {
	err := mapstructure.Decode(createWorkerDto, workerModel)
	helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest, errors.New("bad request")))
}

func MapCreateWorkerWalletDtoIntoWorkerWalletModel(workerWalletModel *model.WorkerWallet, createWorkerWalletDto *workerWalletDto.CreateWorkerWalletDto) {
	err := mapstructure.Decode(createWorkerWalletDto, workerWalletModel)
	helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest, errors.New("bad request")))
}

func MapWorkerModelIntoWorkerResponseDto(workerModel *model.Worker) *dto.WorkerResponseDto {
	var workerResponseDto dto.WorkerResponseDto
	err := mapstructure.Decode(workerModel, &workerResponseDto)
	workerResponseDto.CreatedAt = workerModel.CreatedAt.Format("2006-01-02 15:04:05")
	workerResponseDto.UpdatedAt = workerModel.UpdatedAt.Format("2006-01-02 15:04:05")
	helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest, errors.New("bad request")))
	return &workerResponseDto
}
