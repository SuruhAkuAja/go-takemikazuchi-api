package transaction

import (
	"go-takemikazuchi-api/internal/model"
	"gorm.io/gorm"
)

type Repository interface {
	FindById(gormTransaction *gorm.DB, id string) *model.Transaction
	Create(gormTransaction *gorm.DB, transactionModel *model.Transaction)
	Update(gormTransaction *gorm.DB, transactionModel *model.Transaction)
	FindWithRelationship(gormTransaction *gorm.DB, id string) *model.Transaction
}
