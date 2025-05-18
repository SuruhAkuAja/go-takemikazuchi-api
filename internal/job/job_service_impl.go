package job

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go-takemikazuchi-api/configs"
	"go-takemikazuchi-api/internal/category"
	jobDto "go-takemikazuchi-api/internal/job/dto"
	jobResourceFeature "go-takemikazuchi-api/internal/job_resource"
	"go-takemikazuchi-api/internal/model"
	"go-takemikazuchi-api/internal/storage"
	userFeature "go-takemikazuchi-api/internal/user"
	userDto "go-takemikazuchi-api/internal/user/dto"
	userAddressFeature "go-takemikazuchi-api/internal/user_address"
	"go-takemikazuchi-api/internal/user_address/dto"
	validatorFeature "go-takemikazuchi-api/internal/validator"
	"go-takemikazuchi-api/internal/worker"
	"go-takemikazuchi-api/pkg/exception"
	"go-takemikazuchi-api/pkg/helper"
	"go-takemikazuchi-api/pkg/mapper"
	"googlemaps.github.io/maps"
	"gorm.io/gorm"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
)

type ServiceImpl struct {
	validatorService      validatorFeature.Service
	jobRepository         Repository
	userRepository        userFeature.Repository
	categoryRepository    category.Repository
	dbConnection          *gorm.DB
	jobResourceRepository jobResourceFeature.Repository
	fileStorage           storage.FileStorage
	mapsClient            *maps.Client
	userAddressRepository userAddressFeature.Repository
	workerRepository      worker.Repository
	nominatimHttpClient   *configs.HttpClient
}

func NewService(
	jobRepository Repository,
	userRepository userFeature.Repository,
	categoryRepository category.Repository,
	jobResourceRepository jobResourceFeature.Repository,
	dbConnection *gorm.DB,
	fileStorage storage.FileStorage,
	mapsClient *maps.Client,
	userAddressRepository userAddressFeature.Repository,
	workerRepository worker.Repository,
	validatorService validatorFeature.Service,
	nominatimHttpClient *configs.HttpClient,
) *ServiceImpl {
	return &ServiceImpl{
		jobRepository:         jobRepository,
		userRepository:        userRepository,
		categoryRepository:    categoryRepository,
		dbConnection:          dbConnection,
		jobResourceRepository: jobResourceRepository,
		fileStorage:           fileStorage,
		mapsClient:            mapsClient,
		userAddressRepository: userAddressRepository,
		validatorService:      validatorService,
		workerRepository:      workerRepository,
		nominatimHttpClient:   nominatimHttpClient}
}

func (jobService *ServiceImpl) HandleFindAll() []*jobDto.JobResponseDto {

	var jobResponses []*jobDto.JobResponseDto
	err := jobService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		jobModels := jobService.jobRepository.FindAll(gormTransaction)
		jobResponses = mapper.MapJobModelIntoJobResponseDto(jobModels)
		return nil
	})
	helper.CheckErrorOperation(err, exception.NewClientError(http.StatusInternalServerError, exception.ErrInternalServerError, err))
	return jobResponses
}

