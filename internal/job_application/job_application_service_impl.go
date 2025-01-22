package job_application

import (
	"errors"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go-takemikazuchi-api/internal/job"
	"go-takemikazuchi-api/internal/job_application/dto"
	"go-takemikazuchi-api/internal/model"
	userFeature "go-takemikazuchi-api/internal/user"
	userDto "go-takemikazuchi-api/internal/user/dto"
	"go-takemikazuchi-api/pkg/exception"
	"go-takemikazuchi-api/pkg/helper"
	"go-takemikazuchi-api/pkg/mapper"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type ServiceImpl struct {
	validationInstance       *validator.Validate
	engTranslator            ut.Translator
	jobApplicationRepository Repository
	dbConnection             *gorm.DB
	jobRepository            job.Repository
	userRepository           userFeature.Repository
}

func NewService(
	validationInstance *validator.Validate,
	engTranslator ut.Translator,
	jobApplicationRepository Repository,
	dbConnection *gorm.DB,
	jobRepository job.Repository,
	userRepository userFeature.Repository) *ServiceImpl {
	return &ServiceImpl{
		validationInstance:       validationInstance,
		engTranslator:            engTranslator,
		jobApplicationRepository: jobApplicationRepository,
		dbConnection:             dbConnection,
		jobRepository:            jobRepository,
		userRepository:           userRepository,
	}
}

func (jobApplicationService *ServiceImpl) FindAllApplication(userJwtClaims *userDto.JwtClaimDto, jobId string) []*dto.JobApplicationResponseDto {
	err := jobApplicationService.validationInstance.Var(jobId, "required,number,gt=0")
	exception.ParseValidationError(err, jobApplicationService.engTranslator)
	parsedJobId, err := strconv.ParseUint(jobId, 10, 64)
	var jobApplicationsResponse []*dto.JobApplicationResponseDto
	helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest, errors.New("bad request")))
	err = jobApplicationService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		jobApplicationService.jobRepository.VerifyJobOwner(gormTransaction, userJwtClaims.Email, &parsedJobId)
		jobApplications := jobApplicationService.jobApplicationRepository.FindAllApplication(gormTransaction, &parsedJobId)
		jobApplicationsResponse = mapper.MapJobApplicationModelIntoJobApplicationResponse(jobApplications)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return jobApplicationsResponse
}

func (jobApplicationService *ServiceImpl) HandleApply(userJwtClaims *userDto.JwtClaimDto, applyJobApplicationDto *dto.ApplyJobApplicationDto) {
	err := jobApplicationService.validationInstance.Struct(applyJobApplicationDto)
	exception.ParseValidationError(err, jobApplicationService.engTranslator)
	err = jobApplicationService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		var userModel model.User
		var jobApplicationModel model.JobApplication
		jobApplicationService.userRepository.FindUserByEmail(userJwtClaims.Email, &userModel, gormTransaction)
		isJobExists := jobApplicationService.jobRepository.IsExists(applyJobApplicationDto.JobId, gormTransaction)
		if !isJobExists {
			exception.ThrowClientError(exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest, errors.New("job not exists")))
		}
		jobApplicationModel.ApplicantId = userModel.ID
		jobApplicationModel.JobId = applyJobApplicationDto.JobId
		err = jobApplicationService.dbConnection.Create(&jobApplicationModel).Error
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}
