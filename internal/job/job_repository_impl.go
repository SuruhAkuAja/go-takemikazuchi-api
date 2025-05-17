package job

import (
	"go-takemikazuchi-api/internal/model"
	"go-takemikazuchi-api/pkg/exception"
	"go-takemikazuchi-api/pkg/helper"
	"gorm.io/gorm"
)

type RepositoryImpl struct{}

func NewRepository() *RepositoryImpl {
	return &RepositoryImpl{}
}
func (jobRepository *RepositoryImpl) FindAll(gormTransaction *gorm.DB) []*model.Job {
	var jobModels []*model.Job
	err := gormTransaction.Model(model.Job{}).
		Preload("User").
		Preload("Category").
		Joins("JOIN users ON users.id = jobs.user_id").
		Joins("JOIN categories ON categories.id = jobs.category_id").
		Find(&jobModels).Error
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return jobModels
}
func (jobRepository *RepositoryImpl) FindById(gormTransaction *gorm.DB, id *uint64) (*model.Job, error) {
	var jobModel model.Job
	err := gormTransaction.Model(model.Job{}).Where("id = ?", *id).First(&jobModel).Error
	if err != nil {
		return nil, err
	}
	return &jobModel, nil
}

func (jobRepository *RepositoryImpl) Store(jobModel *model.Job, gormTransaction *gorm.DB) {
	err := gormTransaction.Create(jobModel).Error
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}

func (jobRepository *RepositoryImpl) Update(jobModel *model.Job, gormTransaction *gorm.DB) {
	err := gormTransaction.Debug().Where("id = ?", jobModel.ID).Updates(&jobModel).Error
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}

func (jobRepository *RepositoryImpl) Delete(jobId string, userId uint64, gormTransaction *gorm.DB) {
	err := gormTransaction.Joins("Users").Where("id = ? AND users.id = ?", jobId, userId).Delete(&model.Job{}).Error
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}

func (jobRepository *RepositoryImpl) IsExists(jobId uint64, gormTransaction *gorm.DB) bool {
	var isJobExists bool
	gormTransaction.Model(&model.Job{}).
		Select("COUNT(*) > 0").
		Where("id = ?", jobId).First(&isJobExists)
	return isJobExists
}

func (jobRepository *RepositoryImpl) VerifyJobOwner(gormTransaction *gorm.DB, userEmail *string, jobId *uint64) (bool, error) {
	var isJobValid bool
	err := gormTransaction.Model(&model.Job{}).
		Joins("JOIN users ON users.id = jobs.user_id").
		Select("COUNT(*) > 0").
		Where("jobs.id = ? AND users.email = ?", jobId, userEmail).
		First(&isJobValid).Error
	return isJobValid, err

}

func (jobRepository *RepositoryImpl) VerifyJobWorker(gormTransaction *gorm.DB, workerId *string, jobId *uint64) (bool, error) {
	var isJobValid bool
	err := gormTransaction.Model(&model.Job{}).
		Select("COUNT(*) > 0").
		Where("jobs.id = ? AND workers.id = ?", jobId, workerId).
		First(&isJobValid).Error
	return isJobValid, err
}

func (jobRepository *RepositoryImpl) FindVerifyById(gormTransaction *gorm.DB, userEmail *string, jobId *uint64) (*model.Job, error) {
	var jobModel model.Job
	err := gormTransaction.Model(&model.Job{}).
		Joins("JOIN users ON users.id = jobs.user_id").
		Select("jobs.*").
		Where("jobs.id = ? AND users.email = ?", jobId, userEmail).
		First(&jobModel).Error
	return &jobModel, err
}

func (jobRepository *RepositoryImpl) FindWithRelationship(gormTransaction *gorm.DB, userEmail *string, jobId *uint64) *model.Job {
	var jobModel model.Job
	err := gormTransaction.
		Preload("User").
		Preload("Category").
		Joins("JOIN users ON users.id = jobs.user_id").
		Joins("JOIN categories ON categories.id = jobs.category_id").
		Select(`
			jobs.*, 
			jobs.id AS job_id, 
			users.id AS user_id, users.name AS user_name, users.email AS user_email,
			categories.id AS category_id, categories.name AS category_name
		`).Where("jobs.id = ? AND users.email = ?", jobId, userEmail).
		First(&jobModel).Error
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return &jobModel
}
