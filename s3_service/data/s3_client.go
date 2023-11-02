package data

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func getConfigProfileName() string {
	switch os.Getenv("ENVIRONMENT") {
	case "local":
		return "default"
	case "docker":
		return "docker"
	default:
		return "default"
	}
}

func getS3ClientFromProfile(profile string) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithSharedConfigProfile(profile),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)
	return client, nil
}

func CreateS3Client() (*s3.Client, error) {
	configProfile := getConfigProfileName()
	return getS3ClientFromProfile(configProfile)
}
