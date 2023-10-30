package services

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func UploadAuthorImageToS3Bucket(headshot []byte) error {
	var configProfile string

	switch os.Getenv("ENVIRONMENT") {
	case "local":
		configProfile = "default"
	case "docker":
		configProfile = "docker"
	}

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithSharedConfigProfile(configProfile),
	)

	if err != nil {
		fmt.Printf("Error: %v", err.Error())
		return err
	}

	creds, err := cfg.Credentials.Retrieve(context.TODO())
	if err != nil {
		fmt.Printf("Error: %v", err.Error())
		return err
	}
	fmt.Printf("Access key ID: %v\n", creds.AccessKeyID)
	client := s3.NewFromConfig(cfg)

	output, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	// output, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
	// 	Bucket: aws.String("library-pictures"),
	// })

	fmt.Println(*output.Buckets[0].Name)

	return nil
}
