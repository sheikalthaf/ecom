package utilities

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Handler struct {
	s3Client *s3.Client
	bucket   string
}

func NewS3Handler(bucket string) (*S3Handler, error) {
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

	return &S3Handler{
		s3Client: client,
		bucket:   bucket,
	}, nil
}

func (h *S3Handler) UploadFile(ctx context.Context, key string, file multipart.File) error {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	_, err = h.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(h.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(fileBytes),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}

	return nil
}

func (h *S3Handler) DeleteFile(ctx context.Context, key string) error {
	_, err := h.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(h.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	return nil
}

func (h *S3Handler) GetFileURL(key string) string {
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", h.bucket, key)
}
