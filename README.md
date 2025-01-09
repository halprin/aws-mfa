# aws-mfa

AWS IAM User Credentials using MFA

## Set-up

Modify your `~/.aws/credentials` file.

First, there needs to be a profile section for the long-term credentials associated with the IAM user.

```ini
[default-long-term]
aws_access_key_id     = <Access Key ID>
aws_secret_access_key = <Secret Access Key>
mfa_serial            = <ARN to MFA device assigned to IAM user>
```

Then you add a profile section that you actively use for authenticating to AWS with.

```ini
[default]
aws_access_key_id     = a
aws_secret_access_key = a
aws_session_token     = a
source_profile        = <long-term profile section name, e.g. default-long-term>
```

The `a`s will be replaced with real values after you run `aws-mfa` successfully.

## Usage

```shell
aws-mfa login \
  [--mfa-code=<Current MFA code>] \
  [--profile=<Non-long-term profile name from the credentials file to login with>] \
  [--mfa-device-arn=<ARN of the MFA device used to login>] \
  [--duration=<Length of time, in seconds, to be logged in for>]
```

- `mfa-code` - If unspecified on the command line, it will be queried for.
- `profile` - If unspecified on the command line, it defaults to `default`.
- `mfa-device-arn` - If unspecified on the command line, it is read from the long-term profile's `mfa_serial` value in `~/.aws/credentials`.
- `duration` - If unspecified on the command line, it defaults to 12 hours.
