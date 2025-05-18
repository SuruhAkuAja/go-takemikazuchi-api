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
	validatorFeature "go-takemikazuchi-api/internal/validator"
	workerFeature "go-takemikazuchi-api/internal/worker"
	"go-takemikazuchi-api/pkg/exception"
	"go-takemikazuchi-api/pkg/helper"
	"go-takemikazuchi-api/pkg/mapper"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type ServiceImpl struct {
	validatorService         validatorFeature.Service
	jobApplicationRepository Repository
	dbConnection             *gorm.DB
	jobRepository            job.Repository
	userRepository           userFeature.Repository
	workerRepository         workerFeature.Repository
}

func NewService(
	validationInstance *validator.Validate,
	engTranslator ut.Translator,
	jobApplicationRepository Repository,
	dbConnection *gorm.DB,
	jobRepository job.Repository,
	userRepository userFeature.Repository,
	validatorService validatorFeature.Service,
	workerRepository workerFeature.Repository,
) *ServiceImpl {
	return &ServiceImpl{
		validatorService:         validatorService,
		jobApplicationRepository: jobApplicationRepository,
		dbConnection:             dbConnection,
		jobRepository:            jobRepository,
		userRepository:           userRepository,
		workerRepository:         workerRepository,
	}
}

func (jobApplicationService *ServiceImpl) FindAllApplication(userJwtClaims *userDto.JwtClaimDto, jobId string) []*dto.JobApplicationResponseDto {
	err := jobApplicationService.validatorService.ValidateVar(jobId, "required,number,gt=0")
	jobApplicationService.validatorService.ParseValidationError(err)
	parsedJobId, err := strconv.ParseUint(jobId, 10, 64)
	var jobApplicationsResponse []*dto.JobApplicationResponseDto
	helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest, errors.New("bad request")))
	err = jobApplicationService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		_, err := jobApplicationService.jobRepository.VerifyJobOwner(gormTransaction, userJwtClaims.Email, &parsedJobId)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		jobApplications := jobApplicationService.jobApplicationRepository.FindAllApplication(gormTransaction, &parsedJobId)
		jobApplicationsResponse = mapper.MapJobApplicationModelIntoJobApplicationResponse(jobApplications)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return jobApplicationsResponse
}

func (jobApplicationService *ServiceImpl) HandleApply(userJwtClaims *userDto.JwtClaimDto, applyJobApplicationDto *dto.ApplyJobApplicationDto) {
	err := jobApplicationService.validatorService.ValidateStruct(applyJobApplicationDto)
	jobApplicationService.validatorService.ParseValidationError(err)
	err = jobApplicationService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		var userModel model.User
		var jobApplicationModel model.JobApplication
		jobApplicationService.userRepository.FindUserByEmail(userJwtClaims.Email, &userModel, gormTransaction)
		isJobExists := jobApplicationService.jobRepository.IsExists(applyJobApplicationDto.JobId, gormTransaction)
		if !isJobExists {
			exception.ThrowClientError(exception.NewClientError(http.StatusNotFound, exception.ErrNotFound, errors.New("job not exists")))
		}
		jobApplicationModel.ApplicantId = userModel.ID
		jobApplicationModel.JobId = applyJobApplicationDto.JobId
		err = jobApplicationService.dbConnection.Create(&jobApplicationModel).Error
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}

func (jobApplicationService *ServiceImpl) SelectApplication(userJwtClaims *userDto.JwtClaimDto, selectApplicationDto *dto.SelectApplicationDto) {
	err := jobApplicationService.validatorService.ValidateStruct(selectApplicationDto)
	jobApplicationService.validatorService.ParseValidationError(err)
	err = jobApplicationService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		jobApplicationModel := jobApplicationService.jobApplicationRepository.FindById(gormTransaction, &selectApplicationDto.Id, &selectApplicationDto.JobId)
		id, err := jobApplicationService.jobRepository.FindVerifyById(gormTransaction, userJwtClaims.Email, &selectApplicationDto.JobId)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		workerModel, err := jobApplicationService.workerRepository.FindById(gormTransaction, &jobApplicationModel.ApplicantId)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		jobModel := id
		jobModel.Status = "Process"
		jobModel.WorkerId = &workerModel.ID
		jobApplicationModel.Status = "Accepted"
		jobApplicationService.jobApplicationRepository.BulkRejectUpdate(gormTransaction, &jobModel.ID)
		jobApplicationService.jobApplicationRepository.Update(gormTransaction, jobApplicationModel)
		jobApplicationService.jobRepository.Update(jobModel, gormTransaction)
		return nil
	})
}
