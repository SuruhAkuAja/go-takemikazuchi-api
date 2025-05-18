package mapper

import (
	"github.com/go-viper/mapstructure/v2"
	"go-takemikazuchi-api/internal/category/dto"
	"go-takemikazuchi-api/internal/model"
	"go-takemikazuchi-api/pkg/exception"
	"go-takemikazuchi-api/pkg/helper"
	"net/http"
)

func MapCategoryDtoIntoCategoryModel[T *dto.CreateCategoryDto | *dto.UpdateCategoryDto](categoryModel *model.Category, categoryDto T) {
	err := mapstructure.Decode(categoryDto, &categoryModel)
	helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest, err))
}

func MapCategoryModelIntoCategoryResponse(categoriesModel []model.Category) []dto.CategoryResponseDto {
	var err error
	var categoriesResponseDto []dto.CategoryResponseDto
	for _, categoryModel := range categoriesModel {
		var categoryResponseDto dto.CategoryResponseDto
		err = mapstructure.Decode(categoryModel, &categoryResponseDto)
		categoriesResponseDto = append(categoriesResponseDto, categoryResponseDto)
	}
	helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest, err))
	return categoriesResponseDto
}

func MapCategoryModelIntoCategoryResponseDto(categoryModel *model.Category) *dto.CategoryResponseDto {
	var err error
	var categoryResponseDto dto.CategoryResponseDto
	err = mapstructure.Decode(categoryModel, &categoryResponseDto)
	helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest, err))
	return &categoryResponseDto
}
