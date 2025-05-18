package category

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

func (categoryRepository *RepositoryImpl) FindAll(gormTransaction *gorm.DB) []model.Category {
	var categoriesModel []model.Category
	err := gormTransaction.
		Preload("Jobs").
		Find(&categoriesModel).Error
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return categoriesModel
}

func (categoryRepository *RepositoryImpl) IsCategoryExists(categoryId uint64, gormTransaction *gorm.DB) bool {
	var isCategoryExists bool
	gormTransaction.Model(&model.Category{}).
		Select("COUNT(*) > 0").
		Where("id = ?", categoryId).First(&isCategoryExists)
	return isCategoryExists
}
