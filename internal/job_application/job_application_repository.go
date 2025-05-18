package job_application

import (
	"go-takemikazuchi-api/internal/model"
	"gorm.io/gorm"
)

type Repository interface {
	BulkRejectUpdate(gormTransaction *gorm.DB, jobId *uint64)
	Update(gormTransaction *gorm.DB, jobApplicationModel *model.JobApplication)
	FindAllApplication(gormTransaction *gorm.DB, jobId *uint64) []model.JobApplication
	FindById(gormTransaction *gorm.DB, id *uint64, jobId *uint64) *model.JobApplication
}
