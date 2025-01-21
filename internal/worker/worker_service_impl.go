package worker

import (
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go-takemikazuchi-api/internal/model"
	"go-takemikazuchi-api/internal/storage"
	"go-takemikazuchi-api/internal/user"
	userDto "go-takemikazuchi-api/internal/user/dto"
	"go-takemikazuchi-api/internal/worker/dto"
	workerResource "go-takemikazuchi-api/internal/worker_resource"
	workerResourceDto "go-takemikazuchi-api/internal/worker_resource/dto"
	workerWallet "go-takemikazuchi-api/internal/worker_wallet"
	"go-takemikazuchi-api/pkg/exception"
	"go-takemikazuchi-api/pkg/helper"
	"go-takemikazuchi-api/pkg/mapper"
	"gorm.io/gorm"
)

type ServiceImpl struct {
	workerRepository         Repository
	userRepository           user.Repository
	validatorInstance        *validator.Validate
	engTranslator            ut.Translator
	dbConnection             *gorm.DB
	workerWalletRepository   workerWallet.Repository
	workerResourceRepository workerResource.Repository
	fileStorage              storage.FileStorage
}

func NewService(
	workerRepository Repository,
	validatorInstance *validator.Validate,
	engTranslator ut.Translator,
	dbConnection *gorm.DB,
	userRepository user.Repository,
	workerWalletRepository workerWallet.Repository,
	workerResourceRepository workerResource.Repository,
	fileStorage storage.FileStorage,
) *ServiceImpl {
	return &ServiceImpl{
		workerRepository:         workerRepository,
		validatorInstance:        validatorInstance,
		engTranslator:            engTranslator,
		dbConnection:             dbConnection,
		workerWalletRepository:   workerWalletRepository,
		workerResourceRepository: workerResourceRepository,
		fileStorage:              fileStorage,
		userRepository:           userRepository,
	}
}

func (workerService *ServiceImpl) Create(userJwtClaim *userDto.JwtClaimDto, createWorkerDto *dto.CreateWorkerDto, createWorkerWalletDocumentDto *workerResourceDto.CreateWorkerWalletDocumentDto) {
	err := workerService.validatorInstance.Struct(createWorkerDto)
	exception.ParseValidationError(err, workerService.engTranslator)
	err = workerService.validatorInstance.Struct(createWorkerWalletDocumentDto)
	exception.ParseValidationError(err, workerService.engTranslator)
	err = workerService.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		var workerModel model.Worker
		var workerWalletModel model.WorkerWallet
		var userModel model.User
		workerService.userRepository.FindUserByEmail(userJwtClaim.Email, &userModel, gormTransaction)
		mapper.MapCreateWorkerDtoIntoWorkerModel(&workerModel, createWorkerDto)
		mapper.MapCreateWorkerWalletDtoIntoWorkerWalletModel(&workerWalletModel, &createWorkerDto.WalletInformation)
		workerModel.UserId = userModel.ID
		workerService.workerRepository.Store(gormTransaction, &workerModel)
		workerWalletModel.WorkerID = workerModel.ID
		workerService.workerWalletRepository.Store(gormTransaction, &workerWalletModel)
		driverLicenseFile, _ := createWorkerWalletDocumentDto.DriverLicense.Open()
		identityCardFile, _ := createWorkerWalletDocumentDto.IdentityCard.Open()
		policeCertificateFile, _ := createWorkerWalletDocumentDto.PoliceCertificate.Open()
		defer driverLicenseFile.Close()
		defer identityCardFile.Close()
		defer policeCertificateFile.Close()
		uuidString := uuid.New().String()
		driverLicensePath := fmt.Sprintf("%s-%s-%s", uuidString, "driverLicense", createWorkerWalletDocumentDto.DriverLicense.Filename)
		policeCertificatePath := fmt.Sprintf("%s-%s-%s", uuidString, "policeCertificate", createWorkerWalletDocumentDto.PoliceCertificate.Filename)
		identityCardPath := fmt.Sprintf("%s-%s-%s", uuidString, "identityCard", createWorkerWalletDocumentDto.IdentityCard.Filename)
		workerService.fileStorage.UploadFile(driverLicenseFile, driverLicensePath)
		workerService.fileStorage.UploadFile(policeCertificateFile, policeCertificatePath)
		workerService.fileStorage.UploadFile(identityCardFile, identityCardPath)
		workerResourcesModel := mapper.MapStringIntoWorkerResourceModel(workerModel.ID,
			[]string{driverLicensePath, policeCertificatePath, identityCardPath},
			[]string{"Driver License", "Police Certificate", "Identity Card"},
			3,
		)
		workerService.workerResourceRepository.BulkStore(gormTransaction, workerResourcesModel)
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}
