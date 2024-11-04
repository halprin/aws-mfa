package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type StsService struct {
	MfaDeviceArn string
	Profile      string
	MfaCode      string
	Duration     int32
}

func (service *StsService) GetSessionToken() (string, string, string, error) {
	awsConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(service.Profile))
	if err != nil {
		return "", "", "", fmt.Errorf("failed to load the AWS config: %v", err)
	}

	stsService := sts.NewFromConfig(awsConfig)

	sessionTokenInput := &sts.GetSessionTokenInput{
		DurationSeconds: aws.Int32(service.Duration),
		SerialNumber:    aws.String(service.MfaDeviceArn),
		TokenCode:       aws.String(service.MfaCode),
	}

	sessionTokenOutput, err := stsService.GetSessionToken(context.TODO(), sessionTokenInput)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get session token, %v", err)
	}

	return *sessionTokenOutput.Credentials.AccessKeyId, *sessionTokenOutput.Credentials.SecretAccessKey, *sessionTokenOutput.Credentials.SessionToken, nil
}
