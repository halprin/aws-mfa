package cli

import (
	"aws-mfa/login"
	"fmt"
	"github.com/teris-io/cli"
	"os"
	"strconv"
)

func Start() {

	loginAction := cli.NewCommand("login", "Login with MFA code").
		WithOption(cli.NewOption("mfa-device-arn", "ARN of the MFA device used to login").WithType(cli.TypeString)).
		WithOption(cli.NewOption("profile", "Non-long-term profile name from the credentials file to login with").WithType(cli.TypeString)).
		WithOption(cli.NewOption("mfa-code", "Current MFA code").WithType(cli.TypeString)).
		WithOption(cli.NewOption("duration", "Length of time, in seconds, to be logged in for").WithType(cli.TypeInt)).
		WithAction(func(args []string, options map[string]string) int {
			mfaDeviceArn := ""
			mfaDeviceArnString, exists := options["mfa-device-arn"]
			if exists {
				mfaDeviceArn = mfaDeviceArnString
			}

			profile := ""
			profileString, exists := options["profile"]
			if exists {
				profile = profileString
			}

			mfaCode := ""
			mfaCodeString, exists := options["mfa-code"]
			if exists {
				mfaCode = mfaCodeString
			}

			duration := 0
			durationString, exists := options["duration"]
			if exists {
				var err error
				duration, err = strconv.Atoi(durationString)
				if err != nil {
					fmt.Println("Invalid duration")
					fmt.Println(err.Error())
					return 1
				}
			}

			err := login.GetTokenAndSave(mfaDeviceArn, profile, mfaCode, duration)
			if err != nil {
				fmt.Println("Error!")
				fmt.Println(err.Error())
				return 2
			}

			fmt.Println("Login successful!")

			return 0
		})

	cliApplication := cli.New("MFA login into AWS").
		WithCommand(loginAction)

	os.Exit(cliApplication.Run(os.Args, os.Stdout))
}
