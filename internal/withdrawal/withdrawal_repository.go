package withdrawal

import (
	"go-takemikazuchi-api/internal/model"
	"gorm.io/gorm"
)

type Repository interface {
	Create(gormTransaction *gorm.DB, withdrawalModel *model.Withdrawal)
	FindAll(gormTransaction *gorm.DB) []model.Withdrawal
	FindById(gormTransaction *gorm.DB, withdrawalId *uint64) (*model.Withdrawal, error)
	Update(gormTransaction *gorm.DB, withdrawalModel *model.Withdrawal)
}
