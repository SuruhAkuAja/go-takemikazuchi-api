package user_address

import (
	"fmt"
	"go-takemikazuchi-api/internal/model"
	"go-takemikazuchi-api/pkg/exception"
	"go-takemikazuchi-api/pkg/helper"
	"gorm.io/gorm"
)

type RepositoryImpl struct {
}

func NewUserAddressRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}

func (userAddressRepository *RepositoryImpl) FindById(gormTransaction *gorm.DB, id *uint64, userAddress *model.UserAddress) {
	err := gormTransaction.Where("id = ?", id).First(userAddress).Error
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}

func (userAddressRepository *RepositoryImpl) Store(gormTransaction *gorm.DB, userAddress *model.UserAddress) {
	err := gormTransaction.Create(userAddress).Error
	fmt.Println("Category", err)
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}
