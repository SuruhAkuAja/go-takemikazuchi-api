package user

import (
	"fmt"
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

func (userRepository *RepositoryImpl) IsUserExists(gormTransaction *gorm.DB, queryClause string, argumentClause ...interface{}) (bool, error) {
	var isUserExists bool
	fmt.Println(queryClause, argumentClause)
	err := gormTransaction.Model(model.User{}).
		Select("COUNT(*) > 0").
		Where(queryClause, argumentClause...).
		First(&isUserExists).Error
	return isUserExists, err
}

func (userRepository *RepositoryImpl) FindUserByEmail(userEmail *string, userModel *model.User, gormConnection *gorm.DB) {
	err := gormConnection.Where("email = ?", userEmail).First(userModel).Error
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}

func (userRepository *RepositoryImpl) DynamicUpdate(gormTransaction *gorm.DB, whereClause interface{}, updatedValue interface{}, whereArgument ...interface{}) {
	err := gormTransaction.Model(&model.Worker{}).Debug().Where(whereClause, whereArgument).Updates(updatedValue).Error
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}
