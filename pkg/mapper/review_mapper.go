package mapper

import (
	"errors"
	"github.com/go-viper/mapstructure/v2"
	"go-takemikazuchi-api/internal/model"
	"go-takemikazuchi-api/internal/review/dto"
	"go-takemikazuchi-api/pkg/exception"
	"go-takemikazuchi-api/pkg/helper"
	"net/http"
)

func MapCreateReviewDtoIntoReviewModel(createReviewDto *dto.CreateReviewDto, reviewModel *model.Review) {
	err := mapstructure.Decode(createReviewDto, reviewModel)
	helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest, errors.New("error when parsing")))
}
