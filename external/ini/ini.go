package ini

import (
	"fmt"
	"gopkg.in/ini.v1"
)

type AwsCredentials struct {
	filePath  string
	iniConfig *ini.File
}

func NewAwsCredentials(filePath string) (*AwsCredentials, error) {
	config, err := ini.Load(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS credentials at %s: %v", filePath, err)
	}

	return &AwsCredentials{
		filePath:  filePath,
		iniConfig: config,
	}, nil
}

func (receiver *AwsCredentials) ReadValue(profile string, key string) (string, error) {
	section, err := receiver.iniConfig.GetSection(profile)
	if err != nil {
		return "", fmt.Errorf("failed to find %s profile: %v", profile, err)
	}

	value, err := section.GetKey(key)
	if err != nil {
		return "", fmt.Errorf("failed to find %s key in %s profile: %v", key, profile, err)
	}

	return value.String(), nil
}

func (receiver *AwsCredentials) WriteValue(profile string, key string, value string) error {
	section, err := receiver.iniConfig.GetSection(profile)
	if err != nil {
		return fmt.Errorf("failed to find %s profile: %v", profile, err)
	}

	configValue, err := section.GetKey(key)
	if err != nil {
		return fmt.Errorf("failed to find %s key in %s profile: %v", key, profile, err)
	}

	configValue.SetValue(value)

	return nil
}

func (receiver *AwsCredentials) Save() error {
	return receiver.iniConfig.SaveTo(receiver.filePath)
}