func (jobService *ServiceImpl) HandleCreate(userJwtClaims *userDto.JwtClaimDto, createJobDto *jobDto.CreateJobDto, uploadedFiles []*multipart.FileHeader) *exception.ClientError {
	err := jobService.validatorService.ValidateStruct(createJobDto)
	jobService.validatorService.ParseValidationError(err)
	err = jobService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		var jobModel model.Job
		var userModel model.User
		var userAddress model.UserAddress
		jobService.userRepository.FindUserByEmail(userJwtClaims.Email, &userModel, gormTransaction)
		if createJobDto.AddressId == nil {
			reverseResponse, err := jobService.nominatimHttpClient.HTTPClient.Get(fmt.Sprintf("%s/reverse?lat=%f&lon=%f&format=json", *jobService.nominatimHttpClient.BaseURL, createJobDto.Latitude, createJobDto.Longitude))
			helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest, errors.New("failed call reverse geocoding")))
			responseBody, readErr := io.ReadAll(reverseResponse.Body)
			helper.CheckErrorOperation(readErr, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest, errors.New("failed call reverse geocoding")))
			var userLocation dto.UserLocation
			jsonErr := json.Unmarshal(responseBody, &userLocation)
			helper.CheckErrorOperation(jsonErr, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest, errors.New("failed call reverse geocoding")))
			userAddress = mapper.MapLocationToUserAddress(userLocation, userModel.ID)
			jobService.userAddressRepository.Store(gormTransaction, &userAddress)
			helper.CheckErrorOperation(err, exception.ParseGormError(err))
		} else {
			jobService.userAddressRepository.FindById(gormTransaction, createJobDto.AddressId, &userAddress)
		}
		isCategoryExists := jobService.categoryRepository.IsCategoryExists(createJobDto.CategoryId, gormTransaction)
		if !isCategoryExists {
			exception.ThrowClientError(exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest, errors.New("category not found")))
		}
		mapper.MapJobDtoIntoJobModel(createJobDto, &jobModel)
		jobModel.UserId = userModel.ID
		jobModel.AddressId = userAddress.ID
		jobModel.WorkerId = nil
		jobService.jobRepository.Store(&jobModel, gormTransaction)
		uuidString := uuid.New().String()
		var allFileName []string

		for _, uploadedFile := range uploadedFiles {
			openedFile, _ := uploadedFile.Open()
			driverLicensePath := fmt.Sprintf("%s-%d-%s", uuidString, jobModel.ID, uploadedFile.Filename)
			_, err = jobService.fileStorage.UploadFile(openedFile, driverLicensePath)
			helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest, errors.New("upload file failed")))
			allFileName = append(allFileName, uploadedFile.Filename)
		}
		if len(allFileName) != 0 {
			resourceModel := mapper.MapStringIntoJobResourceModel(jobModel.ID, allFileName)
			jobService.jobResourceRepository.BulkCreate(gormTransaction, resourceModel)

		}
		return nil
	})
	helper.CheckErrorOperation(err, exception.NewClientError(http.StatusInternalServerError, exception.ErrInternalServerError, err))
	return nil
}

func (jobService *ServiceImpl) HandleUpdate(userJwtClaims *userDto.JwtClaimDto, jobId string, updateJobDto *jobDto.UpdateJobDto, uploadedFiles []*multipart.FileHeader) {
	err := jobService.validatorService.ValidateStruct(updateJobDto)
	jobService.validatorService.ParseValidationError(err)
	err = jobService.validatorService.ValidateVar(jobId, "required|gt=1")
	jobService.validatorService.ParseValidationError(err)
	err = jobService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		var userModel model.User
		parsedJobId, err := strconv.ParseUint(jobId, 10, 64)
		helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest, errors.New("invalid job id")))
		jobService.userRepository.FindUserByEmail(userJwtClaims.Email, &userModel, gormTransaction)
		jobModel, err := jobService.jobRepository.FindVerifyById(gormTransaction, &userModel.Email, &parsedJobId)
		if jobModel.CategoryId != updateJobDto.CategoryId {
			isCategoryExists := jobService.categoryRepository.IsCategoryExists(updateJobDto.CategoryId, gormTransaction)
			if !isCategoryExists {
				exception.ThrowClientError(exception.NewClientError(http.StatusNotFound, exception.ErrNotFound, errors.New("category not found")))
			}
		}
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		mapper.MapJobDtoIntoJobModel(updateJobDto, jobModel)
		jobService.jobRepository.Update(jobModel, gormTransaction)
		resourceModel := jobService.UpdateUploadedFiles(uploadedFiles, jobModel.ID)
		if resourceModel != nil {
			jobService.jobResourceRepository.BulkCreate(gormTransaction, resourceModel)
		}
		if len(updateJobDto.DeletedFilesName) != 0 {
			countFile := jobService.jobResourceRepository.CountBulkByName(gormTransaction, jobModel.ID, updateJobDto.DeletedFilesName)
			jobService.DeleteRequestedFile(updateJobDto.DeletedFilesName, countFile)
			jobService.jobResourceRepository.DeleteBulkByName(gormTransaction, jobModel.ID, updateJobDto.DeletedFilesName)
		}
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}

