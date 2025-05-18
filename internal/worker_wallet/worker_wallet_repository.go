package worker_wallet

import (
	"go-takemikazuchi-api/internal/model"
	"gorm.io/gorm"
)

type Repository interface {
	FindById(gormTransaction *gorm.DB, walletId *uint64) (*model.WorkerWallet, error)
	Store(gormTransaction *gorm.DB, workerWalletModel *model.WorkerWallet)
	FindByWorkerId(gormTransaction *gorm.DB, workerId *uint64) (*model.WorkerWallet, error)
	DynamicUpdate(gormTransaction *gorm.DB, whereClause interface{}, updatedValue interface{}, whereArgument ...interface{})
}
