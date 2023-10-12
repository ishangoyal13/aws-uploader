package aws

import (
	"bytes"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/service/s3"
)

type AwsClient interface {
	CreateSession() *s3.S3
	UploadFile(key string, f multipart.File) (string, error)
	ReadFile(key string) *bytes.Buffer
	DeleteFile(key string)
}
