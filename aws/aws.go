package aws

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"os"

	"aws_uploader/pkg/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/private/protocol/rest"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
)

type AWS struct {
	logger *logrus.Logger
	config *config.AwsConfig
}

func NewAwsClient(logger *logrus.Logger, config *config.AwsConfig) AwsClient {
	return &AWS{
		logger: logger,
		config: config,
	}
}

// creating a new AWS session
func (a *AWS) CreateSession() *s3.S3 {

	// start aws session
	sess := session.Must(session.NewSession(
		&aws.Config{
			Region: aws.String(a.config.Region),
			Credentials: credentials.NewStaticCredentials(
				a.config.AccessKeyID,
				a.config.SecretAccessKey,
				"",
			),
		},
	))

	// s3 session
	svc := s3.New(sess)

	return svc
}

// Upload File to s3 bucket
func (a *AWS) UploadFile(key string, f multipart.File) (string, error) {

	svc := a.CreateSession()

	ctx := context.Background()

	_, err := svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(a.config.BucketName),
		Key:    aws.String(key),
		Body:   f,
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {
			a.logger.Error(os.Stderr, "upload canceled due to timeout, %v\n", err)
			return "", err
		} else {
			a.logger.Error(os.Stderr, "failed to upload object, %v\n", err)
			return "", err
		}
	}

	a.logger.Infof("successfully uploaded file to %s/%s\n", a.config.BucketName, key)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(a.config.BucketName),
		Key:    aws.String(key),
	})
	rest.Build(req)
	resourceUrl := req.HTTPRequest.URL.String()

	return resourceUrl, nil
}

// Read a file uploaded to s3 bucket
func (a *AWS) ReadFile(key string) *bytes.Buffer {

	svc := a.CreateSession()

	req, resp := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(a.config.BucketName),
		Key:    aws.String(key),
	})
	err := req.Send()
	if err != nil {
		a.logger.Error(os.Stderr, "Unable to send the request %v\n", err)
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		a.logger.Errorf("Failed to read the file %v\n", err)
	}

	return bytes.NewBuffer(body)
}

// Delete a file uploaded to s3 bucket
func (a *AWS) DeleteFile(key string) {

	svc := a.CreateSession()

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(a.config.BucketName), // s3 bucket name
		Key:    aws.String(key),                 // file name
	}

	// Delete the object
	deletedObject, err := svc.DeleteObject(input)
	if err != nil {
		a.logger.Error(err)
		os.Exit(1)
	}

	a.logger.Info("Deleted object --> ", deletedObject)
}
