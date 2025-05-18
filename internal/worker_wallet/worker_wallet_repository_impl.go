package worker_wallet

import (
	"go-takemikazuchi-api/internal/model"
	"go-takemikazuchi-api/pkg/exception"
	"go-takemikazuchi-api/pkg/helper"
	"gorm.io/gorm"
)

type RepositoryImpl struct {
}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (workerWalletRepository *RepositoryImpl) FindById(gormTransaction *gorm.DB, walletId *uint64) (*model.WorkerWallet, error) {
	var workerWalletModel model.WorkerWallet
	err := gormTransaction.Where("id = ?", walletId).First(&workerWalletModel).Error
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return &workerWalletModel, err
}

func (workerWalletRepository *RepositoryImpl) Store(gormTransaction *gorm.DB, workerWalletModel *model.WorkerWallet) {
	err := gormTransaction.Create(workerWalletModel).Error
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}

func (workerWalletRepository *RepositoryImpl) FindByWorkerId(gormTransaction *gorm.DB, workerId *uint64) (*model.WorkerWallet, error) {
	var workerWalletModel model.WorkerWallet
	err := gormTransaction.Where("worker_id = ?", workerId).First(&workerWalletModel).Error
	return &workerWalletModel, err
}

func (workerWalletRepository *RepositoryImpl) DynamicUpdate(gormTransaction *gorm.DB, whereClause interface{}, updatedValue interface{}, whereArgument ...interface{}) {
	err := gormTransaction.Model(&model.WorkerWallet{}).Debug().Where(whereClause, whereArgument).Updates(updatedValue).Error
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}
