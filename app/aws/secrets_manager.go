// File: aws/secrets_manager.go

package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// SecretsManager holds the AWS Secrets Manager client
type SecretsManager struct {
	client *secretsmanager.Client
}

// S3Credentials represents the structure of S3 credentials
type S3Credentials struct {
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	Region          string `json:"region"`
	Bucket          string `json:"bucket"`
}

// RDSCredentials represents the structure of RDS credentials
type RDSCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	DBName   string `json:"dbname"`
}

// NewSecretsManager creates a new SecretsManager instance
func NewSecretsManager(region string) (*SecretsManager, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-south-1"))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	return &SecretsManager{
		client: secretsmanager.NewFromConfig(cfg),
	}, nil
}

// GetSecret retrieves a secret from AWS Secrets Manager
func (sm *SecretsManager) GetSecret(secretName string) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := sm.client.GetSecretValue(context.TODO(), input)
	if err != nil {
		log.Fatal(err.Error())
		return "", fmt.Errorf("error retrieving secret: %w", err)
	}

	if result.SecretString != nil {
		return *result.SecretString, nil
	}

	return "", fmt.Errorf("secret is binary and not supported in this implementation")
}

// GetS3Credentials retrieves and parses S3 credentials
func (sm *SecretsManager) GetS3Credentials(secretName string) (*S3Credentials, error) {
	secretValue, err := sm.GetSecret(secretName)
	if err != nil {
		return nil, err
	}

	var creds S3Credentials
	if err := json.Unmarshal([]byte(secretValue), &creds); err != nil {
		return nil, fmt.Errorf("error parsing S3 credentials: %w", err)
	}

	return &creds, nil
}

// GetRDSCredentials retrieves and parses RDS credentials
func (sm *SecretsManager) GetRDSCredentials(secretName string) (*RDSCredentials, error) {
	secretValue, err := sm.GetSecret(secretName)
	if err != nil {
		return nil, err
	}

	var creds RDSCredentials
	if err := json.Unmarshal([]byte(secretValue), &creds); err != nil {
		return nil, fmt.Errorf("error parsing RDS credentials: %w", err)
	}

	return &creds, nil
}
