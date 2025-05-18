package job_application

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

func (jobApplicationRepository *RepositoryImpl) FindAllApplication(gormTransaction *gorm.DB, jobId *uint64) []model.JobApplication {
	var jobApplications []model.JobApplication
	err := gormTransaction.
		Preload("Job").
		Preload("User").
		Joins("JOIN jobs ON jobs.id = job_applications.job_id").
		Joins("JOIN users ON users.id = job_applications.applicant_id").
		Select("job_applications.*, jobs.id as job_id, users.id as applicant_id").
		Where("jobs.id = ?", jobId).
		Find(&jobApplications).Error

	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return jobApplications
}

func (jobApplicationRepository *RepositoryImpl) BulkRejectUpdate(gormTransaction *gorm.DB, jobId *uint64) {
	err := gormTransaction.Where("job_id = ?", jobId).Updates(model.JobApplication{
		Status: "Rejected",
	}).Error
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}

func (jobApplicationRepository *RepositoryImpl) FindById(gormTransaction *gorm.DB, id *uint64, jobId *uint64) *model.JobApplication {
	var jobApplication model.JobApplication
	err := gormTransaction.
		Where("job_id = ? AND id = ?", jobId, id).
		First(&jobApplication).Error
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
	return &jobApplication
}

func (jobApplicationRepository *RepositoryImpl) Update(gormTransaction *gorm.DB, jobApplicationModel *model.JobApplication) {
	err := gormTransaction.
		Where("id = ?", jobApplicationModel.ID).
		Updates(jobApplicationModel).Error
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}
