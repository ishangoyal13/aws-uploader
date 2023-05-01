package main

import (
	"aws_uploader/aws"
	"aws_uploader/pkg/config"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {

	logrus.SetReportCaller(true)

	err := godotenv.Load(".env")
	if err != nil {
		logrus.Fatalf("Error loading .env file")
		return
	}

	// initializing aws client
	awsClient := aws.NewAwsClient(logrus.New(), &config.AwsConfig{
		Region:          "ap-south-1",
		BucketName:      os.Getenv("BUCKET_NAME"),
		AccessKeyID:     os.Getenv("ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("SECRET_ACCESS_KEY"),
	})

	logrus.Info(os.Getenv("ACCESS_KEY_ID"))

	file, err := os.Open("assests/sample.txt")
	if err != nil {
		logrus.Error(err)
		return
	}

	uploadUrl, err := awsClient.UploadFile("your_unique_file_key", file)
	if err != nil {
		logrus.Error(err)
		return
	}

	logrus.Info("Upload url --> ", uploadUrl)
}
