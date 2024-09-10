package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Config struct {
	StorageType string
	// Other configuration fields...
}

func LoadConfig() (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-south-1"))
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %v", err)
	}

	client := s3.NewFromConfig(cfg)

	result, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		fmt.Println("Coudn't list bucket")
	}

	println("NUmber of data %d", len(result.Buckets))

	return client, nil
}
