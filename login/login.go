package login

import (
	"aws-mfa/external/aws"
	"aws-mfa/external/ini"
	"fmt"
	"os"
	"path/filepath"
)

func GetTokenAndSave(mfaDeviceArn string, profile string, mfaCode string, duration int) error {

	profile, duration = setSensibleDefaults(profile, duration)

	awsCredentials, err := getAwsCredentialsFile()
	if err != nil {
		return err
	}

	sourceProfile, err := getSourceProfileFromConfig(awsCredentials, profile)
	if err != nil {
		return err
	}

	mfaDeviceArn, err = getMfaDeviceFromConfigUnlessSpecified(mfaDeviceArn, awsCredentials, sourceProfile)
	if err != nil {
		return err
	}

	fmt.Printf("Logging into AWS for profile %s for %d seconds using MFA device %s\n", profile, duration, mfaDeviceArn)

	mfaCode, err = getMfaCodeUnlessSpecified(mfaCode)

	stsService := aws.StsService{
		MfaDeviceArn: mfaDeviceArn,
		Profile:      sourceProfile,
		MfaCode:      mfaCode,
		Duration:     int32(duration),
	}

	accessKey, secretAccessKey, sessionToken, err := stsService.GetSessionToken()
	if err != nil {
		return fmt.Errorf("unable to login: %v", err)
	}

	err = writeValuesToConfig(awsCredentials, profile, accessKey, secretAccessKey, sessionToken)
	if err != nil {
		return err
	}

	err = awsCredentials.Save()
	if err != nil {
		return fmt.Errorf("unable to save login credentials: %v", err)
	}

	return nil
}

func setSensibleDefaults(profile string, duration int) (string, int) {
	if profile == "" {
		profile = "default"
	}

	if duration == 0 {
		duration = 43200 // 12 hours
	}

	return profile, duration
}

func getAwsCredentialsFile() (*ini.AwsCredentials, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("unable to get your user's home directory: %v", err)
	}

	credentialsPath := filepath.Join(homeDir, ".aws", "credentials")

	awsCredentials, err := ini.NewAwsCredentials(credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("unable to load your AWS credentials file: %v", err)
	}

	return awsCredentials, nil
}

func getSourceProfileFromConfig(credentials *ini.AwsCredentials, profile string) (string, error) {
	sourceProfile, err := credentials.ReadValue(profile, "source_profile")
	if err != nil {
		return "", fmt.Errorf("unable to read source_profile from profile %s: %v", profile, err)
	}

	return sourceProfile, nil
}

func getMfaDeviceFromConfigUnlessSpecified(mfaDeviceArn string, credentials *ini.AwsCredentials, sourceProfile string) (string, error) {
	if mfaDeviceArn != "" {
		return mfaDeviceArn, nil
	}

	mfaDeviceArn, err := credentials.ReadValue(sourceProfile, "mfa_serial")
	if err != nil {
		return "", fmt.Errorf("unable to read mfa_serial from profile %s: %v", sourceProfile, err)
	}

	return mfaDeviceArn, nil
}

func getMfaCodeUnlessSpecified(mfaCode string) (string, error) {
	if mfaCode != "" {
		return mfaCode, nil
	}

	fmt.Print("Enter MFA code: ")

	_, err := fmt.Scanln(&mfaCode)
	if err != nil {
		return "", fmt.Errorf("unable to read MFA code from standard in: %v", err)
	}

	return mfaCode, nil
}

func writeValuesToConfig(credentials *ini.AwsCredentials, profile string, accessKey string, secretAccessKey string, sessionToken string) error {
	err := credentials.WriteValue(profile, "aws_access_key_id", accessKey)
	if err != nil {
		return fmt.Errorf("unable to write aws_access_key_id to profile %s: %v", profile, err)
	}

	err = credentials.WriteValue(profile, "aws_secret_access_key", secretAccessKey)
	if err != nil {
		return fmt.Errorf("unable to write aws_secret_access_key to profile %s: %v", profile, err)
	}

	err = credentials.WriteValue(profile, "aws_session_token", sessionToken)
	if err != nil {
		return fmt.Errorf("unable to write aws_session_token to profile %s: %v", profile, err)
	}

	return nil
}
