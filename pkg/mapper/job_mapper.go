package mapper

import (
	"fmt"
	"github.com/go-viper/mapstructure/v2"
	"go-takemikazuchi-api/internal/job/dto"
	jobApplicationDto "go-takemikazuchi-api/internal/job_application/dto"
	"go-takemikazuchi-api/internal/model"
	"go-takemikazuchi-api/pkg/exception"
	"go-takemikazuchi-api/pkg/helper"
	"net/http"
	"strconv"
	"time"
)

func MapJobDtoIntoJobModel[T *dto.CreateJobDto | *dto.UpdateJobDto](jobDto T, jobModel *model.Job) {
	err := mapstructure.Decode(jobDto, &jobModel)
	helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest, err))
}

func MapStringIntoJobResourceModel(jobId uint64, allFilePath []string) []*model.JobResource {
	var jobResourcesModel []*model.JobResource
	for _, filePath := range allFilePath {
		var jobResourceModel model.JobResource
		jobResourceModel.JobId = jobId
		jobResourceModel.ImagePath = filePath
		jobResourcesModel = append(jobResourcesModel, &jobResourceModel)
	}
	return jobResourcesModel
}

func MapJobApplicationModelIntoJobApplicationResponse(jobApplicationsModel []model.JobApplication) []*jobApplicationDto.JobApplicationResponseDto {
	var jobApplicationsResponse []*jobApplicationDto.JobApplicationResponseDto
	for _, jobApplicationModel := range jobApplicationsModel {
		var jobApplicationResponseDto jobApplicationDto.JobApplicationResponseDto
		jobApplicationResponseDto.Id = strconv.FormatUint(jobApplicationModel.ID, 10)

		jobApplicationResponseDto.FullName = jobApplicationModel.User.Name
		jobApplicationResponseDto.AppliedAt = jobApplicationModel.CreatedAt.Format(time.RFC3339)
		jobApplicationsResponse = append(jobApplicationsResponse, &jobApplicationResponseDto)
	}
	return jobApplicationsResponse
}

func MapJobModelIntoJobResponseDto(jobModel []*model.Job) []*dto.JobResponseDto {
	var jobResponseDto []*dto.JobResponseDto
	for _, job := range jobModel {
		var jobResponse dto.JobResponseDto
		jobResponse.ID = job.ID
		jobResponse.Title = job.Title
		jobResponse.Description = job.Description
		jobResponse.CategoryName = job.Category.Name
		jobResponse.Price = job.Price
		jobResponse.Status = job.Status
		jobResponse.CreatedAt = job.CreatedAt.Format(time.RFC3339)
		jobResponse.UpdatedAt = job.UpdatedAt.Format(time.RFC3339)
		jobResponseDto = append(jobResponseDto, &jobResponse)
	}
	return jobResponseDto
}

func MapSingleJobModelIntoSingleJobResponseDto(jobModel *model.Job) *dto.JobResponseDto {
	var jobResponse dto.JobResponseDto
	jobResponse.ID = jobModel.ID
	jobResponse.Title = jobModel.Title
	jobResponse.Description = jobModel.Description
	jobResponse.CategoryName = jobModel.Category.Name
	jobResponse.Price = jobModel.Price
	jobResponse.Status = jobModel.Status
	jobResponse.CreatedAt = jobModel.CreatedAt.Format(time.RFC3339)
	jobResponse.UpdatedAt = jobModel.UpdatedAt.Format(time.RFC3339)
	fmt.Println("Before Mapping User")
	jobResponse.User = MapUserModelIntoUserDto(jobModel.User)
	fmt.Println("Before Mapping User Address")
	jobResponse.UserAddress = MapUserAddressModelIntoUserAddressDto(jobModel.UserAddress)
	if jobModel.Worker != nil {
		fmt.Println("Before Mapping Worker")
		jobResponse.Worker = MapWorkerModelIntoWorkerResponseDto(jobModel.Worker)
	}
	return &jobResponse
}