func (jobService *ServiceImpl) HandleDelete(userJwtClaims *userDto.JwtClaimDto, jobId string) *exception.ClientError {
	err := jobService.validatorService.ValidateVar(jobId, "required|gte=1")
	jobService.validatorService.ParseValidationError(err)
	err = jobService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		var userModel model.User
		parsedJobId, err := strconv.ParseUint(jobId, 10, 64)
		helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest, errors.New("invalid job id")))
		jobService.userRepository.FindUserByEmail(userJwtClaims.Email, &userModel, gormTransaction)
		_, err = jobService.jobRepository.VerifyJobOwner(gormTransaction, userJwtClaims.Email, &parsedJobId)
		if err != nil {
			exception.ThrowClientError(exception.NewClientError(http.StatusUnauthorized, exception.ErrUnauthorized, errors.New("job not belong to user")))
		}
		jobService.jobResourceRepository.DeleteBulkByJobId(gormTransaction, &parsedJobId)
		jobService.jobRepository.Delete(jobId, userModel.ID, gormTransaction)
		return nil
	})
	return nil
}

func (jobService *ServiceImpl) HandleRequestCompleted(userJwtClaims *userDto.JwtClaimDto, jobId *string) {
	err := jobService.validatorService.ValidateVar(jobId, "required")
	jobService.validatorService.ParseValidationError(err)
	err = jobService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		parsedJobId, err := strconv.ParseUint(*jobId, 10, 64)
		var userModel model.User
		helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest, errors.New("invalid job id")))
		jobService.userRepository.FindUserByEmail(userJwtClaims.Email, &userModel, gormTransaction)
		_, err = jobService.jobRepository.VerifyJobOwner(gormTransaction, &userModel.Email, &parsedJobId)
		if err != nil {
			_, err := jobService.jobRepository.VerifyJobWorker(gormTransaction, &userModel.Email, &parsedJobId)
			helper.CheckErrorOperation(err, exception.ParseGormError(err))
		}
		jobModel, err := jobService.jobRepository.FindById(gormTransaction, &parsedJobId)
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		jobModel.Status = "Done"
		jobService.jobRepository.Update(jobModel, gormTransaction)
		jobService.workerRepository.DynamicUpdate(gormTransaction, "id = ?", map[string]interface{}{
			"revenue": jobModel.Price,
		}, jobModel.WorkerId)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}

func (jobService *ServiceImpl) UpdateUploadedFiles(uploadedFiles []*multipart.FileHeader, jobId uint64) []*model.JobResource {
	var resourceModel []*model.JobResource
	uuidString := uuid.New().String()
	if len(uploadedFiles) != 0 {
		var allFileName []string
		for _, uploadedFile := range uploadedFiles {
			openedFile, _ := uploadedFile.Open()
			_, err := jobService.fileStorage.UploadFile(openedFile, fmt.Sprintf("%s-%s", uuidString, uploadedFile.Filename))
			helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest, errors.New("upload file failed")))
			allFileName = append(allFileName, fmt.Sprintf("%s-%s", uuidString, uploadedFile.Filename))
		}
		resourceModel = mapper.MapStringIntoJobResourceModel(jobId, allFileName)
	}
	return resourceModel
}

func (jobService *ServiceImpl) DeleteRequestedFile(deletedFilesName []string, countFile int) {
	if countFile != len(deletedFilesName) {
		exception.ThrowClientError(exception.NewClientError(http.StatusNotFound, "Some files not found", errors.New("count file not equal")))
	}
	for _, deletedFileName := range deletedFilesName {
		_ = jobService.fileStorage.DeleteFile(deletedFileName)
	}
}
