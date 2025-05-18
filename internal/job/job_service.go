package job

import (
	"go-takemikazuchi-api/internal/job/dto"
	userDto "go-takemikazuchi-api/internal/user/dto"
	"go-takemikazuchi-api/pkg/exception"
	"mime/multipart"
)

type Service interface {
	HandleCreate(userJwtClaims *userDto.JwtClaimDto, createJobDto *dto.CreateJobDto, uploadedFiles []*multipart.FileHeader) *exception.ClientError
	HandleUpdate(userJwtClaims *userDto.JwtClaimDto, jobId string, updateJobDto *dto.UpdateJobDto, uploadedFiles []*multipart.FileHeader)
	HandleDelete(userJwtClaims *userDto.JwtClaimDto, jobId string) *exception.ClientError
	HandleRequestCompleted(userJwtClaims *userDto.JwtClaimDto, jobId *string)
	HandleFindAll() []*dto.JobResponseDto
	HandleFindById(jobId uint64) (*dto.JobResponseDto, error)
}
