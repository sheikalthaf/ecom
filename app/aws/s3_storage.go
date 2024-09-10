package aws

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Config struct {
	Bucket string
}

var S3ConfigInstance *S3Config

type S3ImageStorage struct {
	Client *s3.Client
	bucket string
}

func InitS3Config() {
	sm, err := NewSecretsManager("ap-south-1")
	if err != nil {
		log.Fatal(err.Error())
	}

	secret, err := sm.GetSecret("productimages")
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(secret)
	json.Unmarshal([]byte(secret), &S3ConfigInstance)
}

func NewS3Storage(client *s3.Client) *S3ImageStorage {
	return &S3ImageStorage{
		Client: client,
		bucket: S3ConfigInstance.Bucket,
	}
}

func (s *S3ImageStorage) SaveImage(file multipart.File, imageName string) error {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	_, err = s.Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(imageName),
		Body:   bytes.NewReader(fileBytes),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}

	return nil
}

func (s *S3ImageStorage) AppendUrl(imagePath string) string {
	return fmt.Sprintf("https://%s.s3.ap-south-1.amazonaws.com/%s", s.bucket, imagePath)
}

func (s *S3ImageStorage) DeleteImage(imgPath string, thumbnailPath string) error {
	err := s.delete(imgPath)
	if err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	err = s.delete(thumbnailPath)
	if err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	return nil
}

func (s *S3ImageStorage) delete(imageName string) error {
	_, err := s.Client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(imageName),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	return nil
}

func (s *S3ImageStorage) ImageInit() error {
	// Do nothing
	return nil
}
