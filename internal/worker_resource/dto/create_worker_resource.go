package dto

import (
	"mime/multipart"
)

type CreateWorkerWalletDocumentDto struct {
	IdentityCard      *multipart.FileHeader `validate:"obligatoryFile,maxSize=1,extensionFile=jpg pdf txt png jpeg"`
	PoliceCertificate *multipart.FileHeader `validate:"obligatoryFile,maxSize=1,extensionFile=jpg pdf txt png jpeg"`
	DriverLicense     *multipart.FileHeader `validate:"obligatoryFile,maxSize=1,extensionFile=jpg pdf txt png jpeg"`
}
